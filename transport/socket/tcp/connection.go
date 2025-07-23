package tcp

import (
	"context"
	"fmt"
	"net"
	"time"
)

// acceptConnections 接受客户端连接
func (s *Server) acceptConnections(ctx context.Context) {
	for {
		s.mu.RLock()
		if !s.running || s.tcpListener == nil {
			s.mu.RUnlock()
			break
		}
		listener := s.tcpListener
		maxConns := s.config.MaxConnections
		currentConns := s.GetClientCount()
		s.mu.RUnlock()

		// 检查连接数限制
		if maxConns > 0 && currentConns >= maxConns {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// 设置接受超时
		if s.config.ConnectionTimeout > 0 {
			listener.SetDeadline(time.Now().Add(s.config.ConnectionTimeout))
		}

		conn, err := listener.AcceptTCP()
		if err != nil {
			if s.running {
				s.triggerServerError(fmt.Errorf("接受TCP连接失败: %w", err))
			}
			continue
		}

		// 配置连接选项
		if err := s.configureConnection(conn); err != nil {
			conn.Close()
			s.triggerServerError(fmt.Errorf("配置连接失败: %w", err))
			continue
		}

		// 添加到客户端列表
		client := newClientConn(conn)

		s.clients.Store(client.ID, client)

		// 触发连接事件
		s.triggerClientConnected(client)

		// 启动处理连接的协程
		go s.handleConnection(ctx, client)
	}
}

// configureConnection 配置TCP连接选项
func (s *Server) configureConnection(conn *net.TCPConn) error {
	if s.config.EnableKeepalive && s.config.KeepAlive > 0 {
		if err := conn.SetKeepAlive(true); err != nil {
			return err
		}
		if err := conn.SetKeepAlivePeriod(s.config.KeepAlive); err != nil {
			return err
		}
	}

	if s.config.EnableTCPNoDelay {
		if err := conn.SetNoDelay(true); err != nil {
			return err
		}
	}

	if s.config.ReadTimeout > 0 {
		if err := conn.SetReadDeadline(time.Now().Add(s.config.ReadTimeout)); err != nil {
			return err
		}
	}

	if s.config.WriteTimeout > 0 {
		if err := conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout)); err != nil {
			return err
		}
	}

	// 设置缓冲区大小
	if s.config.SendBufferSize > 0 {
		if err := conn.SetWriteBuffer(s.config.SendBufferSize); err != nil {
			return err
		}
	}

	if s.config.ReceiveBufferSize > 0 {
		if err := conn.SetReadBuffer(s.config.ReceiveBufferSize); err != nil {
			return err
		}
	}

	return nil
}

// handleConnection 维护客户端连接状态
func (s *Server) handleConnection(ctx context.Context, client *ClientConn) {
	defer func() {
		// 连接断开时从客户端列表中移除
		s.clients.Delete(client.ID)

		// 触发断开事件
		s.triggerClientDisconnected(client)

		client.Close()
	}()

	// 检查是否需要读取客户端数据
	needDataReading := s.handler != nil && s.handler.OnClientData != nil

	if needDataReading {
		// 需要读取数据的情况
		s.handleConnectionWithDataReading(ctx, client)
	} else {
		// 不需要读取数据的情况，只维护连接
		s.handleConnectionWithoutDataReading(ctx, client)
	}
}

// handleConnectionWithDataReading 处理需要读取数据的连接
func (s *Server) handleConnectionWithDataReading(ctx context.Context, client *ClientConn) {
	// 创建缓冲区用于读取数据
	bufferSize := s.config.ReadBufferSize
	if s.config.ReceiveBufferSize > 0 {
		bufferSize = s.config.ReceiveBufferSize
	}
	buffer := make([]byte, bufferSize)

	// 启动数据读取协程
	dataChan := make(chan []byte, s.config.DataChannelSize)
	errorChan := make(chan error, 1)

	go s.readClientData(client, buffer, dataChan, errorChan)

	// 连接保活和数据处理循环
	ticker := time.NewTicker(s.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !s.isServerRunning() || !client.IsAlive() {
				return
			}
		case data := <-dataChan:
			s.triggerClientData(client, data)
		case err := <-errorChan:
			s.triggerServerError(fmt.Errorf("客户端 %s 读取错误: %w", client.ID, err))
			return
		}
	}
}

// handleConnectionWithoutDataReading 处理不需要读取数据的连接
func (s *Server) handleConnectionWithoutDataReading(ctx context.Context, client *ClientConn) {
	// 连接保活循环（不读取数据）
	ticker := time.NewTicker(s.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !s.isServerRunning() || !client.IsAlive() {
				return
			}
		}
	}
}

// isServerRunning 检查服务器是否还在运行
func (s *Server) isServerRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// readClientData 读取客户端数据的协程
func (s *Server) readClientData(client *ClientConn, buffer []byte, dataChan chan []byte, errorChan chan error) {
	defer close(dataChan)
	defer close(errorChan)

	for {
		// 设置读取超时
		if s.config.ReadTimeout > 0 {
			client.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))
		}

		n, err := client.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 读取超时，继续等待
				continue
			}
			// 其他错误，发送错误信号
			select {
			case errorChan <- err:
			default:
			}
			return
		}

		if n > 0 {
			// 更新最后活跃时间
			client.mu.Lock()
			client.lastSeen = time.Now()
			client.mu.Unlock()

			// 复制数据并发送到数据通道
			data := make([]byte, n)
			copy(data, buffer[:n])

			select {
			case dataChan <- data:
			default:
				// 通道满了，丢弃数据（可以根据需求调整策略）
			}
		}
	}
}
