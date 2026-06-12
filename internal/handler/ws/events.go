package wshandler

import (
	"encoding/json"

	calldomain "ChatServerGolang/internal/domain/call"
	chatdomain "ChatServerGolang/internal/domain/chat"
	messagedomain "ChatServerGolang/internal/domain/message"
	"ChatServerGolang/internal/service"
		"ChatServerGolang/internal/ws"

	"github.com/gin-gonic/gin"
)

type WebSocketEvents struct {
	hub            *ws.Hub
	chatService    service.ChatService
	messageService service.MessageService
	userService    service.UserService
	pushService    service.PushService
	callService    service.CallService
}

func NewWebSocketEvents(
	hub *ws.Hub,
	chatService service.ChatService,
	messageService service.MessageService,
	userService service.UserService,
	pushService service.PushService,
	callService service.CallService,
) *WebSocketEvents {
	e := &WebSocketEvents{
		hub:            hub,
		chatService:    chatService,
		messageService: messageService,
		userService:    userService,
		pushService:    pushService,
		callService:    callService,
	}

	hub.OnSendMessage = e.handleSendMessage
	hub.OnEditMessage = e.handleEditMessage
	hub.OnDeleteMessage = e.handleDeleteMessage
	hub.OnReadMessage = e.handleReadMessage
	hub.OnAddReaction = e.handleAddReaction
	hub.OnRemoveReaction = e.handleRemoveReaction
	hub.OnTogglePin = e.handleTogglePin
	hub.OnStarMessage = e.handleStarMessage
	hub.OnUnstarMessage = e.handleUnstarMessage
	hub.OnForwardMessage = e.handleForwardMessage
	hub.OnCreateChat = e.handleCreateChat
	hub.OnUpdateChat = e.handleUpdateChat
	hub.OnAddParticipant = e.handleAddParticipant
	hub.OnRemoveParticipant = e.handleRemoveParticipant
	hub.OnLeaveChat = e.handleLeaveChat
	hub.OnPinChat = e.handlePinChat
	hub.OnUnpinChat = e.handleUnpinChat
	hub.OnArchiveChat = e.handleArchiveChat
	hub.OnUnarchiveChat = e.handleUnarchiveChat
	hub.OnBlockUser = e.handleBlockUser
	hub.OnUnblockUser = e.handleUnblockUser

	return e
}

func (e *WebSocketEvents) WrapSendMessage(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		chatID := c.Param("id")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
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
		}
	}
}

func (e *WebSocketEvents) WrapEditMessage(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		msgID := c.Param("id")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			msg, err := e.messageService.GetMessageByID(msgID, userID.(string))
			if err == nil && msg != nil {
				chat, _ := e.chatService.GetChat(msg.ChatID, userID.(string))
				if chat != nil {
					var userIDs []string
					for _, p := range chat.Participants {
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
		}
	}
}

func (e *WebSocketEvents) WrapDeleteMessage(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		msgID := c.Param("id")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			chats, _ := e.chatService.ListChats(userID.(string))
			for _, chat := range chats {
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
		}
	}
}

func (e *WebSocketEvents) WrapCreateChat(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			if resp, exists := c.Get("chatResponse"); exists {
				chatResp := resp.(*chatdomain.ChatResponse)
				var userIDs []string
				for _, p := range chatResp.Participants {
					if p.ID != userID.(string) {
						userIDs = append(userIDs, p.ID)
					}
				}
				e.hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
					Type:    ws.MsgChatCreated,
					Payload: chatResp,
				})
			}
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
			e.hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
				Type: ws.MsgChatDeleted,
				Payload: map[string]string{
					"chatId": chatID,
				},
			})
		}
	}
}

func (e *WebSocketEvents) WrapInitiateCall(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			if resp, exists := c.Get("callResponse"); exists {
				if callResp, ok := resp.(*calldomain.Call); ok {
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
						Payload: map[string]string{"chatId": callResp.ChatID, "callerId": userID.(string), "type": string(callResp.Type)},
					})
						e.pushService.SendCallNotification(userID.(string), callResp.ChatID, callResp.ID, "voice")
					}
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
			call, err := e.callService.GetCallByID(callID)
			if err == nil && call != nil {
				chat, _ := e.chatService.GetChat(call.ChatID, userID.(string))
				if chat != nil {
					var userIDs []string
					for _, p := range chat.Participants {
						if p.ID != userID.(string) {
							userIDs = append(userIDs, p.ID)
						}
					}

					e.hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
						Type:    ws.MsgCallAccept,
						Payload: map[string]string{"callId": callID, "userId": userID.(string)},
					})
				}
			}
		}
	}
}

