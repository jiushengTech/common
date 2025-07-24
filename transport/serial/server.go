package serial

import (
	"context"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"net/url"
)

type Server struct {
	serial.OpenOptions
	Conn io.ReadWriteCloser
}

func NewServer(opts ...Option) *Server {
	srv := Server{
		OpenOptions: serial.OpenOptions{
			PortName:        "COM3",
			BaudRate:        9600,
			DataBits:        8,
			StopBits:        1,
			MinimumReadSize: 4,
		},
	}
	for _, o := range opts {
		o(&srv)
	}
	return &srv
}

func (s *Server) Endpoint() (*url.URL, error) {
	addr := s.PortName
	prefix := "serial://"
	addr = prefix + addr
	var endpoint *url.URL
	endpoint, err := url.Parse(addr)
	return endpoint, err
}

func (s *Server) Start(ctx context.Context) error {
	// 打开串口
	conn, err := serial.Open(s.OpenOptions)
	if err != nil {
		panic(fmt.Errorf("open serial fail: %w", err))
		return nil
	}
	s.Conn = conn
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	if s.Conn != nil {
		return s.Conn.Close()
	}
	return nil
}
