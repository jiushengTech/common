package socket

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport"
	"net"
	"net/url"
	"sync"
)

// TCPServer TCP服务器实现
type TCPServer struct {
	mu          sync.RWMutex
	server      *Server
	closed      bool
	closedChan  chan struct{}
	tcpListener *net.TCPListener
	isStarted   bool
}

var (
	_ transport.Server     = (*TCPServer)(nil)
	_ transport.Endpointer = (*TCPServer)(nil)
)

func (s *TCPServer) GetTcpListener() *net.TCPListener {
	for {
		if s.isStarted {
			return s.tcpListener
		}
	}
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(opts ...Option) *TCPServer {
	server := NewServer(opts...)
	return &TCPServer{
		server:     server,
		closedChan: make(chan struct{}),
	}
}

// Start 启动TCP服务器
func (s *TCPServer) Start(ctx context.Context) error {
	addr, err := net.ResolveTCPAddr(s.server.network, s.server.address)
	if err != nil {
		return fmt.Errorf("解析 TCP 地址失败: %w", err)
	}
	tcp, err := net.ListenTCP(s.server.network, addr)
	if err != nil {
		return fmt.Errorf("TCP 监听失败: %w", err)
	}
	s.tcpListener = tcp
	s.isStarted = true
	return nil
}

// Stop 停止TCP服务器
func (s *TCPServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		fmt.Println("TCPServer 已经关闭，无需重复 Stop")
		return nil
	}

	s.closed = true
	close(s.closedChan)

	// 释放连接池
	s.server.connPool.Close()
	if s.tcpListener != nil {
		if err := s.tcpListener.Close(); err != nil {
			return err
		}
		s.tcpListener = nil
	}
	return nil
}

// Endpoint 返回服务器端点
func (s *TCPServer) Endpoint() (*url.URL, error) {
	if s.tcpListener != nil {
		addr := s.tcpListener.Addr().String()
		return url.Parse("socket://" + addr)
	}
	// fallback（启动前调用 Endpoint）
	return url.Parse("socket://" + s.server.address)
}

// SendTo 向指定目标发送数据
func (s *TCPServer) SendTo(targetAddr string, data []byte) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, ErrServerClosed
	}
	return s.server.SendTo(targetAddr, data)
}

// Broadcast 向所有目标广播数据
func (s *TCPServer) Broadcast(data []byte) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, ErrServerClosed
	}
	return s.server.Broadcast(data)
}