func (e *WebSocketEvents) broadcastToChatParticipants(chatID, excludeUserID string, msg ws.WSOutgoingMessage) {
	chat, _ := e.chatService.GetChat(chatID, excludeUserID)
	if chat == nil {
		return
	}
	var userIDs []string
	for _, p := range chat.Participants {
		if p.ID != excludeUserID {
			userIDs = append(userIDs, p.ID)
		}
	}
	e.hub.BroadcastToChat(userIDs, msg)
}

func (e *WebSocketEvents) broadcastToAllParticipants(chatID, senderID string, msg ws.WSOutgoingMessage) {
	chat, _ := e.chatService.GetChat(chatID, senderID)
	if chat == nil {
		return
	}
	var userIDs []string
	for _, p := range chat.Participants {
		userIDs = append(userIDs, p.ID)
	}
	e.hub.BroadcastToChat(userIDs, msg)
}

func (e *WebSocketEvents) handleSendMessage(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		ChatID    string `json:"chatId"`
		Content   string `json:"content"`
		Type      messagedomain.MessageType `json:"type"`
		ReplyToID *string `json:"replyToId"`
	}
	json.Unmarshal(payload, &req)
	if req.ChatID == "" || req.Content == "" {
		return
	}

	sendReq := &messagedomain.SendMessageRequest{
		Content:   req.Content,
		Type:      req.Type,
		ReplyToID: req.ReplyToID,
	}
	msg, err := e.messageService.SendMessage(req.ChatID, userID, sendReq)
	if err != nil {
		return
	}

	chat, _ := e.chatService.GetChat(req.ChatID, userID)
	if chat != nil {
		var userIDs []string
		for _, p := range chat.Participants {
			userIDs = append(userIDs, p.ID)
		}

		hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
			Type:    ws.MsgNewMessage,
			Payload: msg,
		})

		if msg != nil {
			e.pushService.SendMessageNotification(userID, req.ChatID, msg.ID, msg.Content, string(msg.Type))
		}
	}
}

func (e *WebSocketEvents) handleEditMessage(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		MessageID string `json:"messageId"`
		Content   string `json:"content"`
	}
	json.Unmarshal(payload, &req)
	if req.MessageID == "" || req.Content == "" {
		return
	}
	editReq := &messagedomain.EditMessageRequest{Content: req.Content}
	msg, err := e.messageService.EditMessage(req.MessageID, userID, editReq)
	if err != nil {
		return
	}
	e.broadcastToChatParticipants(msg.ChatID, userID, ws.WSOutgoingMessage{
		Type:    ws.MsgEditMessage,
		Payload: msg,
	})
}

func (e *WebSocketEvents) handleDeleteMessage(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		MessageID string `json:"messageId"`
		ChatID    string `json:"chatId"`
	}
	json.Unmarshal(payload, &req)
	if req.MessageID == "" {
		return
	}
	err := e.messageService.DeleteMessage(req.MessageID, userID)
	if err != nil {
		return
	}
	e.broadcastToChatParticipants(req.ChatID, userID, ws.WSOutgoingMessage{
		Type: ws.MsgDeleteMessage,
		Payload: map[string]string{
			"messageId": req.MessageID,
			"chatId":    req.ChatID,
		},
	})
}

func (e *WebSocketEvents) handleReadMessage(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		MessageID string `json:"messageId"`
		ChatID    string `json:"chatId"`
	}
	json.Unmarshal(payload, &req)
	if req.MessageID == "" || req.ChatID == "" {
		return
	}
	err := e.messageService.MarkMessageRead(req.MessageID, userID)
	if err != nil {
		return
	}
	e.broadcastToChatParticipants(req.ChatID, userID, ws.WSOutgoingMessage{
		Type: ws.MsgReadMessage,
		Payload: map[string]string{
			"messageId": req.MessageID,
			"userId":    userID,
			"chatId":    req.ChatID,
		},
	})
}

