package socket

import (
	"context"
	"net/url"
	"sync"
)

// TCPServer TCP服务器实现
type TCPServer struct {
	mu         sync.RWMutex
	server     *Server
	connPool   *ConnectionPool
	closed     bool
	closedChan chan struct{}
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(opts ...Option) *TCPServer {
	server := NewServer(opts...)

	return &TCPServer{
		server:     server,
		connPool:   NewConnectionPool(int(server.MaxConns)),
		closedChan: make(chan struct{}),
	}
}

// Start 启动TCP服务器
func (s *TCPServer) Start(ctx context.Context) error {
	return s.server.Start(ctx)
}

// Stop 停止TCP服务器
func (s *TCPServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	close(s.closedChan)

	// 释放连接池
	s.connPool.Close()

	return s.server.Stop(ctx)
}

// Endpoint 返回服务器端点
func (s *TCPServer) Endpoint() (*url.URL, error) {
	return s.server.Endpoint()
}

// SendTo 向指定目标发送数据
func (s *TCPServer) SendTo(targetAddr string, data []byte) (int, error) {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return 0, ErrServerClosed
	}
	s.mu.RUnlock()

	return s.server.SendTo(targetAddr, data)
}

// Broadcast 向所有目标广播数据
func (s *TCPServer) Broadcast(data []byte) (int, error) {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return 0, ErrServerClosed
	}
	s.mu.RUnlock()

	return s.server.Broadcast(data)
}
