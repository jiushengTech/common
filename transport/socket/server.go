package socket

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Server 是一个简单的 socket 服务器。
type Server struct {
	mu            sync.Mutex
	err           error
	network       string
	address       string
	targetAddr    []string // 支持多个目标地址
	timeout       time.Duration
	deadline      time.Duration
	readDeadline  time.Duration
	writeDeadline time.Duration
	//  设置缓冲区大小
	bufferSize int
	//  设置最大连接数
	maxConns int
	connPool *ConnectionPool
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
	if s.connPool == nil {
		s.connPool = NewConnectionPool(s.maxConns)
	}
}

// Broadcast 向所有目标地址广播数据
func (s *Server) Broadcast(data []byte) (int, error) {
	s.mu.Lock()
	targets := make([]string, len(s.targetAddr))
	copy(targets, s.targetAddr)
	s.mu.Unlock()

	if len(targets) == 0 {
		return 0, fmt.Errorf("没有可广播的目标地址")
	}

	var totalBytes int
	errorMap := make(map[string]error) // 使用map存储每个地址的错误

	for _, target := range targets {
		n, err := s.SendTo(target, data)
		totalBytes += n
		if err != nil {
			errorMap[target] = err
		}
	}

	// 如果有错误发生，将所有错误组合成一个
	if len(errorMap) > 0 {
		var errMsg string
		for target, err := range errorMap {
			errMsg += fmt.Sprintf("目标 [%s]: %v; ", target, err)
		}
		return totalBytes, fmt.Errorf("广播部分失败: %s", errMsg)
	}

	return totalBytes, nil
}

// SendTo 向指定目标地址发送数据
func (s *Server) SendTo(targetAddr string, data []byte) (int, error) {
	s.mu.Lock()
	s.mu.Unlock()

	if s.network == "udp" {
		// --- UDP 简洁发送，不入连接池 ---
		raddr, err := net.ResolveUDPAddr("udp", targetAddr)
		if err != nil {
			return 0, fmt.Errorf("解析目标地址失败: %w", err)
		}

		conn, err := net.DialUDP("udp", nil, raddr)
		if err != nil {
			return 0, fmt.Errorf("连接目标 %s 失败: %w", targetAddr, err)
		}
		defer conn.Close()

		if s.deadline > 0 {
			_ = conn.SetDeadline(time.Now().Add(s.deadline))
		}

		n, err := conn.Write(data)
		if err != nil {
			return 0, fmt.Errorf("写入目标 %s 失败: %w", targetAddr, err)
		}
		return n, nil
	}

	// --- TCP 情况，走连接池 ---
	conn, err := s.connPool.GetConn(s.network, targetAddr, s.timeout)
	if err != nil {
		return 0, err
	}
	if s.deadline > 0 {
		err = conn.SetDeadline(time.Now().Add(s.deadline))
		return 0, err
	}
	n, err := conn.Write(data)
	if err != nil {
		return 0, fmt.Errorf("写入目标 %s 失败: %w", targetAddr, err)
	}
	return n, nil
}
