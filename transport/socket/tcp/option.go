package tcp

import "time"

// Config TCP服务器配置
type Config struct {
	Network           string
	Address           string
	KeepAlive         time.Duration
	KeepAliveInterval time.Duration
	ConnectionTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	MaxConnections    int
	CheckInterval     time.Duration
	EnableTCPNoDelay  bool
	SendBufferSize    int
	ReceiveBufferSize int
	EnableKeepalive   bool
	DataChannelSize   int // 数据通道缓冲区大小
	ReadBufferSize    int // 单次读取缓冲区大小
}

// Option TCP配置选项
type Option func(*Config)

func WithNetwork(network string) Option {
	return func(c *Config) {
		c.Network = network
	}
}

func WithAddress(addr string) Option {
	return func(c *Config) {
		c.Address = addr
	}
}

func WithKeepAlive(keepAlive time.Duration) Option {
	return func(c *Config) {
		c.KeepAlive = keepAlive
	}
}

func WithKeepAliveInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.KeepAliveInterval = interval
	}
}

func WithConnectionTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.ConnectionTimeout = timeout
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.WriteTimeout = timeout
	}
}

func WithMaxConnections(max int) Option {
	return func(c *Config) {
		c.MaxConnections = max
	}
}

func WithCheckInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.CheckInterval = interval
	}
}

func WithTCPNoDelay(enable bool) Option {
	return func(c *Config) {
		c.EnableTCPNoDelay = enable
	}
}

func WithSendBufferSize(size int) Option {
	return func(c *Config) {
		c.SendBufferSize = size
	}
}

func WithReceiveBufferSize(size int) Option {
	return func(c *Config) {
		c.ReceiveBufferSize = size
	}
}

func WithKeepalive(enable bool) Option {
	return func(c *Config) {
		c.EnableKeepalive = enable
	}
}

func WithDataChannelSize(size int) Option {
	return func(c *Config) {
		c.DataChannelSize = size
	}
}

func WithReadBufferSize(size int) Option {
	return func(c *Config) {
		c.ReadBufferSize = size
	}
}