func (e *WebSocketEvents) handleAddReaction(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		MessageID string `json:"messageId"`
		Emoji     string `json:"emoji"`
	}
	json.Unmarshal(payload, &req)
	if req.MessageID == "" || req.Emoji == "" {
		return
	}
	msg, err := e.messageService.AddReaction(req.MessageID, userID, req.Emoji)
	if err != nil {
		return
	}
	e.broadcastToChatParticipants(msg.ChatID, userID, ws.WSOutgoingMessage{
		Type:    ws.MsgReaction,
		Payload: msg,
	})
}

func (e *WebSocketEvents) handleRemoveReaction(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		MessageID string `json:"messageId"`
		Emoji     string `json:"emoji"`
	}
	json.Unmarshal(payload, &req)
	if req.MessageID == "" || req.Emoji == "" {
		return
	}
	msg, err := e.messageService.RemoveReaction(req.MessageID, userID, req.Emoji)
	if err != nil {
		return
	}
	e.broadcastToChatParticipants(msg.ChatID, userID, ws.WSOutgoingMessage{
		Type:    ws.MsgReaction,
		Payload: msg,
	})
}

func (e *WebSocketEvents) handleTogglePin(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		MessageID string `json:"messageId"`
		Pin       bool   `json:"pin"`
	}
	json.Unmarshal(payload, &req)
	if req.MessageID == "" {
		return
	}
	msg, err := e.messageService.TogglePin(req.MessageID, userID, req.Pin)
	if err != nil {
		return
	}
	e.broadcastToChatParticipants(msg.ChatID, userID, ws.WSOutgoingMessage{
		Type:    ws.MsgPinned,
		Payload: msg,
	})
}

func (e *WebSocketEvents) handleStarMessage(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		MessageID string `json:"messageId"`
	}
	json.Unmarshal(payload, &req)
	if req.MessageID == "" {
		return
	}
	e.messageService.StarMessage(req.MessageID, userID)
}

func (e *WebSocketEvents) handleUnstarMessage(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		MessageID string `json:"messageId"`
	}
	json.Unmarshal(payload, &req)
	if req.MessageID == "" {
		return
	}
	e.messageService.UnstarMessage(req.MessageID, userID)
}

func (e *WebSocketEvents) handleForwardMessage(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		MessageID string `json:"messageId"`
		ToChatID  string `json:"toChatId"`
	}
	json.Unmarshal(payload, &req)
	if req.MessageID == "" || req.ToChatID == "" {
		return
	}
	msg, err := e.messageService.ForwardMessage(req.MessageID, "", req.ToChatID, userID)
	if err != nil {
		return
	}
	e.broadcastToAllParticipants(req.ToChatID, userID, ws.WSOutgoingMessage{
		Type:    ws.MsgNewMessage,
		Payload: msg,
	})
}

func (e *WebSocketEvents) handleCreateChat(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req chatdomain.CreateChatRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		return
	}
	chat, err := e.chatService.CreateChat(userID, &req)
	if err != nil {
		return
	}
	var userIDs []string
	for _, p := range chat.Participants {
		if p.ID != userID {
			userIDs = append(userIDs, p.ID)
		}
	}
	hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
		Type:    ws.MsgChatCreated,
		Payload: chat,
	})
}

