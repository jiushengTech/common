package websocket

import (
	"base-service/internal/global"
	"context"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
	"strings"
)

var channelBufSize = 256

type SessionID string

type Session struct {
	id     SessionID
	conn   *ws.Conn
	send   chan []byte
	server *Server
}

func NewSession(conn *ws.Conn, server *Server, path string) *Session {
	if conn == nil {
		panic("conn cannot be nil")
	}

	uuId, _ := uuid.NewUUID()
	sessionID := path + "-" + uuId.String()
	//global.REDIS.
	err := global.REDIS.Set(context.Background(), path, sessionID, 0).Err()
	if err != nil {
		global.LOG.Errorf("redis setEx error: %v", err)
	}
	c := &Session{
		id:     SessionID(sessionID),
		conn:   conn,
		send:   make(chan []byte, channelBufSize),
		server: server,
	}

	return c
}

func (c *Session) Conn() *ws.Conn {
	return c.conn
}

func (c *Session) SessionID() SessionID {
	return c.id
}

func (c *Session) SendMessage(message []byte) {
	select {
	case c.send <- message:
	}
}

func (c *Session) Close() {
	str := string(c.SessionID())
	var key string
	// 使用 strings.Index 函数查找第一个 "-" 的索引
	index := strings.Index(str, "-")
	if index == -1 {
		// 如果字符串中不存在 "-"，则直接使用原始字符串
		key = str
	} else {
		// 使用 strings.SplitN 函数将字符串按照 "-" 分割成多个子字符串
		key = strings.SplitN(str, "-", 2)[0]
	}
	global.REDIS.Del(context.Background(), key)
	c.server.unregister <- c
	c.closeConnect()
}

func (c *Session) Listen() {
	go c.writePump()
	go c.readPump()
}

func (c *Session) closeConnect() {
	//LogInfo(c.SessionID(), " connection closed")
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			global.LOG.WithContext(context.Background()).Errorf("disconnect error: %s", err.Error())
		}
		c.conn = nil
	}
}

func (c *Session) sendPingMessage(message string) error {
	if c.conn == nil {
		return nil
	}
	return c.conn.WriteMessage(ws.PingMessage, []byte(message))
}

func (c *Session) sendPongMessage(message string) error {
	if c.conn == nil {
		return nil
	}
	return c.conn.WriteMessage(ws.PongMessage, []byte(message))
}

func (c *Session) sendTextMessage(message string) error {
	if c.conn == nil {
		return nil
	}
	return c.conn.WriteMessage(ws.TextMessage, []byte(message))
}

func (c *Session) sendBinaryMessage(message []byte) error {
	if c.conn == nil {
		return nil
	}
	return c.conn.WriteMessage(ws.BinaryMessage, message)
}

func (c *Session) writePump() {
	defer c.Close()

	for {
		select {
		case msg := <-c.send:
			var err error
			switch c.server.payloadType {
			case PayloadTypeBinary:
				if err = c.sendBinaryMessage(msg); err != nil {
					global.LOG.WithContext(context.Background()).Error("write binary message error: ", err)
					return
				}
				break

			case PayloadTypeText:
				if err = c.sendTextMessage(string(msg)); err != nil {
					global.LOG.WithContext(context.Background()).Error("write text message error: ", err)
					return
				}
				break
			}

		}
	}
}

func (c *Session) readPump() {
	defer c.Close()

	for {
		if c.conn == nil {
			break
		}

		messageType, data, err := c.conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseNormalClosure, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				global.LOG.WithContext(context.Background()).Errorf("read message error: %v", err)
			}
			return
		}

		switch messageType {
		case ws.CloseMessage:
			return

		case ws.BinaryMessage:
			_ = c.server.messageHandler(c.SessionID(), data)
			break

		case ws.TextMessage:
			_ = c.server.messageHandler(c.SessionID(), data)
			break

		case ws.PingMessage:
			if err = c.sendPongMessage(""); err != nil {
				global.LOG.WithContext(context.Background()).Error("write pong message error: ", err)
				return
			}
			break

		case ws.PongMessage:
			break
		}

	}
}
