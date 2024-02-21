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

// Server is a simple socket server.
type Server struct {
	Conn       net.Conn
	client     net.Conn
	err        error
	network    string
	address    string
	targetAddr string
	timeout    time.Duration
}

// NewServer creates a new Server with the provided options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		timeout: 1 * time.Second,
	}
	srv.init(opts...)
	return srv
}

// init applies the options to the Server.
func (s *Server) init(opts ...ServerOption) {
	for _, o := range opts {
		o(s)
	}
}

// listen starts listening on the specified network and address.
func (s *Server) listen() error {
	if s.address == "" {
		return errors.New("socket初始化失败, address为空")
	}
	switch s.network {
	case "tcp":
		return s.listenTCP()
	case "udp":
		return s.listenUDP()
	default:
		return errors.New("unsupported network type")
	}
}

// listenTCP starts listening on a TCP network.
func (s *Server) listenTCP() error {
	addr, err := net.ResolveTCPAddr(s.network, s.address)
	if err != nil {
		return err
	}
	tcp, err := net.ListenTCP(s.network, addr)
	if err != nil {
		return err
	}
	conn, err := tcp.AcceptTCP()
	if err != nil {
		return err
	}
	s.Conn = conn
	return nil
}

// listenUDP starts listening on a UDP network.
func (s *Server) listenUDP() error {
	udpAddr, err := net.ResolveUDPAddr(s.network, s.address)
	if err != nil {
		return err
	}
	udpConn, err := net.ListenUDP(s.network, udpAddr)
	if err != nil {
		return err
	}
	s.Conn = udpConn
	return nil
}

// Endpoint returns the URL endpoint for the server.
func (s *Server) Endpoint() (*url.URL, error) {
	addr := "socket://" + s.address
	return url.Parse(addr)
}

// Start starts the server and begins listening for incoming connections.
func (s *Server) Start(ctx context.Context) error {
	return s.listen()
}

// Stop stops the server by closing the connection.
func (s *Server) Stop(ctx context.Context) error {
	if s.Conn != nil {
		return s.Conn.Close()
	}
	return nil
}

// Send sends data to the target address.
func (s *Server) Send(data []byte) (int, error) {
	var err error
	switch s.network {
	case "tcp", "udp":
		s.client, err = net.DialTimeout(s.network, s.targetAddr, s.timeout)
	default:
		return 0, errors.New("unsupported network type")
	}
	if err != nil {
		return 0, err
	}
	defer func(client net.Conn) {
		err := client.Close()
		if err != nil {
			return
		}
	}(s.client)

	i, err := s.client.Write(data)
	if err != nil {
		return 0, err
	}
	return i, nil
}