func (e *WebSocketEvents) handleUpdateChat(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		ChatID      string `json:"chatId"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		AvatarURL   string `json:"avatarUrl,omitempty"`
	}
	json.Unmarshal(payload, &req)
	if req.ChatID == "" {
		return
	}
	updateReq := &chatdomain.UpdateGroupRequest{
		Name:        req.Name,
		Description: req.Description,
		AvatarURL:   req.AvatarURL,
	}
	if err := e.chatService.UpdateGroup(req.ChatID, userID, updateReq); err != nil {
		return
	}
	chat, _ := e.chatService.GetChat(req.ChatID, userID)
	if chat != nil {
		hub.BroadcastToChat([]string{}, ws.WSOutgoingMessage{
			Type:    ws.MsgChatUpdated,
			Payload: chat,
		})
	}
}

func (e *WebSocketEvents) handleAddParticipant(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		ChatID      string `json:"chatId"`
		ParticipantID string `json:"userId"`
	}
	json.Unmarshal(payload, &req)
	if req.ChatID == "" || req.ParticipantID == "" {
		return
	}
	if err := e.chatService.AddParticipant(req.ChatID, req.ParticipantID, userID); err != nil {
		return
	}
	chat, _ := e.chatService.GetChat(req.ChatID, userID)
	if chat != nil {
		hub.BroadcastToChat([]string{}, ws.WSOutgoingMessage{
			Type:    ws.MsgChatUpdated,
			Payload: chat,
		})
	}
}

func (e *WebSocketEvents) handleRemoveParticipant(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		ChatID        string `json:"chatId"`
		TargetUserID  string `json:"userId"`
	}
	json.Unmarshal(payload, &req)
	if req.ChatID == "" || req.TargetUserID == "" {
		return
	}
	if err := e.chatService.RemoveParticipant(req.ChatID, req.TargetUserID, userID); err != nil {
		return
	}
	chat, _ := e.chatService.GetChat(req.ChatID, userID)
	if chat != nil {
		userIDs := []string{req.TargetUserID}
		for _, p := range chat.Participants {
			if p.ID != userID {
				userIDs = append(userIDs, p.ID)
			}
		}
		hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
			Type:    ws.MsgChatUpdated,
			Payload: chat,
		})
	}
}

func (e *WebSocketEvents) handleLeaveChat(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		ChatID string `json:"chatId"`
	}
	json.Unmarshal(payload, &req)
	if req.ChatID == "" {
		return
	}
	if err := e.chatService.LeaveGroup(req.ChatID, userID); err != nil {
		return
	}
	chat, _ := e.chatService.GetChat(req.ChatID, userID)
	if chat != nil {
		userIDs := []string{userID}
		for _, p := range chat.Participants {
			userIDs = append(userIDs, p.ID)
		}
		hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
			Type:    ws.MsgChatUpdated,
			Payload: chat,
		})
	}
}

func (e *WebSocketEvents) handlePinChat(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		ChatID string `json:"chatId"`
	}
	json.Unmarshal(payload, &req)
	if req.ChatID == "" {
		return
	}
	e.chatService.PinChat(req.ChatID, userID)
}

func (e *WebSocketEvents) handleUnpinChat(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		ChatID string `json:"chatId"`
	}
	json.Unmarshal(payload, &req)
	if req.ChatID == "" {
		return
	}
	e.chatService.UnpinChat(req.ChatID, userID)
}

func (e *WebSocketEvents) handleArchiveChat(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		ChatID string `json:"chatId"`
	}
	json.Unmarshal(payload, &req)
	if req.ChatID == "" {
		return
	}
	e.chatService.ArchiveChat(req.ChatID, userID)
}

func (e *WebSocketEvents) handleUnarchiveChat(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		ChatID string `json:"chatId"`
	}
	json.Unmarshal(payload, &req)
	if req.ChatID == "" {
		return
	}
	e.chatService.UnarchiveChat(req.ChatID, userID)
}

func (e *WebSocketEvents) handleBlockUser(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		UserID string `json:"userId"`
	}
	json.Unmarshal(payload, &req)
	if req.UserID == "" {
		return
	}
	e.userService.BlockUser(userID, req.UserID)
}

func (e *WebSocketEvents) handleUnblockUser(userID string, payload json.RawMessage, hub *ws.Hub) {
	var req struct {
		UserID string `json:"userId"`
	}
	json.Unmarshal(payload, &req)
	if req.UserID == "" {
		return
	}
	e.userService.UnblockUser(userID, req.UserID)
}

func (e *WebSocketEvents) WrapEndCall(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		callID := c.Param("id")

		handler(c)

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			call, err := e.callService.GetCallByID(callID)
			if err == nil && call != nil {
				chat, _ := e.chatService.GetChat(call.ChatID, userID.(string))
				if chat != nil {
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
			}
		}
	}
}

