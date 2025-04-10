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
	Conns         map[string]net.Conn // 存储连接
	err           error
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
	srv := &Server{
		timeout: 1 * time.Second,
		Conns:   make(map[string]net.Conn), // 初始化 map
	}
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

	// 接受连接并将其存储在 map 中
	go func() {
		for {
			conn, err := tcp.AcceptTCP()
			if err != nil {
				// 记录错误，但不影响继续接受其他连接
				s.err = err
				continue
			}
			remoteAddr := conn.RemoteAddr().String()
			s.mu.Lock()
			s.Conns[remoteAddr] = conn
			s.mu.Unlock()
		}
	}()

	return nil
}

// listenUDP 在 UDP 网络上开始监听。
func (s *Server) listenUDP() error {
	udpAddr, err := net.ResolveUDPAddr(s.network, s.address)
	if err != nil {
		return err
	}
	udpConn, err := net.ListenUDP(s.network, udpAddr)
	if err != nil {
		return err
	}

	// 对于 UDP，RemoteAddr 可能为 nil，因此需要特殊处理
	// 这里使用本地地址作为键
	go func() {
		localAddr := udpConn.LocalAddr().String()
		s.mu.Lock()
		s.Conns[localAddr] = udpConn
		s.mu.Unlock()
	}()

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
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var lastErr error

	// 关闭所有连接
	for addr, conn := range s.Conns {
		if err := conn.Close(); err != nil {
			lastErr = err
		}
		delete(s.Conns, addr)
	}

	return lastErr
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
		conn.Close() // 发送完毕后关闭连接
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
	defer conn.Close()

	// 发送数据
	n, err := conn.Write(data)
	if err != nil {
		return 0, err
	}
	return n, nil
}
