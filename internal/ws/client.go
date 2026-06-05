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
	switch msg.Type {
	case MsgTyping:
		var payload struct {
			ChatID string `json:"chatId"`
		}
		json.Unmarshal(msg.Payload, &payload)
		if payload.ChatID != "" {
			c.Hub.SendToUser(c.UserID, WSOutgoingMessage{
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
			c.Hub.SendToUser(c.UserID, WSOutgoingMessage{
				Type:    MsgStopTyping,
				Payload: map[string]string{"chatId": payload.ChatID, "userId": c.UserID},
			})
		}
	}
}
