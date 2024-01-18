package websocketv2

import "github.com/gorilla/websocket"

// Option 是 WebSocket 客户端选项类型
type Option func(o *WebSocketServer)

// WithURL 设置 WebSocket 连接的 url
func WithURL(u string) Option {
	return func(c *WebSocketServer) {
		c.url = u
	}
}

// WithHeader 设置 WebSocket 连接的请求头
func WithUpgrader(upgrader *websocket.Upgrader) Option {
	return func(c *WebSocketServer) {
		c.upgrader = upgrader
	}
}
