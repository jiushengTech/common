package tcp

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/transport"
)

// Server TCP服务器实现
type Server struct {
	mu          sync.RWMutex
	config      *Config
	tcpListener *net.TCPListener
	clients     sync.Map // map[string]*ClientConn
	running     bool
	handler     *EventHandler
}

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

// NewServer 创建TCP服务器
func NewServer(opts ...Option) *Server {
	config := &Config{
		Network:           "tcp",
		KeepAlive:         30 * time.Second,
		KeepAliveInterval: 15 * time.Second,
		ConnectionTimeout: 30 * time.Second,
		ReadTimeout:       0, // 0表示不设置超时
		WriteTimeout:      30 * time.Second,
		MaxConnections:    1000,
		CheckInterval:     5 * time.Second,
		EnableTCPNoDelay:  true,
		SendBufferSize:    0, // 0表示使用系统默认
		ReceiveBufferSize: 0, // 0表示使用系统默认
		EnableKeepalive:   true,
		DataChannelSize:   100,  // 数据通道缓冲区大小
		ReadBufferSize:    4096, // 单次读取缓冲区大小
	}

	for _, opt := range opts {
		opt(config)
	}

	return &Server{
		config:  config,
		handler: &EventHandler{}, // 默认空处理器
	}
}

// Start 启动TCP服务器并开始接受连接
func (s *Server) Start(ctx context.Context) error {
	addr, err := net.ResolveTCPAddr(s.config.Network, s.config.Address)
	if err != nil {
		return fmt.Errorf("解析 TCP 地址失败: %w", err)
	}

	listener, err := net.ListenTCP(s.config.Network, addr)
	if err != nil {
		return fmt.Errorf("TCP 监听失败: %w", err)
	}

	s.mu.Lock()
	s.tcpListener = listener
	s.running = true
	s.mu.Unlock()

	// 启动接受连接的协程
	go s.acceptConnections(ctx)

	return nil
}

// Stop 停止TCP服务器
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.running = false

	// 关闭监听器
	if s.tcpListener != nil {
		if err := s.tcpListener.Close(); err != nil {
			return err
		}
		s.tcpListener = nil
	}

	// 关闭所有客户端连接
	s.clients.Range(func(key, value interface{}) bool {
		if client, ok := value.(*ClientConn); ok {
			client.Close()
		}
		s.clients.Delete(key)
		return true
	})

	return nil
}

// Endpoint 返回服务器端点
func (s *Server) Endpoint() (*url.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.tcpListener != nil {
		addr := s.tcpListener.Addr().String()
		return url.Parse("tcp://" + addr)
	}
	return url.Parse("tcp://" + s.config.Address)
}

// GetClients 获取所有连接的客户端
func (s *Server) GetClients() map[string]*ClientConn {
	result := make(map[string]*ClientConn)
	s.clients.Range(func(key, value interface{}) bool {
		if clientID, ok := key.(string); ok {
			if client, ok := value.(*ClientConn); ok {
				result[clientID] = client
			}
		}
		return true
	})
	return result
}

// GetClient 根据ID获取客户端连接
func (s *Server) GetClient(clientID string) (*ClientConn, bool) {
	if value, ok := s.clients.Load(clientID); ok {
		if client, ok := value.(*ClientConn); ok {
			return client, true
		}
	}
	return nil, false
}

// GetClientCount 获取连接的客户端数量
func (s *Server) GetClientCount() int {
	count := 0
	s.clients.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// SendToClient 向指定客户端发送数据
func (s *Server) SendToClient(clientID string, data []byte) (int, error) {
	client, exists := s.GetClient(clientID)
	if !exists {
		return 0, fmt.Errorf("客户端 %s 不存在", clientID)
	}

	return client.Send(data)
}

// Broadcast 向所有客户端广播数据
func (s *Server) Broadcast(data []byte) (int, error) {
	var totalBytes int
	var errorList []string

	// 使用 sync.Map.Range 直接遍历，避免创建临时 map
	s.clients.Range(func(key, value interface{}) bool {
		if clientID, ok := key.(string); ok {
			if client, ok := value.(*ClientConn); ok {
				n, err := client.Send(data)
				totalBytes += n
				if err != nil {
					errorList = append(errorList, fmt.Sprintf("客户端 [%s]: %v", clientID, err))
				}
			}
		}
		return true // 继续遍历
	})

	if totalBytes == 0 {
		return 0, fmt.Errorf("没有连接的客户端")
	}

	if len(errorList) > 0 {
		return totalBytes, fmt.Errorf("广播部分失败: %s", joinErrors(errorList))
	}

	return totalBytes, nil
}

// joinErrors 连接错误字符串
func joinErrors(errors []string) string {
	if len(errors) == 0 {
		return ""
	}
	result := errors[0]
	for i := 1; i < len(errors); i++ {
		result += "; " + errors[i]
	}
	return result
}

// CloseClient 关闭指定客户端连接
func (s *Server) CloseClient(clientID string) error {
	if value, ok := s.clients.LoadAndDelete(clientID); ok {
		if client, ok := value.(*ClientConn); ok {
			return client.Close()
		}
	}
	return fmt.Errorf("客户端 %s 不存在", clientID)
}
