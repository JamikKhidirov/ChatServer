package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ChatServerGolang/config"
	_ "ChatServerGolang/docs"
	"ChatServerGolang/internal/database"
	"ChatServerGolang/internal/handler"
	"ChatServerGolang/internal/middleware"
	"ChatServerGolang/internal/repository"
	"ChatServerGolang/internal/service"
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
	defer db.Close()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	callRepo := repository.NewCallRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo, chatRepo)
	callService := service.NewCallService(callRepo, chatRepo, userRepo, userService)
	messageService := service.NewMessageService(messageRepo, chatRepo, userRepo, userService)
	chatService := service.NewChatService(chatRepo, userRepo, messageRepo, userService)
	pushService := service.NewPushService(userRepo, cfg)

	// WebSocket hub
	hub := ws.NewHub()
	go hub.Run()

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, pushService)
	chatHandler := handler.NewChatHandler(chatService)
	messageHandler := handler.NewMessageHandler(messageService)
	callHandler := handler.NewCallHandler(callService)
	wsHandler := handler.NewWSHandler(hub, authService, userRepo, chatRepo)
	wsEvents := handler.NewWebSocketEvents(hub, chatService, messageService, userService, pushService)

	// Gin engine
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Rate limiter for API (100 requests/min per IP)
	apiLimiter := middleware.NewRateLimiter(100, time.Minute)

	// Serve uploaded files
	r.Static("/uploads", "./uploads")

	// API routes
	api := r.Group("/api")
	{
		// Auth (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Authorized routes
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware(authService))
		{
			// Auth
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

			// Chats
			authorized.GET("/chats", chatHandler.ListChats)
			authorized.POST("/chats", wsEvents.WrapCreateChat(chatHandler.CreateChat))
			authorized.GET("/chats/:id", chatHandler.GetChat)
			authorized.DELETE("/chats/:id", wsEvents.WrapDeleteChat(chatHandler.DeleteChat))
			authorized.POST("/chats/:id/participants", chatHandler.AddParticipant)
			authorized.DELETE("/chats/:id/participants/:userId", chatHandler.RemoveParticipant)
			authorized.POST("/chats/:id/read", chatHandler.MarkAsRead)
			authorized.PUT("/chats/:id/participants/:userId/role", chatHandler.PromoteAdmin)
			authorized.POST("/chats/:id/leave", chatHandler.LeaveGroup)
			authorized.PUT("/chats/:id", chatHandler.UpdateGroup)
			authorized.PUT("/chats/:id/notifications", userHandler.SetNotificationMuted)
			authorized.GET("/chats/:id/notifications", userHandler.IsNotificationMuted)

			// Messages
			authorized.GET("/chats/:id/messages", messageHandler.GetMessages)
			authorized.GET("/chats/:id/messages/search", messageHandler.SearchMessages)
			authorized.POST("/chats/:id/messages", wsEvents.WrapSendMessage(messageHandler.SendMessage))
			authorized.POST("/chats/:id/messages/file", messageHandler.UploadFile)
			authorized.POST("/chats/:id/messages/:msgId/resend", messageHandler.ResendMessage)
			authorized.PUT("/messages/:id", wsEvents.WrapEditMessage(messageHandler.EditMessage))
			authorized.DELETE("/messages/:id", wsEvents.WrapDeleteMessage(messageHandler.DeleteMessage))
			authorized.GET("/messages/:id", messageHandler.GetMessageByID)
			authorized.POST("/messages/:id/reactions", messageHandler.AddReaction)
			authorized.DELETE("/messages/:id/reactions", messageHandler.RemoveReaction)
			authorized.PUT("/messages/:id/pin", messageHandler.TogglePin)
			authorized.POST("/messages/:id/read", messageHandler.MarkMessageRead)
			authorized.GET("/chats/:id/pinned", messageHandler.GetPinned)

			// Calls
			authorized.POST("/calls/initiate", wsEvents.WrapInitiateCall(callHandler.InitiateCall))
			authorized.POST("/calls/:id/respond", wsEvents.WrapRespondCall(callHandler.RespondCall))
			authorized.POST("/calls/:id/end", wsEvents.WrapEndCall(callHandler.EndCall))
			authorized.GET("/calls/:id", callHandler.GetCall)
			authorized.GET("/calls/history/:chatId", callHandler.GetCallHistory)
		}
	}

	// WebSocket
	r.GET("/ws", wsHandler.HandleWebSocket)

	// Healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Apply rate limiter to API routes
	api.Use(apiLimiter)

	// Startup banner
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║         CHAT MESSENGER SERVER STARTED               ║")
	fmt.Println("╠══════════════════════════════════════════════════════╣")
	fmt.Printf("║  Server:  http://localhost:%s                       \n", cfg.ServerPort)
	fmt.Printf("║  Swagger: http://localhost:%s/swagger/index.html    \n", cfg.ServerPort)
	fmt.Printf("║  Health:  http://localhost:%s/health                \n", cfg.ServerPort)
	fmt.Println("╚══════════════════════════════════════════════════════╝")
	fmt.Println()

	go func() {
		if err := r.Run(":" + cfg.ServerPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
