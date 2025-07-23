package udp

import "time"

// Config UDP服务器配置
type Config struct {
	Network         string
	Address         string
	TargetAddrs     []string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	BufferSize      int
	MaxPacketSize   int
	EnableBroadcast bool
}

// Option UDP配置选项
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

func WithMaxPacketSize(size int) Option {
	return func(c *Config) {
		c.MaxPacketSize = size
	}
}

func WithBroadcast(enable bool) Option {
	return func(c *Config) {
		c.EnableBroadcast = enable
	}
}
