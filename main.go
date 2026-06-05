package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	callRepo := repository.NewCallRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	callService := service.NewCallService(callRepo, chatRepo, userRepo, userService)
	messageService := service.NewMessageService(messageRepo, chatRepo, userRepo, userService)
	chatService := service.NewChatService(chatRepo, userRepo, messageRepo, userService)
	pushService := service.NewPushService(userRepo, cfg)

	hub := ws.NewHub()
	go hub.Run()

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, pushService)
	chatHandler := handler.NewChatHandler(chatService)
	messageHandler := handler.NewMessageHandler(messageService)
	callHandler := handler.NewCallHandler(callService)
	wsHandler := handler.NewWSHandler(hub, authService, userRepo, chatRepo)
	wsEvents := handler.NewWebSocketEvents(hub, chatService, messageService, userService, pushService)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware(authService))
		{
			authorized.GET("/users/profile", userHandler.GetProfile)
			authorized.PUT("/users/profile", userHandler.UpdateProfile)
			authorized.PUT("/users/push-token", userHandler.UpdatePushToken)
			authorized.POST("/users/push-test", userHandler.TestPush)
			authorized.PUT("/users/status", userHandler.UpdateStatus)
			authorized.GET("/users/search", userHandler.SearchUsers)
			authorized.GET("/users/:id", userHandler.GetUserByID)
			authorized.GET("/users/username/:username", userHandler.GetUserByUsername)

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

			authorized.GET("/chats/:id/messages", messageHandler.GetMessages)
			authorized.POST("/chats/:id/messages", wsEvents.WrapSendMessage(messageHandler.SendMessage))
			authorized.PUT("/messages/:id", wsEvents.WrapEditMessage(messageHandler.EditMessage))
			authorized.DELETE("/messages/:id", wsEvents.WrapDeleteMessage(messageHandler.DeleteMessage))

			authorized.POST("/calls/initiate", wsEvents.WrapInitiateCall(callHandler.InitiateCall))
			authorized.POST("/calls/:id/respond", wsEvents.WrapRespondCall(callHandler.RespondCall))
			authorized.POST("/calls/:id/end", wsEvents.WrapEndCall(callHandler.EndCall))
			authorized.GET("/calls/:id", callHandler.GetCall)
			authorized.GET("/calls/history/:chatId", callHandler.GetCallHistory)
		}

		r.GET("/ws", wsHandler.HandleWebSocket)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
