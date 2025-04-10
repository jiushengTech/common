package socket

import (
	"context"
	"errors"
	klog "github.com/jiushengTech/common/log/klog/logger"
	"io"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)

	ErrServerStopped  = errors.New("socket: server has been stopped")
	ErrInvalidNetwork = errors.New("socket: unsupported network type")
	ErrEmptyAddress   = errors.New("socket: address cannot be empty")
)

// Server 是一个简单的 socket 服务器。
type Server struct {
	mu            sync.RWMutex
	conns         map[string]net.Conn // 存储连接
	listener      net.Listener        // TCP 监听器
	packetConn    net.PacketConn      // UDP 包连接
	network       string
	address       string
	targetAddrs   []string // 支持多个目标地址
	timeout       time.Duration
	readDeadline  time.Duration
	writeDeadline time.Duration
	connPool      sync.Pool     // 连接池以复用连接
	done          chan struct{} // 用于通知goroutine停止
	logger        log.Logger    // 日志记录器
	running       bool          // 服务器运行状态
}

// NewServer 使用提供的选项创建一个新的 Server。
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		timeout:       1 * time.Second,
		readDeadline:  30 * time.Second,
		writeDeadline: 30 * time.Second,
		conns:         make(map[string]net.Conn),
		done:          make(chan struct{}),
		logger:        log.DefaultLogger,
	}

	// 设置连接池
	srv.connPool.New = func() interface{} {
		return &connWrapper{
			srv: srv,
		}
	}

	srv.init(opts...)
	return srv
}

// connWrapper 是对连接的包装，用于连接池
type connWrapper struct {
	conn net.Conn
	srv  *Server
}

// init 应用选项到服务器。
func (s *Server) init(opts ...ServerOption) {
	for _, o := range opts {
		o(s)
	}
}

// listen 开始在指定的网络和地址上监听。
func (s *Server) listen() error {
	if s.address == "" {
		return ErrEmptyAddress
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return nil // 服务器已在运行
	}

	var err error
	switch s.network {
	case "tcp", "tcp4", "tcp6":
		err = s.listenTCP()
	case "udp", "udp4", "udp6":
		err = s.listenUDP()
	default:
		return ErrInvalidNetwork
	}

	if err == nil {
		s.running = true
	}
	return err
}

// listenTCP 在 TCP 网络上开始监听。
func (s *Server) listenTCP() error {
	addr, err := net.ResolveTCPAddr(s.network, s.address)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP(s.network, addr)
	if err != nil {
		return err
	}
	s.listener = listener

	// 接受连接并将其存储在 map 中
	go s.acceptLoop()

	return nil
}

// acceptLoop 接受新的TCP连接
func (s *Server) acceptLoop() {
	for {
		select {
		case <-s.done:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				var ne net.Error
				if errors.As(err, &ne) && ne.Temporary() {
					// 临时错误，稍后重试
					time.Sleep(100 * time.Millisecond)
					continue
				}
				klog.Log.Error("accept error:", err)
				return // 严重错误，退出循环
			}

			// 设置连接超时
			if s.readDeadline > 0 {
				conn.SetReadDeadline(time.Now().Add(s.readDeadline))
			}
			if s.writeDeadline > 0 {
				conn.SetWriteDeadline(time.Now().Add(s.writeDeadline))
			}

			remoteAddr := conn.RemoteAddr().String()
			s.mu.Lock()
			s.conns[remoteAddr] = conn
			s.mu.Unlock()
			klog.Log.Info("new connection addr", remoteAddr)
		}
	}
}

// listenUDP 在 UDP 网络上开始监听。
func (s *Server) listenUDP() error {
	udpAddr, err := net.ResolveUDPAddr(s.network, s.address)
	if err != nil {
		return err
	}

	packetConn, err := net.ListenUDP(s.network, udpAddr)
	if err != nil {
		return err
	}
	s.packetConn = packetConn

	// 设置UDP连接超时
	if s.readDeadline > 0 {
		packetConn.SetReadDeadline(time.Now().Add(s.readDeadline))
	}
	if s.writeDeadline > 0 {
		packetConn.SetWriteDeadline(time.Now().Add(s.writeDeadline))
	}

	// 对于 UDP，使用本地地址作为键
	localAddr := packetConn.LocalAddr().String()
	s.conns[localAddr] = packetConn

	// 启动UDP读取循环
	go s.readUDPLoop(packetConn)

	return nil
}

