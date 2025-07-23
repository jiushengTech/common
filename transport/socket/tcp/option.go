package tcp

import "time"

// Config TCP服务器配置
type Config struct {
	Network      string
	Address      string
	TargetAddrs  []string
	Timeout      time.Duration
	KeepAlive    time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	BufferSize   int
	MaxRetries   int
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

func WithTargetAddrs(addrs []string) Option {
	return func(c *Config) {
		c.TargetAddrs = addrs
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithKeepAlive(keepAlive time.Duration) Option {
	return func(c *Config) {
		c.KeepAlive = keepAlive
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

func WithBufferSize(size int) Option {
	return func(c *Config) {
		c.BufferSize = size
	}
}

func WithMaxRetries(retries int) Option {
	return func(c *Config) {
		c.MaxRetries = retries
	}
}
