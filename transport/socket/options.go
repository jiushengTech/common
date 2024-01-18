package socket

import (
	"time"
)

type ServerOption func(o *RadarServer)

func WithNetwork(network string) ServerOption {
	return func(s *RadarServer) {
		s.network = network
	}
}

func WithAddress(addr string) ServerOption {
	return func(s *RadarServer) {
		s.address = addr
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *RadarServer) {
		s.timeout = timeout
	}
}

func WithTargetAddr(targetAddr string) ServerOption {
	return func(s *RadarServer) {
		s.targetAddr = targetAddr
	}
}
