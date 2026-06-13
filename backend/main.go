// @title Chat Messenger API
// @version 2.0.0
// @description Сервер мессенджера на Go. Поддерживает: личные и групповые чаты, голосовые и видеозвонки, stories, голосовые комнаты, каналы, стикеры, опросы, ботов, геолокацию, кастомные эмодзи, эффекты сообщений, избранное и многое другое.
// @description WebSocket: ws://localhost:8080/ws?token={jwt_token}
// @description
// @description ===================== ДОКУМЕНТАЦИЯ WEB-SOCKET =====================
// @description
// @description Клиент подключается по URL: ws://localhost:8080/ws?token={JWT_TOKEN}
// @description После подключения сервер и клиент обмениваются JSON-сообщениями.
// @description
// @description Формат ВСЕХ сообщений (как от клиента, так и от сервера):
// @description { "type": "название_события", "payload": { ... поля ... } }
// @description
// @description ----------------------------------------------------------------
// @description СОБЫТИЯ ОТ СЕРВЕРА К КЛИЕНТУ (сервер отправляет автоматически):
// @description ----------------------------------------------------------------
// @description
// @description === Сообщения ===
// @description message:new — Новое сообщение. Payload: полный объект Message (id, chatId, senderId, content, type, createdAt и т.д.)
// @description message:edited — Сообщение отредактировано. Payload: Message
// @description message:deleted — Сообщение удалено. Payload: { messageId, chatId }
// @description message:read — Сообщение прочитано. Payload: { messageId, userId, chatId }
// @description message:reaction — Добавлена/удалена реакция. Payload: Message (обновлённый объект с reactions)
// @description message:pinned — Сообщение закреплено/откреплено. Payload: Message
// @description message:starred — Сообщение добавлено в избранное (без payload)
// @description message:forward — Сообщение переслано (без payload)
// @description
// @description === Пользователи ===
// @description user:online — Пользователь стал онлайн. Payload: { userId, online: true }
// @description user:offline — Пользователь стал офлайн. Payload: { userId, online: false }
// @description user:typing — Пользователь печатает. Payload: { chatId, userId }
// @description user:stop_typing — Пользователь перестал печатать. Payload: { chatId, userId }
// @description user:keyboard_opened — Клавиатура открыта (мобильные). Payload: { chatId, userId }
// @description user:keyboard_closed — Клавиатура закрыта. Payload: { chatId, userId }
// @description
// @description === Чаты ===
// @description chat:created — Создан новый чат. Payload: полный объект Chat
// @description chat:updated — Чат обновлён (название, аватар, участники). Payload: Chat
// @description chat:deleted — Чат удалён (без payload)
// @description
// @description === Звонки (WebRTC) ===
// @description call:offer — Входящий звонок. Payload: { chatId, callId, sdp }
// @description call:answer — Ответ на звонок. Payload: { chatId, callId, sdp }
// @description call:ice — ICE-кандидат. Payload: { callId, candidate }
// @description call:end — Звонок завершён. Payload: { callId, userId }
// @description call:missed — Пропущенный звонок (без payload)
// @description call:accept — Звонок принят (без payload)
// @description call:reject — Звонок отклонён. Payload: { callId, userId }
// @description
// @description ----------------------------------------------------------------
// @description СОБЫТИЯ ОТ КЛИЕНТА К СЕРВЕРУ (клиент отправляет, чтобы что-то сделать):
// @description ----------------------------------------------------------------
// @description
// @description === Работа с сообщениями ===
// @description message:send — Отправить сообщение. Payload: { chatId, content, type, replyToId? }
// @description   type: "text" | "image" | "file" | "gif" | "voice" | "video" | "audio" | "location" | "system"
// @description   Пример: { "type": "message:send", "payload": { "chatId": "123", "content": "Привет!", "type": "text" } }
// @description
// @description message:edit — Редактировать сообщение. Payload: { messageId, content }
// @description message:delete — Удалить сообщение. Payload: { messageId, chatId }
// @description message:read — Отметить как прочитанное. Payload: { messageId, chatId }
// @description message:react — Добавить реакцию. Payload: { messageId, emoji } (emoji: "👍", "❤️", "😆", "😮", "😢", "🙏")
// @description message:unreact — Удалить реакцию. Payload: { messageId, emoji }
// @description message:pin — Закрепить/открепить. Payload: { messageId, pin: true/false }
// @description message:star — В избранное. Payload: { messageId }
// @description message:unstar — Из избранного. Payload: { messageId }
// @description message:forward — Переслать сообщение. Payload: { messageId, toChatId }
// @description
// @description === Управление чатами ===
// @description chat:create — Создать чат. Payload: { type, name?, participantIds, description? }
// @description   type: "private" | "group" | "channel". Пример: { "type": "chat:create", "payload": { "type": "group", "name": "Friends", "participantIds": ["user1","user2"] } }
// @description chat:update — Обновить чат. Payload: { chatId, name?, description?, avatarUrl? }
// @description chat:add_participant — Добавить участника. Payload: { chatId, userId }
// @description chat:remove_participant — Удалить участника. Payload: { chatId, userId }
// @description chat:leave — Покинуть чат. Payload: { chatId }
// @description chat:pin — Закрепить чат в списке. Payload: { chatId }
// @description chat:unpin — Открепить чат. Payload: { chatId }
// @description chat:archive — Архивировать чат. Payload: { chatId }
// @description chat:unarchive — Разархивировать чат. Payload: { chatId }
// @description
// @description === Статус пользователя ===
// @description user:typing — Отправить индикатор печатания. Payload: { chatId }
// @description user:stop_typing — Прекратить индикатор. Payload: { chatId }
// @description user:keyboard_opened — Клавиатура открыта. Payload: { chatId }
// @description user:keyboard_closed — Клавиатура закрыта. Payload: { chatId }
// @description user:block — Заблокировать пользователя. Payload: { userId }
// @description user:unblock — Разблокировать пользователя. Payload: { userId }
// @description
// @description === WebRTC звонки ===
// @description call:offer — Отправить WebRTC offer. Payload: { chatId, callId, sdp }
// @description call:answer — WebRTC answer. Payload: { chatId, callId, sdp }
// @description call:ice — ICE candidate. Payload: { callId, candidate }
// @description
// @description =================================================================
// @termsOfService http://localhost:8080/terms
// @contact.name API Support
// @contact.url http://localhost:8080/support
// @contact.email support@chatserver.local
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите JWT токен в формате: Bearer <token>
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

	"ChatServerGolang/backend/internal/config"
	_ "ChatServerGolang/docs"
	"ChatServerGolang/backend/internal/database"
	"ChatServerGolang/backend/internal/email"
	"ChatServerGolang/backend/internal/handler/auth"
	"ChatServerGolang/backend/internal/handler/bot"
	"ChatServerGolang/backend/internal/handler/call"
	"ChatServerGolang/backend/internal/handler/channel"
	"ChatServerGolang/backend/internal/handler/chat"
	"ChatServerGolang/backend/internal/handler/contact"
	"ChatServerGolang/backend/internal/handler/draft"
	"ChatServerGolang/backend/internal/handler/folder"
	"ChatServerGolang/backend/internal/handler/gif"
	"ChatServerGolang/backend/internal/handler/groupcall"
	"ChatServerGolang/backend/internal/handler/link"
	"ChatServerGolang/backend/internal/handler/login"
	"ChatServerGolang/backend/internal/handler/message"
	"ChatServerGolang/backend/internal/handler/poll"
	"ChatServerGolang/backend/internal/handler/schedmsg"
	"ChatServerGolang/backend/internal/handler/session"
	"ChatServerGolang/backend/internal/handler/sticker"
	"ChatServerGolang/backend/internal/handler/story"
	"ChatServerGolang/backend/internal/handler/user"
	"ChatServerGolang/backend/internal/handler/verification"
	"ChatServerGolang/backend/internal/handler/ws"
	"ChatServerGolang/backend/internal/handler/savedmsg"
	"ChatServerGolang/backend/internal/handler/emoji"
	"ChatServerGolang/backend/internal/handler/voicechat"
	"ChatServerGolang/backend/internal/middleware"
	"ChatServerGolang/backend/internal/repository/account"
	"ChatServerGolang/backend/internal/repository/bot"
	"ChatServerGolang/backend/internal/repository/call"
	"ChatServerGolang/backend/internal/repository/chat"
	"ChatServerGolang/backend/internal/repository/contact"
	"ChatServerGolang/backend/internal/repository/draft"
	"ChatServerGolang/backend/internal/repository/folder"
	"ChatServerGolang/backend/internal/repository/gif"
	"ChatServerGolang/backend/internal/repository/channel"
	"ChatServerGolang/backend/internal/repository/groupcall"
	"ChatServerGolang/backend/internal/repository/link"
	"ChatServerGolang/backend/internal/repository/message"
	"ChatServerGolang/backend/internal/repository/poll"
	"ChatServerGolang/backend/internal/repository/schedmsg"
	"ChatServerGolang/backend/internal/repository/session"
	"ChatServerGolang/backend/internal/repository/sticker"
	"ChatServerGolang/backend/internal/repository/story"
	"ChatServerGolang/backend/internal/repository/user"
	"ChatServerGolang/backend/internal/repository/verification"
	"ChatServerGolang/backend/internal/repository/savedmsg"
	"ChatServerGolang/backend/internal/repository/emoji"
	"ChatServerGolang/backend/internal/repository/voicechat"
	"ChatServerGolang/backend/internal/service/auth"
	"ChatServerGolang/backend/internal/service/bot"
	"ChatServerGolang/backend/internal/service/call"
	"ChatServerGolang/backend/internal/service/channel"
	"ChatServerGolang/backend/internal/service/chat"
	"ChatServerGolang/backend/internal/service/contact"
	"ChatServerGolang/backend/internal/service/draft"
	"ChatServerGolang/backend/internal/service/folder"
	"ChatServerGolang/backend/internal/service/gif"
	"ChatServerGolang/backend/internal/service/groupcall"
	"ChatServerGolang/backend/internal/service/link"
	"ChatServerGolang/backend/internal/service/mention"
	"ChatServerGolang/backend/internal/service/message"
	"ChatServerGolang/backend/internal/service/poll"
	"ChatServerGolang/backend/internal/service/push"
	"ChatServerGolang/backend/internal/service/schedmsg"
	"ChatServerGolang/backend/internal/service/session"
	"ChatServerGolang/backend/internal/service/sticker"
	"ChatServerGolang/backend/internal/service/story"
	"ChatServerGolang/backend/internal/service/systemmsg"
	"ChatServerGolang/backend/internal/service/savedmsg"
	"ChatServerGolang/backend/internal/service/emoji"
	"ChatServerGolang/backend/internal/service/voicechat"
	"ChatServerGolang/backend/internal/service/typing"
	"ChatServerGolang/backend/internal/service/user"
	"ChatServerGolang/backend/internal/service/verification"
	"ChatServerGolang/backend/internal/sms"
	"ChatServerGolang/backend/internal/ws"

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

	userRepo := userrepo.NewUserRepository(db)
	chatRepo := chatrepo.NewChatRepository(db)
	messageRepo := messagerepo.NewMessageRepository(db)
	callRepo := callrepo.NewCallRepository(db)
	accRepo := accountrepo.NewAccountSettingRepository(db)
	contactRepo := contactrepo.NewContactRepository(db)
	pollRepo := pollrepo.NewPollRepository(db)
	stickerRepo := stickerrepo.NewStickerRepository(db)
	draftRepo := draftrepo.NewDraftRepository(db)
	schedMsgRepo := schedmsgrepo.NewScheduledMessageRepository(db)
	sessionRepo := sessionrepo.NewSessionRepository(db)
	botRepo := botrepo.NewBotRepository(db)
	gifRepo := gifrepo.NewSavedGifRepository(db)
	verRepo := verrepo.NewVerificationRepository(db)
	linkRepo := linkrepo.NewInviteLinkRepository(db)
	folderRepo := folderrepo.NewChatFolderRepository(db)
	storyRepo := storyrepo.NewStoryRepository(db)
	groupCallRepo := groupcallrepo.NewGroupCallRepository(db)
	savedMsgRepo := savedmsgrepo.NewSavedMessageRepository(db)
	customEmojiRepo := emojirepo.NewCustomEmojiRepository(db)
	voiceChatRepo := voicechatrepo.NewVoiceChatRepository(db)

	hub := ws.NewHub()
	go hub.Run()

	authService := authservice.NewAuthService(userRepo, cfg)
	userService := userservice.NewUserService(userRepo, chatRepo, accRepo)
	callService := callservice.NewCallService(callRepo, chatRepo, userRepo, userService)
	messageService := messageservice.NewMessageService(messageRepo, chatRepo, userRepo, userService)
	chatService := chatservice.NewChatService(chatRepo, userRepo, messageRepo, userService)
	pushService := pushservice.NewPushService(userRepo, cfg)
	contactService := contactservice.NewContactService(contactRepo)
	sysMsgService := systemmsgservice.NewSystemMessageService(messageRepo, chatRepo, hub)
	_ = typingservice.NewTypingService(hub)
	_ = mentionservice.NewMentionService(userRepo, messageRepo)
	pollService := pollservice.NewPollService(pollRepo, chatRepo, messageRepo, sysMsgService)
	stickerService := stickerservice.NewStickerService(stickerRepo)
	draftService := draftservice.NewDraftService(draftRepo)
	schedMsgService := schedmsgservice.NewScheduledMessageService(schedMsgRepo, messageRepo, chatRepo)
	sessionService := sessionservice.NewSessionService(sessionRepo)
	botService := botservice.NewBotService(botRepo)
	gifService := gifservice.NewSavedGifService(gifRepo)
	linkService := linkservice.NewInviteLinkService(linkRepo, chatRepo)
	folderService := folderservice.NewChatFolderService(folderRepo, chatRepo)
	storyService := storyservice.NewStoryService(storyRepo, userRepo, chatRepo)
	groupCallService := groupcallservice.NewGroupCallService(groupCallRepo, chatRepo, userRepo)
	channelRepo := channelrepo.NewChannelSubscriberRepository(db)
	channelService := channelservice.NewChannelService(channelRepo, chatRepo, userRepo)
	savedMsgService := savedmsgservice.NewSavedMessageService(savedMsgRepo, messageRepo, chatRepo, userRepo)
	customEmojiService := emojiservice.NewCustomEmojiService(customEmojiRepo)
	voiceChatService := voicechatservice.NewVoiceChatService(voiceChatRepo, chatRepo, userRepo)
	emailSender := email.NewSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)
	smsSender := sms.NewSender(cfg.TwilioAccountSID, cfg.TwilioAuthToken, cfg.TwilioPhone)
	verService := verservice.NewVerificationService(verRepo, userRepo, emailSender, smsSender)

	loginCodeHandler := loginhandler.NewLoginCodeHandler(verService, authService)
	verHandler := verhandler.NewVerificationHandler(verService)
	authHandler := authhandler.NewAuthHandler(authService)
	userHandler := userhandler.NewUserHandler(userService, pushService)
	chatHandler := chathandler.NewChatHandler(chatService, hub)
	messageHandler := messagehandler.NewMessageHandler(messageService)
	callHandler := callhandler.NewCallHandler(callService)
	contactHandler := contacthandler.NewContactHandler(contactService)
	pollHandler := pollhandler.NewPollHandler(pollService)
	stickerHandler := stickerhandler.NewStickerHandler(stickerService)
	draftHandler := drafthandler.NewDraftHandler(draftService)
	sessionHandler := sessionhandler.NewSessionHandler(sessionService)
	botHandler := bothandler.NewBotHandler(botService)
	gifHandler := gifhandler.NewGifHandler(gifService)
	schedMsgHandler := schedmsghandler.NewScheduledMessageHandler(schedMsgService)
	linkHandler := linkhandler.NewInviteLinkHandler(linkService)
	folderHandler := folderhandler.NewChatFolderHandler(folderService)
	storyHandler := storyhandler.NewStoryHandler(storyService)
	groupCallHandler := groupcallhandler.NewGroupCallHandler(groupCallService)
	channelHandler := channelhandler.NewChannelHandler(channelService, chatService)
	savedMsgHandler := savedmsghandler.NewSavedMessageHandler(savedMsgService)
	emojiHandler := emojihandler.NewEmojiHandler(customEmojiService)
	voiceChatHandler := voicechathandler.NewVoiceChatHandler(voiceChatService)
	wsHandler := wshandler.NewWSHandler(hub, authService, userRepo, chatRepo)
	wsEvents := wshandler.NewWebSocketEvents(hub, chatService, messageService, userService, pushService, callService)

	// Scheduler for delayed messages
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			schedMsgService.SchedulerProcess()
		}
	}()

	apiLimiter := middleware.NewRateLimiter(100, time.Minute)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware(cfg))

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	api.Use(apiLimiter)
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/admin/register", authHandler.RegisterAdmin)
			auth.POST("/login", authHandler.Login)
			auth.POST("/login/email", loginCodeHandler.SendEmailCode)
			auth.POST("/login/email/verify", loginCodeHandler.VerifyEmailCode)
			auth.POST("/login/phone", loginCodeHandler.SendPhoneCode)
			auth.POST("/login/phone/verify", loginCodeHandler.VerifyPhoneCode)
		}

		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware(authService))
		{
			authorized.GET("/auth/refresh", authHandler.RefreshToken)
			authorized.PUT("/auth/change-password", authHandler.ChangePassword)

			authorized.GET("/users/profile", userHandler.GetProfile)
			authorized.PUT("/users/profile", userHandler.UpdateProfile)
			authorized.PUT("/users/username", userHandler.ChangeUsername)
			authorized.PUT("/users/email", userHandler.ChangeEmail)
			authorized.DELETE("/users/account", userHandler.DeleteAccount)
			authorized.PUT("/users/push-token", userHandler.UpdatePushToken)
			authorized.POST("/users/push-test", userHandler.TestPush)
			authorized.PUT("/users/status", userHandler.UpdateStatus)
			authorized.GET("/users/search", userHandler.SearchUsers)
			authorized.GET("/users/:id", userHandler.GetUserByID)
			authorized.GET("/users/username/:username", userHandler.GetUserByUsername)
			authorized.GET("/users/:id/last-seen", userHandler.GetLastSeen)
			authorized.POST("/users/block", userHandler.BlockUser)
			authorized.DELETE("/users/block/:userId", userHandler.UnblockUser)
			authorized.GET("/users/blocked", userHandler.GetBlockedUsers)

			authorized.GET("/account/settings", userHandler.GetAccountSetting)
			authorized.PUT("/account/settings", userHandler.UpdateAccountSetting)

			authorized.POST("/contacts/sync", contactHandler.SyncContacts)
			authorized.GET("/contacts", contactHandler.GetContacts)
			authorized.GET("/contacts/search", contactHandler.SearchByPhone)
			authorized.GET("/contacts/registered", contactHandler.FindRegistered)
			authorized.POST("/contacts/photo", contactHandler.UpdateContactPhoto)

			authorized.POST("/users/avatar", userHandler.UploadAvatar)

			authorized.GET("/chats", chatHandler.ListChats)
			authorized.GET("/chats/search", chatHandler.SearchChats)
			authorized.GET("/chats/archived", chatHandler.ListArchivedChats)
			authorized.POST("/chats", wsEvents.WrapCreateChat(chatHandler.CreateChat))
			authorized.POST("/chats/start/:userId", chatHandler.StartPrivateChat)
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
			authorized.POST("/chats/:id/promote", chatHandler.PromoteToAdmin)
			authorized.POST("/chats/:id/demote", chatHandler.DemoteFromAdmin)
			authorized.POST("/chats/:id/photo", chatHandler.UploadChatPhoto)
			authorized.POST("/chats/:id/wallpaper", chatHandler.SetChatWallpaper)
			authorized.PUT("/chats/:id/permissions", chatHandler.SetChatPermissions)
			authorized.GET("/chats/:id/online", chatHandler.GetOnlineMembers)
			authorized.PUT("/chats/:id/notifications", userHandler.SetNotificationMuted)
			authorized.GET("/chats/:id/notifications", userHandler.IsNotificationMuted)
			authorized.PUT("/chats/:id/slow-mode", chatHandler.SetSlowMode)

			// Invite links
			authorized.POST("/chats/:id/invite-links", linkHandler.CreateInviteLink)
			authorized.GET("/chats/:id/invite-links", linkHandler.GetInviteLinks)
			authorized.DELETE("/chats/:id/invite-links/:linkId", linkHandler.DeleteInviteLink)
			authorized.POST("/chats/join", linkHandler.JoinByInviteLink)

			// Chat folders
			authorized.GET("/folders", folderHandler.ListFolders)
			authorized.POST("/folders", folderHandler.CreateFolder)
			authorized.PUT("/folders/:id", folderHandler.UpdateFolder)
			authorized.DELETE("/folders/:id", folderHandler.DeleteFolder)

			authorized.GET("/chats/:id/messages", messageHandler.GetMessages)
			authorized.GET("/chats/:id/messages/search", messageHandler.SearchMessages)
			authorized.GET("/chats/:id/media", messageHandler.GetChatMedia)
			authorized.POST("/chats/:id/messages", wsEvents.WrapSendMessage(messageHandler.SendMessage))
			authorized.POST("/chats/:id/messages/file", messageHandler.UploadFile)
			authorized.POST("/chats/:id/messages/voice", messageHandler.UploadVoice)
			authorized.POST("/chats/:id/messages/:msgId/resend", messageHandler.ResendMessage)
			authorized.GET("/chats/:id/pinned", messageHandler.GetPinned)
			authorized.GET("/chats/:id/export", messageHandler.ExportChat)
			authorized.GET("/messages/search", messageHandler.SearchAllMessages)
			authorized.GET("/messages/starred", messageHandler.GetStarredMessages)
			authorized.POST("/messages/read/bulk", messageHandler.BulkMarkRead)
			authorized.DELETE("/messages/bulk", messageHandler.BulkDeleteMessages)
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
			authorized.POST("/messages/:id/self-destruct", messageHandler.SetSelfDestruct)
			authorized.POST("/messages/:id/read", messageHandler.MarkMessageRead)
			authorized.POST("/messages/:id/report", messageHandler.ReportMessage)
			authorized.GET("/messages/:id/history", messageHandler.GetMessageHistory)
			authorized.GET("/files/:filename", messageHandler.DownloadFile)

			authorized.POST("/chats/:id/polls", pollHandler.CreatePoll)
			authorized.GET("/chats/:id/polls", pollHandler.GetPolls)
			authorized.POST("/polls/:pollId/vote", pollHandler.Vote)
			authorized.POST("/polls/:pollId/close", pollHandler.ClosePoll)
			
			authorized.GET("/stickers/packs", stickerHandler.ListPacks)
			authorized.GET("/stickers/packs/my", stickerHandler.GetMyPacks)
			authorized.POST("/stickers/packs", stickerHandler.CreatePack)
			authorized.GET("/stickers/packs/:id", stickerHandler.GetPack)
			authorized.POST("/stickers/packs/:id/stickers", stickerHandler.AddSticker)
			authorized.DELETE("/stickers/packs/:id", stickerHandler.DeletePack)
			authorized.GET("/stickers/library", stickerHandler.GetLibrary)
			authorized.POST("/stickers/library", stickerHandler.AddToLibrary)

			authorized.POST("/drafts", draftHandler.SaveDraft)
			authorized.GET("/drafts", draftHandler.GetDraft)
			authorized.DELETE("/drafts/:id", draftHandler.DeleteDraft)

			authorized.GET("/sessions", sessionHandler.GetSessions)
			authorized.DELETE("/sessions/:id", sessionHandler.DeleteSession)
			authorized.DELETE("/sessions", sessionHandler.DeleteAllSessions)

			authorized.POST("/bots", botHandler.CreateBot)
			authorized.GET("/bots", botHandler.GetMyBots)
			authorized.PUT("/bots/:id", botHandler.UpdateBot)
			authorized.DELETE("/bots/:id", botHandler.DeleteBot)
			authorized.POST("/bots/:id/regenerate-token", botHandler.RegenerateToken)

			authorized.POST("/gifs", gifHandler.SaveGif)
			authorized.GET("/gifs", gifHandler.GetSavedGifs)
			authorized.DELETE("/gifs", gifHandler.DeleteGif)

			// Stories
			authorized.POST("/stories", storyHandler.CreateStory)
			authorized.GET("/stories", storyHandler.GetFollowingStories)
			authorized.GET("/stories/my", storyHandler.GetMyStories)
			authorized.GET("/stories/:id", storyHandler.GetStoryByID)
			authorized.DELETE("/stories/:id", storyHandler.DeleteStory)
			authorized.GET("/stories/:id/views", storyHandler.GetStoryViews)

			// Group calls
			authorized.POST("/calls/group/initiate", groupCallHandler.InitiateGroupCall)
			authorized.POST("/calls/group/respond", groupCallHandler.JoinGroupCall)
			authorized.POST("/calls/group/:id/end", groupCallHandler.EndGroupCall)
			authorized.GET("/calls/group/:id", groupCallHandler.GetGroupCall)
			authorized.GET("/chats/:id/active-calls", groupCallHandler.GetActiveGroupCalls)

			// Broadcast channels
			authorized.POST("/channels/subscribe", channelHandler.Subscribe)
			authorized.POST("/channels/unsubscribe", channelHandler.Unsubscribe)
			authorized.GET("/channels", channelHandler.GetMyChannels)
			authorized.GET("/channels/:id/subscribers", channelHandler.GetSubscribers)
			authorized.GET("/channels/:id/subscribed", channelHandler.IsSubscribed)
			authorized.PUT("/channels/:id/subscribers/:userId/role", channelHandler.SetSubscriberRole)

			// Video circle messages
			authorized.POST("/chats/:id/messages/video-circle", messageHandler.UploadVideoCircle)

			// Location messages
			authorized.POST("/chats/:id/messages/location", messageHandler.SendLocation)

			// Message effects (via SendMessage with effect field)
			// Saved messages
			authorized.POST("/messages/:id/save", savedMsgHandler.SaveMessage)
			authorized.GET("/saved-messages", savedMsgHandler.GetSavedMessages)
			authorized.DELETE("/saved-messages/:id", savedMsgHandler.DeleteSavedMessage)

			// Custom emojis
			authorized.POST("/emojis", emojiHandler.CreateEmoji)
			authorized.GET("/emojis", emojiHandler.GetAllEmojis)
			authorized.GET("/emojis/my", emojiHandler.GetMyEmojis)
			authorized.DELETE("/emojis/:id", emojiHandler.DeleteEmoji)

			// Voice chats
			authorized.POST("/chats/:id/voice-chat", voiceChatHandler.CreateVoiceChat)
			authorized.GET("/chats/:id/voice-chats/active", voiceChatHandler.GetActiveVoiceChats)
			authorized.GET("/chats/:id/voice-chats/history", voiceChatHandler.GetVoiceChatHistory)
			authorized.GET("/voice-chats/:id", voiceChatHandler.GetVoiceChat)
			authorized.POST("/voice-chats/:id/join", voiceChatHandler.JoinVoiceChat)
			authorized.POST("/voice-chats/:id/leave", voiceChatHandler.LeaveVoiceChat)
			authorized.POST("/voice-chats/:id/end", voiceChatHandler.EndVoiceChat)
			authorized.POST("/voice-chats/:id/mute", voiceChatHandler.MuteParticipant)

			authorized.POST("/calls/initiate", wsEvents.WrapInitiateCall(callHandler.InitiateCall))
			authorized.POST("/calls/:id/respond", wsEvents.WrapRespondCall(callHandler.RespondCall))
			authorized.POST("/calls/:id/end", wsEvents.WrapEndCall(callHandler.EndCall))
			authorized.GET("/calls/:id", callHandler.GetCall)
			authorized.GET("/calls/history/:chatId", callHandler.GetCallHistory)

			authorized.POST("/verification/email/send", verHandler.SendEmail)
			authorized.POST("/verification/email/verify", verHandler.VerifyEmail)
			authorized.POST("/verification/phone/send", verHandler.SendPhone)
			authorized.POST("/verification/phone/verify", verHandler.VerifyPhone)
		}
	}

	r.GET("/ws", wsHandler.HandleWebSocket)
	r.Static("/assets", "./frontend/dist/assets")
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
	r.GET("/app", func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "version": "2.0.0"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.StaticFile("/postman", "./docs/postman_collection.json")

	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    CHAT MESSENGER SERVER v2.0                       ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  Server:   http://localhost:%s                                        \n", cfg.ServerPort)
	fmt.Printf("║  Frontend: http://localhost:%s/                                         \n", cfg.ServerPort)
	fmt.Printf("║  Swagger:  http://localhost:%s/swagger/index.html                     \n", cfg.ServerPort)
	fmt.Printf("║  Postman:  http://localhost:%s/postman                                \n", cfg.ServerPort)
	fmt.Printf("║  WebSocket: ws://localhost:%s/ws?token={jwt}                          \n", cfg.ServerPort)
	fmt.Printf("║  Health:   http://localhost:%s/health                                 \n", cfg.ServerPort)
	fmt.Println("╠══════════════════════════════════════════════════════════════════════╣")
	fmt.Println("║  Import Postman: Download from /postman, import into Postman        ║")
	fmt.Println("║  Auth: Register or Login → copy token → paste as Bearer token       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

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
