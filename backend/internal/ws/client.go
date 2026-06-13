package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 65536
)

type Client struct {
	Hub     *Hub
	Conn    *websocket.Conn
	Send    chan WSOutgoingMessage
	UserID  string
	done    chan struct{}
	onClose func()
}

func NewClient(hub *Hub, conn *websocket.Conn, userID string, onClose func()) *Client {
	if onClose == nil {
		onClose = func() {}
	}
	return &Client{
		Hub:     hub,
		Conn:    conn,
		Send:    make(chan WSOutgoingMessage, 256),
		UserID:  userID,
		done:    make(chan struct{}),
		onClose: onClose,
	}
}

func (c *Client) Start() {
	c.Hub.RegisterClient(c)
	go c.writePump()
	go c.readPump()
}

func (c *Client) Close() {
	close(c.done)
	c.Hub.UnregisterClient(c.UserID)
	c.Conn.Close()
	c.onClose()
}

func (c *Client) readPump() {
	defer c.Close()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var wsMsg WsMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Invalid WS message: %v", err)
			continue
		}

		c.handleMessage(wsMsg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteJSON(message)

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.done:
			return
		}
	}
}

func (c *Client) handleMessage(msg WsMessage) {
	dispatch := func(handler MessageActionHandler) {
		if handler != nil {
			handler(c.UserID, msg.Payload, c.Hub)
		}
	}

	switch msg.Type {
	case MsgKeyboardOpened:
		var payload struct {
			ChatID string `json:"chatId"`
		}
		json.Unmarshal(msg.Payload, &payload)
		if payload.ChatID != "" {
			c.Hub.BroadcastToChatExcept([]string{c.UserID}, WSOutgoingMessage{
				Type:    MsgKeyboardOpened,
				Payload: map[string]string{"chatId": payload.ChatID, "userId": c.UserID},
			})
		}

	case MsgKeyboardClosed:
		var payload struct {
			ChatID string `json:"chatId"`
		}
		json.Unmarshal(msg.Payload, &payload)
		if payload.ChatID != "" {
			c.Hub.BroadcastToChatExcept([]string{c.UserID}, WSOutgoingMessage{
				Type:    MsgKeyboardClosed,
				Payload: map[string]string{"chatId": payload.ChatID, "userId": c.UserID},
			})
		}

	case MsgTyping:
		var payload struct {
			ChatID string `json:"chatId"`
		}
		json.Unmarshal(msg.Payload, &payload)
		if payload.ChatID != "" {
			c.Hub.BroadcastToChatExcept([]string{c.UserID}, WSOutgoingMessage{
				Type:    MsgTyping,
				Payload: map[string]string{"chatId": payload.ChatID, "userId": c.UserID},
			})
		}

	case MsgStopTyping:
		var payload struct {
			ChatID string `json:"chatId"`
		}
		json.Unmarshal(msg.Payload, &payload)
		if payload.ChatID != "" {
			c.Hub.BroadcastToChatExcept([]string{c.UserID}, WSOutgoingMessage{
				Type:    MsgStopTyping,
				Payload: map[string]string{"chatId": payload.ChatID, "userId": c.UserID},
			})
		}

	// WebRTC signaling relay
	case MsgCallOffer, MsgCallAnswer, MsgCallICE:
		var payload struct {
			ChatID  string `json:"chatId"`
			CallID  string `json:"callId"`
			SDP     string `json:"sdp,omitempty"`
			Candidate string `json:"candidate,omitempty"`
		}
		json.Unmarshal(msg.Payload, &payload)
		if payload.ChatID != "" {
			c.Hub.BroadcastToChatExcept([]string{c.UserID}, WSOutgoingMessage{
				Type:    MsgCallOffer,
				Payload: msg.Payload,
			})
		} else if payload.CallID != "" {
			c.Hub.BroadcastToChatExcept([]string{c.UserID}, WSOutgoingMessage{
				Type:    msg.Type,
				Payload: msg.Payload,
			})
		}

	case MsgSendMessage:
		dispatch(c.Hub.OnSendMessage)
	case MsgEditMessageReq:
		dispatch(c.Hub.OnEditMessage)
	case MsgDeleteMessageReq:
		dispatch(c.Hub.OnDeleteMessage)
	case MsgReadMessageReq:
		dispatch(c.Hub.OnReadMessage)
	case MsgAddReaction:
		dispatch(c.Hub.OnAddReaction)
	case MsgRemoveReaction:
		dispatch(c.Hub.OnRemoveReaction)
	case MsgTogglePin:
		dispatch(c.Hub.OnTogglePin)
	case MsgStarMessage:
		dispatch(c.Hub.OnStarMessage)
	case MsgUnstarMessage:
		dispatch(c.Hub.OnUnstarMessage)
	case MsgForwardMessage:
		dispatch(c.Hub.OnForwardMessage)
	case MsgCreateChat:
		dispatch(c.Hub.OnCreateChat)
	case MsgUpdateChat:
		dispatch(c.Hub.OnUpdateChat)
	case MsgAddParticipant:
		dispatch(c.Hub.OnAddParticipant)
	case MsgRemoveParticipant:
		dispatch(c.Hub.OnRemoveParticipant)
	case MsgLeaveChat:
		dispatch(c.Hub.OnLeaveChat)
	case MsgPinChat:
		dispatch(c.Hub.OnPinChat)
	case MsgUnpinChat:
		dispatch(c.Hub.OnUnpinChat)
	case MsgArchiveChat:
		dispatch(c.Hub.OnArchiveChat)
	case MsgUnarchiveChat:
		dispatch(c.Hub.OnUnarchiveChat)
	case MsgBlockUser:
		dispatch(c.Hub.OnBlockUser)
	case MsgUnblockUser:
		dispatch(c.Hub.OnUnblockUser)

	case MsgCallReject:
		var payload struct {
			CallID string `json:"callId"`
			ChatID string `json:"chatId"`
		}
		json.Unmarshal(msg.Payload, &payload)
		c.Hub.BroadcastToChatExcept([]string{c.UserID}, WSOutgoingMessage{
			Type:    MsgCallReject,
			Payload: map[string]string{"callId": payload.CallID, "userId": c.UserID},
		})
	}
}


