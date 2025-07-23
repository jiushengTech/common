package tcp

// EventHandler 事件处理器
type EventHandler struct {
	OnClientConnected    func(client *ClientConn)
	OnClientDisconnected func(client *ClientConn)
	OnClientData         func(client *ClientConn, data []byte)
	OnServerError        func(err error)
}

// SetEventHandler 设置事件处理器
func (s *Server) SetEventHandler(handler *EventHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if handler != nil {
		s.handler = handler
	}
}

// triggerClientConnected 触发客户端连接事件
func (s *Server) triggerClientConnected(client *ClientConn) {
	if s.handler != nil && s.handler.OnClientConnected != nil {
		s.handler.OnClientConnected(client)
	}
}

// triggerClientDisconnected 触发客户端断开事件
func (s *Server) triggerClientDisconnected(client *ClientConn) {
	if s.handler != nil && s.handler.OnClientDisconnected != nil {
		s.handler.OnClientDisconnected(client)
	}
}

// triggerClientData 触发客户端数据事件
func (s *Server) triggerClientData(client *ClientConn, data []byte) {
	if s.handler != nil && s.handler.OnClientData != nil {
		s.handler.OnClientData(client, data)
	}
}

// triggerServerError 触发服务器错误事件
func (s *Server) triggerServerError(err error) {
	if s.handler != nil && s.handler.OnServerError != nil {
		s.handler.OnServerError(err)
	}
}
