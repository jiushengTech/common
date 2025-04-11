package socket

import "time"

type Option func(o *Server)

func WithNetwork(network string) Option {
	return func(s *Server) {
		s.network = network
	}
}

func WithAddress(addr string) Option {
	return func(s *Server) {
		s.address = addr
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithTargetAddr(targetAddr []string) Option {
	return func(s *Server) {
		s.targetAddr = targetAddr
	}
}

func WithDeadline(deadline time.Duration) Option {
	return func(s *Server) {
		s.deadline = deadline
	}
}

func WithReadDeadline(readDeadline time.Duration) Option {
	return func(s *Server) {
		s.readDeadline = readDeadline
	}
}

func WithWriteDeadline(writeDeadline time.Duration) Option {
	return func(s *Server) {
		s.writeDeadline = writeDeadline
	}
}

// WithBufferSize 设置缓冲区大小
func WithBufferSize(size int) Option {
	return func(s *Server) {
		s.bufferSize = size
	}
}

// WithMaxConns 设置最大连接数
func WithMaxConns(maxConns int) Option {
	return func(s *Server) {
		s.maxConns = maxConns
	}
}
