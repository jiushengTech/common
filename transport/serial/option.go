package serial

import "github.com/jacobsa/go-serial/serial"

type Option func(o *Server)

func WithPortName(portName string) Option {
	return func(s *Server) {
		s.PortName = portName
	}
}

func WithBaudRate(baudRate uint) Option {
	return func(s *Server) {
		s.BaudRate = baudRate
	}
}

func WithDataBits(dataBits uint) Option {
	return func(s *Server) {
		s.DataBits = dataBits
	}
}

func WithStopBits(stopBits uint) Option {
	return func(s *Server) {
		s.StopBits = stopBits
	}
}

func WithParityMode(parityMode serial.ParityMode) Option {
	return func(s *Server) {
		s.ParityMode = parityMode
	}
}

func WithRTSCTSFlowControl(rtsCtsFlowControl bool) Option {
	return func(s *Server) {
		s.RTSCTSFlowControl = rtsCtsFlowControl
	}
}

func WithInterCharacterTimeout(interCharacterTimeout uint) Option {
	return func(s *Server) {
		s.InterCharacterTimeout = interCharacterTimeout
	}
}

func WithMinimumReadSize(minimumReadSize uint) Option {
	return func(s *Server) {
		s.MinimumReadSize = minimumReadSize
	}
}

func WithRs485Enable(rs485Enable bool) Option {
	return func(s *Server) {
		s.Rs485Enable = rs485Enable
	}
}

func WithRs485RtsHighDuringSend(rtsHighDuringSend bool) Option {
	return func(s *Server) {
		s.Rs485RtsHighDuringSend = rtsHighDuringSend
	}
}

func WithRs485RtsHighAfterSend(rtsHighAfterSend bool) Option {
	return func(s *Server) {
		s.Rs485RtsHighAfterSend = rtsHighAfterSend
	}
}

func WithRs485RxDuringTx(rxDuringTx bool) Option {
	return func(s *Server) {
		s.Rs485RxDuringTx = rxDuringTx
	}
}

func WithRs485DelayRtsBeforeSend(delayRtsBeforeSend int) Option {
	return func(s *Server) {
		s.Rs485DelayRtsBeforeSend = delayRtsBeforeSend
	}
}

func WithRs485DelayRtsAfterSend(delayRtsAfterSend int) Option {
	return func(s *Server) {
		s.Rs485DelayRtsAfterSend = delayRtsAfterSend
	}
}
