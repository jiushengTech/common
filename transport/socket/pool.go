package socket

import (
	"net"
	"sync"
	"time"
)

// 连接键
type connKey struct {
	network string
	address string
}

// ConnectionPool 连接池
type ConnectionPool struct {
	mu       sync.Mutex
	conns    map[connKey][]net.Conn
	maxConns int
}

// NewConnectionPool 创建连接池
func NewConnectionPool(maxConns int) *ConnectionPool {
	if maxConns <= 0 {
		maxConns = 100
	}

	return &ConnectionPool{
		conns:    make(map[connKey][]net.Conn),
		maxConns: maxConns,
	}
}

// Get 获取连接
func (p *ConnectionPool) Get(network, address string, timeout time.Duration) (net.Conn, error) {
	key := connKey{network: network, address: address}

	p.mu.Lock()
	if conns, ok := p.conns[key]; ok && len(conns) > 0 {
		conn := conns[len(conns)-1]
		p.conns[key] = conns[:len(conns)-1]
		p.mu.Unlock()
		return conn, nil
	}
	p.mu.Unlock()

	// 如果没有可用连接，创建新连接
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, &NetworkError{Op: "创建连接", Err: err}
	}

	return conn, nil
}

// Put 归还连接
func (p *ConnectionPool) Put(conn net.Conn) {
	if conn == nil {
		return
	}

	key := connKey{
		network: conn.RemoteAddr().Network(),
		address: conn.RemoteAddr().String(),
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// 检查连接池是否已满
	if conns, ok := p.conns[key]; ok {
		if len(conns) >= p.maxConns {
			conn.Close()
			return
		}
		p.conns[key] = append(conns, conn)
	} else {
		p.conns[key] = []net.Conn{conn}
	}
}

// Close 关闭所有连接
func (p *ConnectionPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, conns := range p.conns {
		for _, conn := range conns {
			conn.Close()
		}
	}

	p.conns = make(map[connKey][]net.Conn)
}
