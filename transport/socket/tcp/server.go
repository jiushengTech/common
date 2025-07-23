package tcp

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/jiushengTech/common/transport/socket"

	"github.com/go-kratos/kratos/v2/transport"
)

// Server TCP服务器实现
type Server struct {
	mu          sync.RWMutex
	config      *Config
	closed      bool
	closedChan  chan struct{}
	tcpListener *net.TCPListener
	isStarted   bool
}

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

// NewServer 创建TCP服务器
func NewServer(opts ...Option) *Server {
	config := &Config{
		Network:    "tcp",
		Timeout:    30 * time.Second,
		KeepAlive:  30 * time.Second,
		BufferSize: 4096,
		MaxRetries: 3,
	}

	for _, opt := range opts {
		opt(config)
	}

	return &Server{
		config:     config,
		closedChan: make(chan struct{}),
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
	addr, err := net.ResolveTCPAddr(s.config.Network, s.config.Address)
	if err != nil {
		return fmt.Errorf("解析 TCP 地址失败: %w", err)
	}
	tcp, err := net.ListenTCP(s.config.Network, addr)
	if err != nil {
		return fmt.Errorf("TCP 监听失败: %w", err)
	}

	s.tcpListener = tcp
	s.isStarted = true
	return nil
}

// Stop 停止TCP服务器
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		fmt.Println("TCP Server 已经关闭，无需重复 Stop")
		return nil
	}

	s.closed = true
	close(s.closedChan)

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
	return url.Parse("tcp://" + s.config.Address)
}

// SendTo 向指定目标发送TCP数据
func (s *Server) SendTo(targetAddr string, data []byte) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, socket.ErrServerClosed
	}

	conn, err := s.createTCPConnection(targetAddr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return s.writeWithRetry(conn, data)
}

// Broadcast 向所有目标广播TCP数据
func (s *Server) Broadcast(data []byte) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, socket.ErrServerClosed
	}

	if len(s.config.TargetAddrs) == 0 {
		return 0, fmt.Errorf("没有可广播的目标地址")
	}

	var totalBytes int
	errorMap := make(map[string]error)

	for _, target := range s.config.TargetAddrs {
		n, err := s.SendTo(target, data)
		totalBytes += n
		if err != nil {
			errorMap[target] = err
		}
	}

	if len(errorMap) > 0 {
		var errMsg string
		for target, err := range errorMap {
			errMsg += fmt.Sprintf("目标 [%s]: %v; ", target, err)
		}
		return totalBytes, fmt.Errorf("TCP广播部分失败: %s", errMsg)
	}

	return totalBytes, nil
}

// createTCPConnection 创建TCP连接
func (s *Server) createTCPConnection(targetAddr string) (*net.TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		return nil, fmt.Errorf("解析TCP地址失败: %w", err)
	}

	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, fmt.Errorf("TCP连接失败: %w", err)
	}

	// 设置TCP特有的选项
	if s.config.KeepAlive > 0 {
		err = conn.SetKeepAlive(true)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("设置KeepAlive失败: %w", err)
		}
		err = conn.SetKeepAlivePeriod(s.config.KeepAlive)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("设置KeepAlive周期失败: %w", err)
		}
	}

	// 设置读写超时
	if s.config.ReadTimeout > 0 {
		err = conn.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("设置读超时失败: %w", err)
		}
	}

	if s.config.WriteTimeout > 0 {
		err = conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("设置写超时失败: %w", err)
		}
	}

	return conn, nil
}

// writeWithRetry 带重试的写入
func (s *Server) writeWithRetry(conn *net.TCPConn, data []byte) (int, error) {
	var lastErr error

	for i := 0; i <= s.config.MaxRetries; i++ {
		n, err := conn.Write(data)
		if err == nil {
			return n, nil
		}
		lastErr = err

		if i < s.config.MaxRetries {
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond) // 指数退避
		}
	}

	return 0, fmt.Errorf("TCP写入重试%d次后失败: %w", s.config.MaxRetries, lastErr)
}
