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
