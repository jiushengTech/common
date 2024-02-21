package socket

import (
	"time"
)

type ServerOption func(o *Server)

//func WithReadBuffer(readBuffer int) ServerOption {
//	return func(s *Server) {
//		s.readBuffer = readBuffer
//	}
//}

func WithNetwork(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithTargetAddr(targetAddr string) ServerOption {
	return func(s *Server) {
		s.targetAddr = targetAddr
	}
}

func WithDeadline(Deadline time.Duration) ServerOption {
	return func(s *Server) {
		s.readDeadline = Deadline
	}
}

func WithReadDeadline(readDeadline time.Duration) ServerOption {
	return func(s *Server) {
		s.readDeadline = readDeadline
	}
}

func WithWriteDeadline(writeDeadline time.Duration) ServerOption {
	return func(s *Server) {
		s.writeDeadline = writeDeadline
	}
}
