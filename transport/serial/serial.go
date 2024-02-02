package serial

import (
	"context"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"net/url"
)

type Server struct {
	serial.OpenOptions
	port io.ReadWriteCloser
	//PortName                string
	//BaudRate                uint
	//DataBits                uint
	//StopBits                uint
	//ParityMode              uint
	//RTSCTSFlowControl       bool
	//InterCharacterTimeout   uint
	//MinimumReadSize         uint
	//Rs485Enable             bool
	//Rs485RtsHighDuringSend  bool
	//Rs485RtsHighAfterSend   bool
	//Rs485RxDuringTx         bool
	//Rs485DelayRtsBeforeSend uint
	//Rs485DelayRtsAfterSend  uint
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
	port, err := serial.Open(s.OpenOptions)
	if err != nil {
		return err
	}
	s.port = port
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	if s.port != nil {
		return s.port.Close()
	}
	return nil
}
