package socket

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport"
	"net"
	"net/url"
	"sync"
	"time"
)

// UDPServer UDP服务器实现
type UDPServer struct {
	mu         sync.RWMutex
	server     *Server
	closed     bool
	closedChan chan struct{}
	udpConn    *net.UDPConn
}

var (
	_ transport.Server     = (*UDPServer)(nil)
	_ transport.Endpointer = (*UDPServer)(nil)
)

// NewUDPServer 创建UDP服务器
func NewUDPServer(opts ...Option) *UDPServer {
	// 确保网络类型是UDP
	server := NewServer(append(opts, WithNetwork("udp"))...)

	return &UDPServer{
		server:     server,
		closedChan: make(chan struct{}),
	}
}

func (s *UDPServer) GetUdpConn() *net.UDPConn {
	return s.udpConn
}

// Start 启动UDP服务器
func (s *UDPServer) Start(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr(s.server.network, s.server.address)
	if err != nil {
		return fmt.Errorf("解析 UDP 地址失败: %w", err)
	}
	udp, err := net.ListenUDP(s.server.network, addr)
	if err != nil {
		return fmt.Errorf("UDP 监听失败: %w", err)
	}
	s.udpConn = udp
	if s.server.readDeadline > 0 {
		_ = s.udpConn.SetReadDeadline(time.Now().Add(s.server.readDeadline))
	}
	return nil
}

// Stop 停止UDP服务器
func (s *UDPServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		fmt.Println("UDPServer 已经关闭，无需重复 Stop")
		return nil
	}

	s.closed = true
	close(s.closedChan)

	if s.udpConn != nil {
		if err := s.udpConn.Close(); err != nil {
			return err
		}
		s.udpConn = nil
	}
	return nil
}

// Endpoint 返回服务器端点
func (s *UDPServer) Endpoint() (*url.URL, error) {
	if s.udpConn != nil {
		addr := s.udpConn.LocalAddr().String()
		return url.Parse("socket://" + addr)
	}
	return url.Parse("socket://" + s.server.address)
}

// SendTo 向指定目标发送UDP数据
func (s *UDPServer) SendTo(targetAddr string, data []byte) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, ErrServerClosed
	}
	return s.server.SendTo(targetAddr, data)
}

// Broadcast 向所有目标广播UDP数据
func (s *UDPServer) Broadcast(data []byte) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, ErrServerClosed
	}
	return s.server.Broadcast(data)
}
