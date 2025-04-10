package socket

import (
	"context"
	"errors"
	"fmt"
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
	//  设置缓冲区大小
	BufferSize int32
	//  设置最大连接数
	MaxConns int32
}

// NewServer 使用提供的选项创建一个新的 Server。
func NewServer(opts ...Option) *Server {
	srv := &Server{}
	srv.init(opts...)
	return srv
}

// init 应用配置选项
func (s *Server) init(opts ...Option) {
	for _, o := range opts {
		o(s)
	}
	if s.network == "" {
		s.network = "tcp"
	}
}

func (s *Server) GetTcpListener() *net.TCPListener {
	return s.tcpListener
}

func (s *Server) GetUdpConn() *net.UDPConn {
	return s.udpConn
}

// listen 启动监听服务
func (s *Server) listen() error {
	if s.address == "" {
		return errors.New("socket 初始化失败: address 为空")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.network {
	case "tcp":
		return s.listenTCP()
	case "udp":
		return s.listenUDP()
	default:
		return fmt.Errorf("不支持的网络类型: %s", s.network)
	}
}

// listenTCP 在 TCP 网络上开始监听
func (s *Server) listenTCP() error {
	addr, err := net.ResolveTCPAddr(s.network, s.address)
	if err != nil {
		return fmt.Errorf("解析 TCP 地址失败: %w", err)
	}
	tcp, err := net.ListenTCP(s.network, addr)
	if err != nil {
		return fmt.Errorf("TCP 监听失败: %w", err)
	}
	s.tcpListener = tcp
	return nil
}

// listenUDP 在 UDP 网络上开始监听
func (s *Server) listenUDP() error {
	addr, err := net.ResolveUDPAddr(s.network, s.address)
	if err != nil {
		return fmt.Errorf("解析 UDP 地址失败: %w", err)
	}
	udp, err := net.ListenUDP(s.network, addr)
	if err != nil {
		return fmt.Errorf("UDP 监听失败: %w", err)
	}
	s.udpConn = udp
	if s.readDeadline > 0 {
		_ = s.udpConn.SetReadDeadline(time.Now().Add(s.readDeadline))
	}
	return nil
}

// Endpoint 返回 Server 的 URL 形式
func (s *Server) Endpoint() (*url.URL, error) {
	addr := "socket://" + s.address
	return url.Parse(addr)
}

// Start 启动服务
func (s *Server) Start(ctx context.Context) error {
	return s.listen()
}

// Stop 停止服务并释放资源
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tcpListener != nil {
		if err := s.tcpListener.Close(); err != nil {
			return err
		}
		s.tcpListener = nil
	}
	if s.udpConn != nil {
		if err := s.udpConn.Close(); err != nil {
			return err
		}
		s.udpConn = nil
	}
	return nil
}

// Broadcast 向所有目标地址广播数据
func (s *Server) Broadcast(data []byte) (int, error) {
	// 先获取所需的配置信息，避免长时间持有锁
	s.mu.Lock()
	targets := make([]string, len(s.targetAddr))
	copy(targets, s.targetAddr)
	network := s.network
	timeout := s.timeout
	deadline := s.deadline
	s.mu.Unlock()

	if len(targets) == 0 {
		return 0, fmt.Errorf("没有设置目标地址")
	}

	var totalSent int
	var errors []error
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 并发发送数据
	for _, target := range targets {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()

			conn, err := net.DialTimeout(network, addr, timeout)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("连接目标 %s 失败: %w", addr, err))
				mu.Unlock()
				return
			}

			defer conn.Close()
			if deadline > 0 {
				_ = conn.SetDeadline(time.Now().Add(deadline))
			}

			n, err := conn.Write(data)

			mu.Lock()
			if err != nil {
				errors = append(errors, fmt.Errorf("写入目标 %s 失败: %w", addr, err))
			} else {
				totalSent += n
			}
			mu.Unlock()
		}(target)
	}

	wg.Wait()

	// 处理错误
	if len(errors) > 0 {
		var errMsg string
		if len(errors) == 1 {
			errMsg = errors[0].Error()
		} else {
			errMsg = fmt.Sprintf("广播时发生了 %d 个错误:", len(errors))
			for i, err := range errors {
				if i < 3 || i == len(errors)-1 { // 只显示前3个和最后一个错误
					errMsg += "\n- " + err.Error()
				} else if i == 3 {
					errMsg += fmt.Sprintf("\n- 还有 %d 个错误...", len(errors)-4)
				}
			}
		}
		return totalSent, fmt.Errorf(errMsg)
	}

	return totalSent, nil
}

// SendTo 向指定目标地址发送数据
func (s *Server) SendTo(targetAddr string, data []byte) (int, error) {
	// 先获取所需的配置信息，避免长时间持有锁
	s.mu.Lock()
	network := s.network
	timeout := s.timeout
	deadline := s.deadline
	s.mu.Unlock()

	conn, err := net.DialTimeout(network, targetAddr, timeout)
	if err != nil {
		return 0, fmt.Errorf("连接目标 %s 失败: %w", targetAddr, err)
	}
	defer conn.Close()

	if deadline > 0 {
		_ = conn.SetDeadline(time.Now().Add(deadline))
	}
	n, err := conn.Write(data)
	if err != nil {
		return 0, fmt.Errorf("写入目标 %s 失败: %w", targetAddr, err)
	}
	return n, nil
}
