package tcp

import (
	"net"
	"sync"
	"time"
)

// ClientConn 客户端连接封装
type ClientConn struct {
	*net.TCPConn
	ID       string
	RemoteIP string
	lastSeen time.Time
	mu       sync.Mutex
}

// Send 向客户端发送数据
func (c *ClientConn) Send(data []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastSeen = time.Now()
	return c.TCPConn.Write(data)
}

// Close 关闭客户端连接
func (c *ClientConn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.TCPConn.Close()
}

// IsAlive 检查连接是否活跃
func (c *ClientConn) IsAlive() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 尝试写入0字节数据来检查连接状态
	_, err := c.TCPConn.Write([]byte{})
	return err == nil
}

// LastSeen 获取最后活跃时间
func (c *ClientConn) LastSeen() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lastSeen
}

// newClientConn 创建新的客户端连接
func newClientConn(conn *net.TCPConn) *ClientConn {
	return &ClientConn{
		TCPConn:  conn,
		ID:       conn.RemoteAddr().String(),
		RemoteIP: conn.RemoteAddr().(*net.TCPAddr).IP.String(),
		lastSeen: time.Now(),
	}
}
