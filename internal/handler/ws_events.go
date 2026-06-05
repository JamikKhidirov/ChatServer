package handler

import (
	"encoding/json"

	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/internal/ws"

	"github.com/gin-gonic/gin"
)

type WebSocketEvents struct {
	hub           *ws.Hub
	chatService   *service.ChatService
	messageService *service.MessageService
	userService   *service.UserService
	pushService   *service.PushService
}

func NewWebSocketEvents(
	hub *ws.Hub,
	chatService *service.ChatService,
	messageService *service.MessageService,
	userService *service.UserService,
	pushService *service.PushService,
) *WebSocketEvents {
	return &WebSocketEvents{
		hub:           hub,
		chatService:   chatService,
		messageService: messageService,
		userService:   userService,
		pushService:   pushService,
	}
}

func (e *WebSocketEvents) WrapSendMessage(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		chatID := c.Param("id")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			go func() {
				chat, _ := e.chatService.GetChat(chatID, userID.(string))
				if chat != nil && chat.LastMessage != nil {
					var userIDs []string
					for _, p := range chat.Participants {
						if p.ID != userID.(string) {
							userIDs = append(userIDs, p.ID)
						}
					}

					e.hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
						Type:    ws.MsgNewMessage,
						Payload: chat.LastMessage,
					})

					e.pushService.SendMessageNotification(userID.(string), chatID, chat.LastMessage.ID, chat.LastMessage.Content, string(chat.LastMessage.Type))
				}
			}()
		}
	}
}

func (e *WebSocketEvents) WrapEditMessage(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		msgID := c.Param("id")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			go func() {
				msg, err := e.messageService.GetMessageByID(msgID, userID.(string))
				if err == nil && msg != nil {
					participants, _ := e.chatService.GetChat(msg.ChatID, userID.(string))
					if participants != nil {
						var userIDs []string
						for _, p := range participants.Participants {
							if p.ID != userID.(string) {
								userIDs = append(userIDs, p.ID)
							}
						}
						e.hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
							Type:    ws.MsgEditMessage,
							Payload: msg,
						})
					}
				}
			}()
		}
	}
}

func (e *WebSocketEvents) WrapDeleteMessage(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		msgID := c.Param("id")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			go func() {
				participants, _ := e.chatService.ListChats(userID.(string))
				for _, chat := range participants {
					var userIDs []string
					for _, p := range chat.Participants {
						if p.ID != userID.(string) {
							userIDs = append(userIDs, p.ID)
						}
					}
					e.hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
						Type: ws.MsgDeleteMessage,
						Payload: map[string]string{
							"messageId": msgID,
							"chatId":    chat.ID,
						},
					})
				}
			}()
		}
	}
}

func (e *WebSocketEvents) WrapCreateChat(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			var resp struct {
				Data *domain.ChatResponse `json:"data"`
			}
			if data, exists := c.Get("response"); exists {
				resp.Data = data.(*domain.ChatResponse)
			}

			go func() {
				if resp.Data != nil {
					var userIDs []string
					for _, p := range resp.Data.Participants {
						if p.ID != userID.(string) {
							userIDs = append(userIDs, p.ID)
						}
					}
					payload, _ := json.Marshal(resp.Data)
					_ = payload
					e.hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
						Type:    ws.MsgChatCreated,
						Payload: resp.Data,
					})
				}
			}()
		}
	}
}

func (e *WebSocketEvents) WrapDeleteChat(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatID := c.Param("id")
		userID, _ := c.Get("userID")

		chat, _ := e.chatService.GetChat(chatID, userID.(string))
		var userIDs []string
		if chat != nil {
			for _, p := range chat.Participants {
				userIDs = append(userIDs, p.ID)
			}
		}

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			go func() {
				e.hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
					Type: ws.MsgChatDeleted,
					Payload: map[string]string{
						"chatId": chatID,
					},
				})
			}()
		}
	}
}

func (e *WebSocketEvents) WrapInitiateCall(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			if resp, exists := c.Get("callResponse"); exists {
				if callResp, ok := resp.(*domain.Call); ok {
					go func() {
						chat, _ := e.chatService.GetChat(callResp.ChatID, userID.(string))
						if chat != nil {
							var userIDs []string
							for _, p := range chat.Participants {
								if p.ID != userID.(string) {
									userIDs = append(userIDs, p.ID)
								}
							}
							e.hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
								Type:    ws.MsgCallOffer,
								Payload: map[string]string{"chatId": callResp.ChatID, "callerId": userID.(string)},
							})
							e.pushService.SendCallNotification(userID.(string), callResp.ChatID, callResp.ID, "voice")
						}
					}()
				}
			}
		}
	}
}

func (e *WebSocketEvents) WrapRespondCall(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		callID := c.Param("id")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			go func() {
				call, err := e.chatService.GetChat(callID, userID.(string))
				if err == nil && call != nil {
					var userIDs []string
					for _, p := range call.Participants {
						if p.ID != userID.(string) {
							userIDs = append(userIDs, p.ID)
						}
					}

					var actionType ws.MessageType
					actionType = ws.MsgCallAccept

					e.hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
						Type:    actionType,
						Payload: map[string]string{"callId": callID, "userId": userID.(string)},
					})
				}
			}()
		}
	}
}

func (e *WebSocketEvents) WrapEndCall(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		callID := c.Param("id")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			go func() {
				chats, _ := e.chatService.ListChats(userID.(string))
				for _, chat := range chats {
					var userIDs []string
					for _, p := range chat.Participants {
						if p.ID != userID.(string) {
							userIDs = append(userIDs, p.ID)
						}
					}
					e.hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
						Type:    ws.MsgCallEnd,
						Payload: map[string]string{"callId": callID, "userId": userID.(string)},
					})
				}
			}()
		}
	}
}
