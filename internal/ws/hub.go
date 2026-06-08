package ws

import (
	"encoding/json"
	"log"
	"sync"
)

type MessageType string

const (
	MsgNewMessage    MessageType = "message:new"
	MsgEditMessage   MessageType = "message:edited"
	MsgDeleteMessage MessageType = "message:deleted"
	MsgReadMessage   MessageType = "message:read"
	MsgTyping        MessageType = "user:typing"
	MsgStopTyping    MessageType = "user:stop_typing"
	MsgOnline        MessageType = "user:online"
	MsgOffline       MessageType = "user:offline"
	MsgChatCreated   MessageType = "chat:created"
	MsgChatUpdated   MessageType = "chat:updated"
	MsgChatDeleted   MessageType = "chat:deleted"
	MsgCallOffer     MessageType = "call:offer"
	MsgCallAnswer    MessageType = "call:answer"
	MsgCallICE       MessageType = "call:ice"
	MsgCallEnd       MessageType = "call:end"
	MsgCallMissed    MessageType = "call:missed"
	MsgCallReject    MessageType = "call:reject"
	MsgCallAccept    MessageType = "call:accept"
)

type WSOutgoingMessage struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

type WsMessage struct {
	Type    MessageType      `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Hub struct {
	mu       sync.RWMutex
	clients  map[string]*Client
	register chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:  make(map[string]*Client),
		register: make(chan *Client, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("Client connected: %s", client.UserID)
		}
	}
}

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(userID string) {
	h.mu.Lock()
	delete(h.clients, userID)
	h.mu.Unlock()
	log.Printf("Client disconnected: %s", userID)
}

func (h *Hub) SendToUser(userID string, msg WSOutgoingMessage) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()
	if ok {
		select {
		case client.Send <- msg:
		default:
		}
	}
}

func (h *Hub) SendToUsers(userIDs []string, msg WSOutgoingMessage) {
	for _, uid := range userIDs {
		h.SendToUser(uid, msg)
	}
}

func (h *Hub) IsOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

func (h *Hub) BroadcastToChat(chatParticipants []string, msg WSOutgoingMessage) {
	h.SendToUsers(chatParticipants, msg)
}

func (h *Hub) BroadcastToChatExcept(excludeIDs []string, msg WSOutgoingMessage) {
	exclude := make(map[string]bool)
	for _, id := range excludeIDs {
		exclude[id] = true
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for userID, client := range h.clients {
		if !exclude[userID] {
			select {
			case client.Send <- msg:
			default:
			}
		}
	}
}

func (h *Hub) GetConnectedUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	users := make([]string, 0, len(h.clients))
	for uid := range h.clients {
		users = append(users, uid)
	}
	return users
}
