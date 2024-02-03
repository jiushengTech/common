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
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

type Server struct {
	UdpConn    *net.UDPConn
	err        error
	network    string
	address    string
	targetAddr string
	timeout    time.Duration
	readBuffer int
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "udp",
		address: "0.0.0.0:30003",
		timeout: 1 * time.Second,
	}
	srv.init(opts...)
	if srv.readBuffer != 0 {
		err := srv.UdpConn.SetReadBuffer(srv.readBuffer)
		if err != nil {
			panic(err)
		}
	}
	return srv
}

func (s *Server) init(opts ...ServerOption) {
	for _, o := range opts {
		o(s)
	}
}

func (s *Server) listen() error {
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

func (s *Server) Endpoint() (*url.URL, error) {
	addr := s.address
	prefix := "socket://"
	addr = prefix + addr
	var endpoint *url.URL
	endpoint, s.err = url.Parse(addr)
	return endpoint, nil
}

func (s *Server) Start(ctx context.Context) error {
	err := s.listen()
	if err != nil {
		return err
	}
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	if s.UdpConn != nil {
		return s.UdpConn.Close()
	}
	s.UdpConn = nil
	return nil
}

func (s *Server) Send(data []byte) (int, error) {
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