// readUDPLoop 持续从UDP连接读取数据
func (s *Server) readUDPLoop(conn *net.UDPConn) {
	buffer := make([]byte, 4096)
	for {
		select {
		case <-s.done:
			return
		default:
			n, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				var ne net.Error
				if errors.As(err, &ne) && ne.Temporary() {
					continue
				}
				klog.Log.Info("msg UDP read error", err)
				return
			}

			// 处理收到的UDP数据
			s.handleUDPPacket(buffer[:n], addr)
		}
	}
}

// handleUDPPacket 处理接收到的UDP数据包
func (s *Server) handleUDPPacket(data []byte, addr *net.UDPAddr) {
	// 这里可以实现UDP数据包处理逻辑
	klog.Log.Info("received UDP packet from", addr.String(), "size", len(data))
}

// Endpoint 返回服务器的URL端点。
func (s *Server) Endpoint() (*url.URL, error) {
	addr := "socket://" + s.address
	return url.Parse(addr)
}

// Start 启动服务器并开始监听传入连接。
func (s *Server) Start(ctx context.Context) error {
	return s.listen()
}

// Stop 通过关闭连接停止服务器。
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	// 通知所有goroutine停止
	close(s.done)
	s.running = false

	var multiErr error

	// 关闭TCP监听器
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			multiErr = errors.Join(multiErr, err)
		}
		s.listener = nil
	}

	// 关闭UDP连接
	if s.packetConn != nil {
		if err := s.packetConn.Close(); err != nil {
			multiErr = errors.Join(multiErr, err)
		}
		s.packetConn = nil
	}

	// 关闭所有连接
	for addr, conn := range s.conns {
		if err := conn.Close(); err != nil {
			multiErr = errors.Join(multiErr, err)
		}
		delete(s.conns, addr)
	}

	return multiErr
}

// Broadcast 向所有目标地址发送数据。
func (s *Server) Broadcast(data []byte) (int, error) {
	s.mu.RLock()
	if !s.running {
		s.mu.RUnlock()
		return 0, ErrServerStopped
	}
	targets := make([]string, len(s.targetAddrs))
	copy(targets, s.targetAddrs)
	s.mu.RUnlock()

	var (
		totalSent int
		multiErr  error
		wg        sync.WaitGroup
	)

	results := make(chan struct {
		n   int
		err error
	}, len(targets))

	// 并行向所有目标地址发送数据
	for _, target := range targets {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			n, err := s.sendToWithConnPool(target, data)
			results <- struct {
				n   int
				err error
			}{n, err}
		}(target)
	}

	// 等待所有发送完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	for result := range results {
		totalSent += result.n
		if result.err != nil {
			multiErr = errors.Join(multiErr, result.err)
		}
	}

	return totalSent, multiErr
}

// sendToWithConnPool 使用连接池发送数据
func (s *Server) sendToWithConnPool(targetAddr string, data []byte) (int, error) {
	// 从连接池获取一个连接包装器
	wrapper := s.connPool.Get().(*connWrapper)
	defer s.connPool.Put(wrapper)

	// 如果连接不存在或已关闭，创建新连接
	if wrapper.conn == nil {
		conn, err := net.DialTimeout(s.network, targetAddr, s.timeout)
		if err != nil {
			return 0, err
		}
		wrapper.conn = conn
	}

	// 设置写入超时
	if s.writeDeadline > 0 {
		wrapper.conn.SetWriteDeadline(time.Now().Add(s.writeDeadline))
	}

	// 发送数据
	n, err := wrapper.conn.Write(data)

	// 处理连接错误
	if err != nil {
		// 如果是连接相关错误，关闭连接并创建新连接重试
		if isConnectionError(err) {
			wrapper.conn.Close()
			conn, dialErr := net.DialTimeout(s.network, targetAddr, s.timeout)
			if dialErr != nil {
				return 0, dialErr
			}
			wrapper.conn = conn

			// 重试一次
			if s.writeDeadline > 0 {
				wrapper.conn.SetWriteDeadline(time.Now().Add(s.writeDeadline))
			}
			return wrapper.conn.Write(data)
		}
		return n, err
	}

	return n, nil
}

// SendTo 向特定目标地址发送数据。
func (s *Server) SendTo(targetAddr string, data []byte) (int, error) {
	s.mu.RLock()
	running := s.running
	s.mu.RUnlock()

	if !running {
		return 0, ErrServerStopped
	}

	return s.sendToWithConnPool(targetAddr, data)
}

// GetConn 获取与特定地址关联的连接
func (s *Server) GetConn(addr string) (net.Conn, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conn, exists := s.conns[addr]
	return conn, exists
}

// isConnectionError 检查错误是否是连接相关错误
func isConnectionError(err error) bool {
	if err == io.EOF {
		return true
	}
	if ne, ok := err.(net.Error); ok {
		return ne.Timeout() || !ne.Temporary()
	}
	return false
}
