package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ChatServerGolang/config"
	_ "ChatServerGolang/docs"
	"ChatServerGolang/internal/database"
	"ChatServerGolang/internal/email"
	"ChatServerGolang/internal/handler"
	"ChatServerGolang/internal/middleware"
	"ChatServerGolang/internal/repository"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/internal/sms"
	"ChatServerGolang/internal/ws"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	callRepo := repository.NewCallRepository(db)
	accRepo := repository.NewAccountSettingRepository(db)
	contactRepo := repository.NewContactRepository(db)
	pollRepo := repository.NewPollRepository(db)
	stickerRepo := repository.NewStickerRepository(db)
	draftRepo := repository.NewDraftRepository(db)
	schedMsgRepo := repository.NewScheduledMessageRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	botRepo := repository.NewBotRepository(db)
	gifRepo := repository.NewSavedGifRepository(db)
	verRepo := repository.NewVerificationRepository(db)

	// WebSocket hub
	hub := ws.NewHub()
	go hub.Run()

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo, chatRepo, accRepo)
	callService := service.NewCallService(callRepo, chatRepo, userRepo, userService)
	messageService := service.NewMessageService(messageRepo, chatRepo, userRepo, userService)
	chatService := service.NewChatService(chatRepo, userRepo, messageRepo, userService)
	pushService := service.NewPushService(userRepo, cfg)
	contactService := service.NewContactService(contactRepo)
	sysMsgService := service.NewSystemMessageService(messageRepo, chatRepo, hub)
	_ = service.NewTypingService(hub)
	_ = service.NewMentionService(userRepo, messageRepo)
	pollService := service.NewPollService(pollRepo, chatRepo, messageRepo, sysMsgService)
	stickerService := service.NewStickerService(stickerRepo)
	draftService := service.NewDraftService(draftRepo)
	schedMsgService := service.NewScheduledMessageService(schedMsgRepo, messageRepo, chatRepo)
	sessionService := service.NewSessionService(sessionRepo)
	botService := service.NewBotService(botRepo)
	gifService := service.NewSavedGifService(gifRepo)
	emailSender := email.NewSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)
	smsSender := sms.NewSender(cfg.TwilioAccountSID, cfg.TwilioAuthToken, cfg.TwilioPhone)
	verService := service.NewVerificationService(verRepo, userRepo, emailSender, smsSender)

	// Handlers
	loginCodeHandler := handler.NewLoginCodeHandler(verService, authService)
	verHandler := handler.NewVerificationHandler(verService)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, pushService)
	chatHandler := handler.NewChatHandler(chatService)
	messageHandler := handler.NewMessageHandler(messageService)
	callHandler := handler.NewCallHandler(callService)
	contactHandler := handler.NewContactHandler(contactService)
	pollHandler := handler.NewPollHandler(pollService)
	stickerHandler := handler.NewStickerHandler(stickerService)
	draftHandler := handler.NewDraftHandler(draftService)
	sessionHandler := handler.NewSessionHandler(sessionService)
	botHandler := handler.NewBotHandler(botService)
	gifHandler := handler.NewGifHandler(gifService)
	schedMsgHandler := handler.NewScheduledMessageHandler(schedMsgService)
	wsHandler := handler.NewWSHandler(hub, authService, userRepo, chatRepo)
	wsEvents := handler.NewWebSocketEvents(hub, chatService, messageService, userService, pushService, callService)

	// Scheduler for delayed messages
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			schedMsgService.SchedulerProcess()
		}
	}()

	// Rate limiter
	apiLimiter := middleware.NewRateLimiter(100, time.Minute)

	// Gin engine
	r := gin.Default()
	r.Use(middleware.CORSMiddleware(cfg))

	r.Static("/uploads", "./uploads")

	// API routes with rate limiter applied early
	api := r.Group("/api")
	api.Use(apiLimiter)
	{
		// Auth (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/login/email", loginCodeHandler.SendEmailCode)
			auth.POST("/login/email/verify", loginCodeHandler.VerifyEmailCode)
			auth.POST("/login/phone", loginCodeHandler.SendPhoneCode)
			auth.POST("/login/phone/verify", loginCodeHandler.VerifyPhoneCode)
		}

		// Authorized routes
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware(authService))
		{
			// Auth (authenticated)
			authorized.GET("/auth/refresh", authHandler.RefreshToken)
			authorized.PUT("/auth/change-password", authHandler.ChangePassword)

			// Users
			authorized.GET("/users/profile", userHandler.GetProfile)
			authorized.PUT("/users/profile", userHandler.UpdateProfile)
			authorized.DELETE("/users/account", userHandler.DeleteAccount)
			authorized.PUT("/users/push-token", userHandler.UpdatePushToken)
			authorized.POST("/users/push-test", userHandler.TestPush)
			authorized.PUT("/users/status", userHandler.UpdateStatus)
			authorized.GET("/users/search", userHandler.SearchUsers)
			authorized.GET("/users/:id", userHandler.GetUserByID)
			authorized.GET("/users/username/:username", userHandler.GetUserByUsername)
			authorized.POST("/users/block", userHandler.BlockUser)
			authorized.DELETE("/users/block/:userId", userHandler.UnblockUser)
			authorized.GET("/users/blocked", userHandler.GetBlockedUsers)

			// Account settings
			authorized.GET("/account/settings", userHandler.GetAccountSetting)
			authorized.PUT("/account/settings", userHandler.UpdateAccountSetting)

			// Contacts
			authorized.POST("/contacts/sync", contactHandler.SyncContacts)
			authorized.GET("/contacts", contactHandler.GetContacts)
			authorized.GET("/contacts/search", contactHandler.SearchByPhone)
			authorized.GET("/contacts/registered", contactHandler.FindRegistered)

			// Avatar upload
			authorized.POST("/users/avatar", userHandler.UploadAvatar)

			// Chats
			authorized.GET("/chats", chatHandler.ListChats)
			authorized.GET("/chats/search", chatHandler.SearchChats)
			authorized.GET("/chats/archived", chatHandler.ListArchivedChats)
			authorized.POST("/chats", wsEvents.WrapCreateChat(chatHandler.CreateChat))
			authorized.GET("/chats/:id", chatHandler.GetChat)
			authorized.PUT("/chats/:id", chatHandler.UpdateGroup)
			authorized.DELETE("/chats/:id", wsEvents.WrapDeleteChat(chatHandler.DeleteChat))
			authorized.POST("/chats/:id/participants", chatHandler.AddParticipant)
			authorized.DELETE("/chats/:id/participants/:userId", chatHandler.RemoveParticipant)
			authorized.PUT("/chats/:id/participants/:userId/role", chatHandler.SetRole)
			authorized.POST("/chats/:id/leave", chatHandler.LeaveGroup)
			authorized.POST("/chats/:id/read", chatHandler.MarkAsRead)
			authorized.POST("/chats/:id/pin", chatHandler.PinChat)
			authorized.DELETE("/chats/:id/pin", chatHandler.UnpinChat)
			authorized.POST("/chats/:id/archive", chatHandler.ArchiveChat)
			authorized.POST("/chats/:id/unarchive", chatHandler.UnarchiveChat)
			authorized.POST("/chats/:id/hide", chatHandler.HideChat)
			authorized.POST("/chats/:id/transfer-ownership", chatHandler.TransferOwnership)
			authorized.PUT("/chats/:id/notifications", userHandler.SetNotificationMuted)
			authorized.GET("/chats/:id/notifications", userHandler.IsNotificationMuted)

			// Messages
			authorized.GET("/chats/:id/messages", messageHandler.GetMessages)
			authorized.GET("/chats/:id/messages/search", messageHandler.SearchMessages)
			authorized.GET("/chats/:id/media", messageHandler.GetChatMedia)
			authorized.POST("/chats/:id/messages", wsEvents.WrapSendMessage(messageHandler.SendMessage))
			authorized.POST("/chats/:id/messages/file", messageHandler.UploadFile)
			authorized.POST("/chats/:id/messages/:msgId/resend", messageHandler.ResendMessage)
			authorized.GET("/chats/:id/pinned", messageHandler.GetPinned)
			authorized.GET("/chats/:id/export", messageHandler.ExportChat)
			authorized.GET("/messages/search", messageHandler.SearchAllMessages)
			authorized.GET("/messages/starred", messageHandler.GetStarredMessages)
			authorized.POST("/messages/forward", messageHandler.ForwardMessage)
			authorized.POST("/messages/schedule", schedMsgHandler.Schedule)
			authorized.GET("/messages/scheduled", schedMsgHandler.GetScheduled)
			authorized.DELETE("/messages/scheduled/:id", schedMsgHandler.CancelScheduled)
			authorized.GET("/messages/:id", messageHandler.GetMessageByID)
			authorized.PUT("/messages/:id", wsEvents.WrapEditMessage(messageHandler.EditMessage))
			authorized.DELETE("/messages/:id", wsEvents.WrapDeleteMessage(messageHandler.DeleteMessage))
			authorized.DELETE("/messages/:id/for-me", messageHandler.DeleteMessageForMe)
			authorized.POST("/messages/:id/reactions", messageHandler.AddReaction)
			authorized.DELETE("/messages/:id/reactions", messageHandler.RemoveReaction)
			authorized.PUT("/messages/:id/pin", messageHandler.TogglePin)
			authorized.POST("/messages/:id/star", messageHandler.StarMessage)
			authorized.DELETE("/messages/:id/star", messageHandler.UnstarMessage)
			authorized.POST("/messages/:id/read", messageHandler.MarkMessageRead)
			authorized.GET("/files/:filename", messageHandler.DownloadFile)

			// Polls
			authorized.POST("/chats/:id/polls", pollHandler.CreatePoll)
			authorized.GET("/chats/:id/polls", pollHandler.GetPolls)
			authorized.POST("/polls/:pollId/vote", pollHandler.Vote)
			authorized.POST("/polls/:pollId/close", pollHandler.ClosePoll)

			// Stickers
			authorized.GET("/stickers/packs", stickerHandler.ListPacks)
			authorized.GET("/stickers/packs/my", stickerHandler.GetMyPacks)
			authorized.POST("/stickers/packs", stickerHandler.CreatePack)
			authorized.GET("/stickers/packs/:id", stickerHandler.GetPack)
			authorized.POST("/stickers/packs/:id/stickers", stickerHandler.AddSticker)
			authorized.DELETE("/stickers/packs/:id", stickerHandler.DeletePack)
			authorized.GET("/stickers/library", stickerHandler.GetLibrary)
			authorized.POST("/stickers/library", stickerHandler.AddToLibrary)

			// Drafts
			authorized.POST("/drafts", draftHandler.SaveDraft)
			authorized.GET("/drafts", draftHandler.GetDraft)
			authorized.DELETE("/drafts/:id", draftHandler.DeleteDraft)

			// Sessions
			authorized.GET("/sessions", sessionHandler.GetSessions)
			authorized.DELETE("/sessions/:id", sessionHandler.DeleteSession)
			authorized.DELETE("/sessions", sessionHandler.DeleteAllSessions)

			// Bots
			authorized.POST("/bots", botHandler.CreateBot)
			authorized.GET("/bots", botHandler.GetMyBots)
			authorized.PUT("/bots/:id", botHandler.UpdateBot)
			authorized.DELETE("/bots/:id", botHandler.DeleteBot)
			authorized.POST("/bots/:id/regenerate-token", botHandler.RegenerateToken)

			// Saved GIFs
			authorized.POST("/gifs", gifHandler.SaveGif)
			authorized.GET("/gifs", gifHandler.GetSavedGifs)
			authorized.DELETE("/gifs", gifHandler.DeleteGif)

			// Calls
			authorized.POST("/calls/initiate", wsEvents.WrapInitiateCall(callHandler.InitiateCall))
			authorized.POST("/calls/:id/respond", wsEvents.WrapRespondCall(callHandler.RespondCall))
			authorized.POST("/calls/:id/end", wsEvents.WrapEndCall(callHandler.EndCall))
			authorized.GET("/calls/:id", callHandler.GetCall)
			authorized.GET("/calls/history/:chatId", callHandler.GetCallHistory)

			// Verification
			authorized.POST("/verification/email/send", verHandler.SendEmail)
			authorized.POST("/verification/email/verify", verHandler.VerifyEmail)
			authorized.POST("/verification/phone/send", verHandler.SendPhone)
			authorized.POST("/verification/phone/verify", verHandler.VerifyPhone)
		}
	}

	// WebSocket
	r.GET("/ws", wsHandler.HandleWebSocket)

	// Frontend static
	r.Static("/app", "./frontend")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/app/")
	})

	// Healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Startup banner
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║         CHAT MESSENGER SERVER                       ║")
	fmt.Println("╠══════════════════════════════════════════════════════╣")
	fmt.Printf("║  Server:   http://localhost:%s                       \n", cfg.ServerPort)
	fmt.Printf("║  Frontend: http://localhost:%s/app/                  \n", cfg.ServerPort)
	fmt.Printf("║  Swagger:  http://localhost:%s/swagger/index.html    \n", cfg.ServerPort)
	fmt.Printf("║  Health:   http://localhost:%s/health                \n", cfg.ServerPort)
	fmt.Println("╚══════════════════════════════════════════════════════╝")
	fmt.Println()

	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	db.Close()
	log.Println("Server exited gracefully")
}
