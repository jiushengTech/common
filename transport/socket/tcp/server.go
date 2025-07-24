package tcp

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"github.com/go-kratos/kratos/v2/transport"
)

// Server TCP服务器实现
type Server struct {
	tcpListener *net.TCPListener
	isStarted   bool
	addr        string
}

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

// NewServer 创建TCP服务器
func NewServer(addr string) *Server {
	return &Server{
		tcpListener: nil,
		isStarted:   false,
		addr:        addr,
	}
}

func (s *Server) GetTcpListener() *net.TCPListener {
	for {
		if s.isStarted {
			return s.tcpListener
		}
	}
}

// Start 启动TCP服务器
func (s *Server) Start(ctx context.Context) error {
	addr, err := net.ResolveTCPAddr("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("解析 TCP 地址失败: %w", err)
	}
	tcp, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return fmt.Errorf("TCP 监听失败: %w", err)
	}

	s.tcpListener = tcp
	s.isStarted = true
	return nil
}

// Stop 停止TCP服务器
func (s *Server) Stop(ctx context.Context) error {
	if s.tcpListener != nil {
		if err := s.tcpListener.Close(); err != nil {
			return err
		}
		s.tcpListener = nil
	}
	return nil
}

// Endpoint 返回服务器端点
func (s *Server) Endpoint() (*url.URL, error) {
	if s.tcpListener != nil {
		addr := s.tcpListener.Addr().String()
		return url.Parse("tcp://" + addr)
	}
	return url.Parse("tcp://" + s.addr)
}
