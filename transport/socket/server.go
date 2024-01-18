package socket

import (
	"context"
	"errors"
	"net"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/transport"
)

var (
	_ transport.Server     = (*RadarServer)(nil)
	_ transport.Endpointer = (*RadarServer)(nil)
)

type RadarServer struct {
	UdpConn    *net.UDPConn
	err        error
	network    string
	address    string
	targetAddr string
	timeout    time.Duration
}

func NewServer(opts ...ServerOption) *RadarServer {
	srv := &RadarServer{
		network: "udp",
		address: "0.0.0.0:30003",
		timeout: 1 * time.Second,
	}
	srv.init(opts...)
	return srv
}

func (s *RadarServer) init(opts ...ServerOption) {
	for _, o := range opts {
		o(s)
	}
}

func (s *RadarServer) listen() error {
	if s.address == "" {
		return errors.New("socket初始化失败, address为空")
	}
	udpAddr, err := net.ResolveUDPAddr(s.network, s.address)
	if errors.Is(err, net.UnknownNetworkError(s.address)) {
		return err
	}
	udpConn, err := net.ListenUDP(s.network, udpAddr)
	if err != nil {
		return err
	}
	s.UdpConn = udpConn
	return err
}

func (s *RadarServer) Endpoint() (*url.URL, error) {
	addr := s.address
	prefix := "socket://"
	addr = prefix + addr
	var endpoint *url.URL
	endpoint, s.err = url.Parse(addr)
	return endpoint, nil
}

func (s *RadarServer) Start(ctx context.Context) error {
	err := s.listen()
	if err != nil {
		return err
	}
	return err
}

func (s *RadarServer) Stop(ctx context.Context) error {
	return s.UdpConn.Close()
}

func (s *RadarServer) Send(data []byte) (int, error) {
	targetAddr, err := net.ResolveUDPAddr(s.network, s.targetAddr)
	if err != nil {
		return 0, err
	}
	i, err := s.UdpConn.WriteToUDP(data, targetAddr)
	if err != nil {
		return 0, err
	}
	return i, nil
}
