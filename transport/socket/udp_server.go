package socket

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/transport"
)

// UDPServerConfig UDP服务器配置
type UDPServerConfig struct {
	Network         string
	Address         string
	TargetAddrs     []string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	BufferSize      int
	MaxPacketSize   int
	EnableBroadcast bool
}

// UDPServer UDP服务器实现
type UDPServer struct {
	mu         sync.RWMutex
	config     *UDPServerConfig
	closed     bool
	closedChan chan struct{}
	udpConn    *net.UDPConn
	isStarted  bool
}

var (
	_ transport.Server     = (*UDPServer)(nil)
	_ transport.Endpointer = (*UDPServer)(nil)
)

// NewUDPServer 创建UDP服务器
func NewUDPServer(opts ...UDPOption) *UDPServer {
	config := &UDPServerConfig{
		Network:         "udp",
		BufferSize:      4096,
		MaxPacketSize:   65507, // UDP最大数据包大小
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    10 * time.Second,
		EnableBroadcast: false,
	}

	for _, opt := range opts {
		opt(config)
	}

	return &UDPServer{
		config:     config,
		closedChan: make(chan struct{}),
	}
}

func (s *UDPServer) GetUdpConn() *net.UDPConn {
	for {
		if s.isStarted {
			return s.udpConn
		}
	}
}

// Start 启动UDP服务器
func (s *UDPServer) Start(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr(s.config.Network, s.config.Address)
	if err != nil {
		return fmt.Errorf("解析 UDP 地址失败: %w", err)
	}

	udp, err := net.ListenUDP(s.config.Network, addr)
	if err != nil {
		return fmt.Errorf("UDP 监听失败: %w", err)
	}

	s.udpConn = udp

	// 设置读超时
	if s.config.ReadTimeout > 0 {
		_ = s.udpConn.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))
	}

	s.isStarted = true
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
		return url.Parse("udp://" + addr)
	}
	return url.Parse("udp://" + s.config.Address)
}

// SendTo 向指定目标发送UDP数据
func (s *UDPServer) SendTo(targetAddr string, data []byte) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, ErrServerClosed
	}

	// 检查数据包大小
	if len(data) > s.config.MaxPacketSize {
		return 0, fmt.Errorf("数据包大小 %d 超过最大限制 %d", len(data), s.config.MaxPacketSize)
	}

	return s.sendUDPPacket(targetAddr, data)
}

// Broadcast 向所有目标广播UDP数据
func (s *UDPServer) Broadcast(data []byte) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, ErrServerClosed
	}

	if len(s.config.TargetAddrs) == 0 {
		return 0, fmt.Errorf("没有可广播的目标地址")
	}

	// 检查数据包大小
	if len(data) > s.config.MaxPacketSize {
		return 0, fmt.Errorf("广播数据包大小 %d 超过最大限制 %d", len(data), s.config.MaxPacketSize)
	}

	var totalBytes int
	errorMap := make(map[string]error)

	// UDP广播可以并发发送，提高性能
	type result struct {
		target string
		bytes  int
		err    error
	}

	resultChan := make(chan result, len(s.config.TargetAddrs))

	for _, target := range s.config.TargetAddrs {
		go func(addr string) {
			n, err := s.sendUDPPacket(addr, data)
			resultChan <- result{target: addr, bytes: n, err: err}
		}(target)
	}

	// 收集结果
	for i := 0; i < len(s.config.TargetAddrs); i++ {
		res := <-resultChan
		totalBytes += res.bytes
		if res.err != nil {
			errorMap[res.target] = res.err
		}
	}

	if len(errorMap) > 0 {
		var errMsg string
		for target, err := range errorMap {
			errMsg += fmt.Sprintf("目标 [%s]: %v; ", target, err)
		}
		return totalBytes, fmt.Errorf("UDP广播部分失败: %s", errMsg)
	}

	return totalBytes, nil
}

// MulticastTo 向组播地址发送数据
func (s *UDPServer) MulticastTo(multicastAddr string, data []byte) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, ErrServerClosed
	}

	if len(data) > s.config.MaxPacketSize {
		return 0, fmt.Errorf("组播数据包大小 %d 超过最大限制 %d", len(data), s.config.MaxPacketSize)
	}

	return s.sendUDPPacket(multicastAddr, data)
}

// sendUDPPacket 发送UDP数据包
func (s *UDPServer) sendUDPPacket(targetAddr string, data []byte) (int, error) {
	raddr, err := net.ResolveUDPAddr("udp", targetAddr)
	if err != nil {
		return 0, fmt.Errorf("解析UDP地址失败: %w", err)
	}

	// UDP是无连接的，每次都创建新的连接
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return 0, fmt.Errorf("UDP连接失败: %w", err)
	}
	defer conn.Close()

	// 设置写超时
	if s.config.WriteTimeout > 0 {
		err = conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
		if err != nil {
			return 0, fmt.Errorf("设置UDP写超时失败: %w", err)
		}
	}

	n, err := conn.Write(data)
	if err != nil {
		return 0, fmt.Errorf("UDP写入目标 %s 失败: %w", targetAddr, err)
	}

	return n, nil
}

// BatchSendTo 批量发送UDP数据到多个目标（高性能版本）
func (s *UDPServer) BatchSendTo(targets []string, data []byte) (map[string]int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, ErrServerClosed
	}

	if len(data) > s.config.MaxPacketSize {
		return nil, fmt.Errorf("批量发送数据包大小 %d 超过最大限制 %d", len(data), s.config.MaxPacketSize)
	}

	type result struct {
		target string
		bytes  int
		err    error
	}

	resultChan := make(chan result, len(targets))

	// 并发发送到所有目标
	for _, target := range targets {
		go func(addr string) {
			n, err := s.sendUDPPacket(addr, data)
			resultChan <- result{target: addr, bytes: n, err: err}
		}(target)
	}

	// 收集结果
	results := make(map[string]int)
	var errors []string

	for i := 0; i < len(targets); i++ {
		res := <-resultChan
		if res.err == nil {
			results[res.target] = res.bytes
		} else {
			errors = append(errors, fmt.Sprintf("%s: %v", res.target, res.err))
		}
	}

	if len(errors) > 0 {
		return results, fmt.Errorf("批量发送部分失败: %v", errors)
	}

	return results, nil
}

// UDPOption UDP配置选项
type UDPOption func(*UDPServerConfig)

func WithUDPNetwork(network string) UDPOption {
	return func(c *UDPServerConfig) {
		c.Network = network
	}
}

func WithUDPAddress(addr string) UDPOption {
	return func(c *UDPServerConfig) {
		c.Address = addr
	}
}

func WithUDPTargetAddrs(addrs []string) UDPOption {
	return func(c *UDPServerConfig) {
		c.TargetAddrs = addrs
	}
}

func WithUDPReadTimeout(timeout time.Duration) UDPOption {
	return func(c *UDPServerConfig) {
		c.ReadTimeout = timeout
	}
}

func WithUDPWriteTimeout(timeout time.Duration) UDPOption {
	return func(c *UDPServerConfig) {
		c.WriteTimeout = timeout
	}
}

func WithUDPBufferSize(size int) UDPOption {
	return func(c *UDPServerConfig) {
		c.BufferSize = size
	}
}

func WithUDPMaxPacketSize(size int) UDPOption {
	return func(c *UDPServerConfig) {
		c.MaxPacketSize = size
	}
}

func WithUDPBroadcast(enable bool) UDPOption {
	return func(c *UDPServerConfig) {
		c.EnableBroadcast = enable
	}
}
