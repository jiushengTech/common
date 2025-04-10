package socket

import (
	"context"
	"net/url"
	"sync"
)

// UDPServer UDP服务器实现
type UDPServer struct {
	mu         sync.RWMutex
	server     *Server
	closed     bool
	closedChan chan struct{}
}

// NewUDPServer 创建UDP服务器
func NewUDPServer(opts ...Option) *UDPServer {
	// 确保网络类型是UDP
	server := NewServer(append(opts, WithNetwork("udp"))...)

	return &UDPServer{
		server:     server,
		closedChan: make(chan struct{}),
	}
}

// Start 启动UDP服务器
func (s *UDPServer) Start(ctx context.Context) error {
	return s.server.Start(ctx)
}

// Stop 停止UDP服务器
func (s *UDPServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	close(s.closedChan)

	return s.server.Stop(ctx)
}

// Endpoint 返回服务器端点
func (s *UDPServer) Endpoint() (*url.URL, error) {
	return s.server.Endpoint()
}

// SendTo 向指定目标发送UDP数据
func (s *UDPServer) SendTo(targetAddr string, data []byte) (int, error) {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return 0, ErrServerClosed
	}
	s.mu.RUnlock()

	return s.server.SendTo(targetAddr, data)
}

// Broadcast 向所有目标广播UDP数据
func (s *UDPServer) Broadcast(data []byte) (int, error) {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return 0, ErrServerClosed
	}
	s.mu.RUnlock()

	return s.server.Broadcast(data)
}
