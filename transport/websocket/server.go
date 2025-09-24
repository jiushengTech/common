package websocket

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/gorilla/websocket"
)

var (
	_ transport.Server     = (*WebSocketServer)(nil)
	_ transport.Endpointer = (*WebSocketServer)(nil)
)

// WebSocketServer 是 WebSocket 服务器结构体
type WebSocketServer struct {
	addr     string
	conn     *websocket.Conn
	url      string
	upgrader *websocket.Upgrader
}

// NewWebSocketServer 创建新的 WebSocket 服务器
func NewWebSocketServer(options ...Option) *WebSocketServer {
	server := &WebSocketServer{
		addr: "0.0.0.0",
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	// 应用选项
	for _, option := range options {
		option(server)
	}
	return server
}

// Stop 关闭 WebSocket 连接
func (s *WebSocketServer) Stop(ctx context.Context) error {
	return s.conn.Close()
}

// Start 启动 WebSocket 服务器
func (s *WebSocketServer) Start(ctx context.Context) error {
	http.HandleFunc(s.url, s.wsHandler)
	return http.ListenAndServe(s.addr, nil)
}

// Endpoint 返回 WebSocket 服务器的地址
func (s *WebSocketServer) Endpoint() (*url.URL, error) {
	prefix := "websocket://"
	addr := prefix + s.addr
	endpoint, err := url.Parse(addr)
	return endpoint, err
}

// WsHandler 处理 WebSocket 连接
func (s *WebSocketServer) wsHandler(res http.ResponseWriter, req *http.Request) {
	// 连接 WebSocket
	conn, err := s.upgrader.Upgrade(res, req, nil)
	if err != nil {
		return
	}
	s.conn = conn
}
