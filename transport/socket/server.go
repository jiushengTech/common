package socket

import (
	"context"
	"errors"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/transport"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

// Server 是一个简单的 socket 服务器。
type Server struct {
	mu            sync.Mutex
	err           error
	tcpListener   *net.TCPListener
	udpConn       *net.UDPConn
	network       string
	address       string
	targetAddr    []string // 支持多个目标地址
	timeout       time.Duration
	deadline      time.Duration
	readDeadline  time.Duration
	writeDeadline time.Duration
}

// NewServer 使用提供的选项创建一个新的 Server。
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{}
	srv.init(opts...)
	return srv
}

// init applies the options to the Server.
func (s *Server) init(opts ...ServerOption) {
	for _, o := range opts {
		o(s)
	}
}

// listen starts listening on the specified network and address.
func (s *Server) listen() error {
	if s.address == "" {
		return errors.New("socket初始化失败, address为空")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.network {
	case "tcp":
		return s.listenTCP()
	case "udp":
		return s.listenUDP()
	default:
		return errors.New("unsupported network type")
	}
}

// listenTCP 在 TCP 网络上开始监听。
func (s *Server) listenTCP() error {
	addr, err := net.ResolveTCPAddr(s.network, s.address)
	if err != nil {
		return err
	}
	tcp, err := net.ListenTCP(s.network, addr)
	if err != nil {
		return err
	}
	s.tcpListener = tcp
	return nil
}

// listenUDP 在 UDP 网络上开始监听。
func (s *Server) listenUDP() error {
	udpAddr, err := net.ResolveUDPAddr(s.network, s.address)
	if err != nil {
		return err
	}
	s.udpConn, err = net.ListenUDP(s.network, udpAddr)
	if err != nil {
		return err
	}

	return nil
}

// Endpoint returns the URL endpoint for the server.
func (s *Server) Endpoint() (*url.URL, error) {
	addr := "socket://" + s.address
	return url.Parse(addr)
}

// Start starts the server and begins listening for incoming connections.
func (s *Server) Start(ctx context.Context) error {
	return s.listen()
}

// Stop stops the server by closing the connection.
// Stop stops the server by closing the connection and cleaning up resources.
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 关闭 TCP 监听器
	if s.tcpListener != nil {
		if err := s.tcpListener.Close(); err != nil {
			return err
		}
		s.tcpListener = nil
	}

	// 关闭 UDP 连接
	if s.udpConn != nil {
		if err := s.udpConn.Close(); err != nil {
			return err
		}
		s.udpConn = nil
	}

	return nil
}

// Broadcast sends data to all target addresses.
func (s *Server) Broadcast(data []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var totalSent int
	var lastErr error

	// 遍历所有目标地址，向每个地址发送数据
	for _, target := range s.targetAddr {
		conn, err := net.DialTimeout(s.network, target, s.timeout)
		if err != nil {
			lastErr = err
			continue
		}
		n, err := conn.Write(data)
		if err != nil {
			lastErr = err
		} else {
			totalSent += n
		}
	}

	if lastErr != nil {
		return totalSent, lastErr
	}
	return totalSent, nil
}

// SendTo sends data to a specific target address.
func (s *Server) SendTo(targetAddr string, data []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 创建到指定 target 的连接
	conn, err := net.DialTimeout(s.network, targetAddr, s.timeout)
	if err != nil {
		return 0, err
	}
	// 发送数据
	n, err := conn.Write(data)
	if err != nil {
		return 0, err
	}
	return n, nil
}
