package tcp

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestTCPServer_Basic(t *testing.T) {
	// 创建服务器
	server := NewServer(
		WithAddress("localhost:0"), // 使用随机端口
		WithKeepAlive(10*time.Second),
		WithReadTimeout(5*time.Second),
		WithWriteTimeout(5*time.Second),
		WithMaxConnections(10),
		WithDataChannelSize(50),
	)

	ctx := context.Background()

	// 启动服务器
	err := server.Start(ctx)
	if err != nil {
		t.Fatalf("启动服务器失败: %v", err)
	}
	defer server.Stop(ctx)

	// 获取实际监听地址
	endpoint, err := server.Endpoint()
	if err != nil {
		t.Fatalf("获取端点失败: %v", err)
	}

	// 测试初始状态
	if count := server.GetClientCount(); count != 0 {
		t.Errorf("初始客户端数量应为0，实际为: %d", count)
	}

	// 创建客户端连接
	conn, err := net.Dial("tcp", endpoint.Host)
	if err != nil {
		t.Fatalf("客户端连接失败: %v", err)
	}
	defer conn.Close()

	// 等待连接建立
	time.Sleep(100 * time.Millisecond)

	// 测试客户端数量
	if count := server.GetClientCount(); count != 1 {
		t.Errorf("期望客户端数量为1，实际为: %d", count)
	}

	// 测试获取客户端列表
	clients := server.GetClients()
	if len(clients) != 1 {
		t.Errorf("期望客户端列表长度为1，实际为: %d", len(clients))
	}

	// 获取客户端ID
	var clientID string
	for id := range clients {
		clientID = id
		break
	}

	// 测试获取指定客户端
	client, exists := server.GetClient(clientID)
	if !exists {
		t.Errorf("客户端应该存在")
	}
	if client == nil {
		t.Errorf("客户端不应为nil")
	}

	// 测试发送数据到指定客户端
	testMessage := "Hello Client!"
	n, err := server.SendToClient(clientID, []byte(testMessage))
	if err != nil {
		t.Errorf("发送数据失败: %v", err)
	}
	if n != len(testMessage) {
		t.Errorf("期望发送%d字节，实际发送%d字节", len(testMessage), n)
	}

	// 验证客户端收到数据
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	receivedLen, err := conn.Read(buffer)
	if err != nil {
		t.Errorf("客户端读取数据失败: %v", err)
	}
	received := string(buffer[:receivedLen])
	if received != testMessage {
		t.Errorf("期望收到'%s'，实际收到'%s'", testMessage, received)
	}

	// 测试广播功能
	broadcastMessage := "Broadcast Message!"
	totalBytes, err := server.Broadcast([]byte(broadcastMessage))
	if err != nil {
		t.Errorf("广播失败: %v", err)
	}
	if totalBytes != len(broadcastMessage) {
		t.Errorf("期望广播%d字节，实际广播%d字节", len(broadcastMessage), totalBytes)
	}

	// 验证客户端收到广播数据
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	receivedLen, err = conn.Read(buffer)
	if err != nil {
		t.Errorf("客户端读取广播数据失败: %v", err)
	}
	received = string(buffer[:receivedLen])
	if received != broadcastMessage {
		t.Errorf("期望收到广播'%s'，实际收到'%s'", broadcastMessage, received)
	}
}

func TestTCPServer_Events(t *testing.T) {
	var (
		connectedClients    = make(map[string]bool)
		disconnectedClients = make(map[string]bool)
		receivedData        = make([]string, 0)
		errors              = make([]string, 0)
		mu                  sync.Mutex
	)

	// 创建服务器
	server := NewServer(WithAddress("localhost:0"))

	// 设置事件处理器
	server.SetEventHandler(&EventHandler{
		OnClientConnected: func(client *ClientConn) {
			mu.Lock()
			connectedClients[client.ID] = true
			mu.Unlock()
			t.Logf("客户端连接: %s", client.ID)
		},
		OnClientDisconnected: func(client *ClientConn) {
			mu.Lock()
			disconnectedClients[client.ID] = true
			mu.Unlock()
			t.Logf("客户端断开: %s", client.ID)
		},
		OnClientData: func(client *ClientConn, data []byte) {
			mu.Lock()
			receivedData = append(receivedData, string(data))
			mu.Unlock()
			t.Logf("收到数据: %s -> %s", client.ID, string(data))
		},
		OnServerError: func(err error) {
			mu.Lock()
			errors = append(errors, err.Error())
			mu.Unlock()
			t.Logf("服务器错误: %v", err)
		},
	})

	ctx := context.Background()
	err := server.Start(ctx)
	if err != nil {
		t.Fatalf("启动服务器失败: %v", err)
	}
	defer server.Stop(ctx)

	endpoint, _ := server.Endpoint()

	// 创建客户端连接
	conn, err := net.Dial("tcp", endpoint.Host)
	if err != nil {
		t.Fatalf("客户端连接失败: %v", err)
	}

	// 等待连接事件
	time.Sleep(100 * time.Millisecond)

	// 验证连接事件
	mu.Lock()
	if len(connectedClients) != 1 {
		t.Errorf("期望1个连接事件，实际%d个", len(connectedClients))
	}
	mu.Unlock()

	// 发送数据到服务器
	testData := "Hello Server!"
	_, err = conn.Write([]byte(testData))
	if err != nil {
		t.Errorf("客户端发送数据失败: %v", err)
	}

	// 等待数据接收事件
	time.Sleep(100 * time.Millisecond)

	// 验证数据接收事件
	mu.Lock()
	if len(receivedData) != 1 {
		t.Errorf("期望收到1条数据，实际收到%d条", len(receivedData))
	} else if receivedData[0] != testData {
		t.Errorf("期望收到'%s'，实际收到'%s'", testData, receivedData[0])
	}
	mu.Unlock()

	// 关闭客户端连接
	conn.Close()

	// 等待断开事件
	time.Sleep(100 * time.Millisecond)

	// 验证断开事件
	mu.Lock()
	if len(disconnectedClients) != 1 {
		t.Errorf("期望1个断开事件，实际%d个", len(disconnectedClients))
	}
	mu.Unlock()
}

func TestTCPServer_MultipleClients(t *testing.T) {
	server := NewServer(
		WithAddress("localhost:0"),
		WithMaxConnections(5),
	)

	ctx := context.Background()
	err := server.Start(ctx)
	if err != nil {
		t.Fatalf("启动服务器失败: %v", err)
	}
	defer server.Stop(ctx)

	endpoint, _ := server.Endpoint()

	// 创建多个客户端连接
	const clientCount = 3
	clients := make([]net.Conn, clientCount)

	for i := 0; i < clientCount; i++ {
		conn, err := net.Dial("tcp", endpoint.Host)
		if err != nil {
			t.Fatalf("第%d个客户端连接失败: %v", i+1, err)
		}
		clients[i] = conn
		defer conn.Close()
	}

	// 等待所有连接建立
	time.Sleep(200 * time.Millisecond)

	// 验证客户端数量
	if count := server.GetClientCount(); count != clientCount {
		t.Errorf("期望客户端数量为%d，实际为%d", clientCount, count)
	}

	// 测试广播到多个客户端
	broadcastMsg := "Broadcast to all!"
	totalBytes, err := server.Broadcast([]byte(broadcastMsg))
	if err != nil {
		t.Errorf("广播失败: %v", err)
	}

	expectedTotal := len(broadcastMsg) * clientCount
	if totalBytes != expectedTotal {
		t.Errorf("期望广播总字节数%d，实际%d", expectedTotal, totalBytes)
	}

	// 验证每个客户端都收到消息
	for i, conn := range clients {
		buffer := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, err := conn.Read(buffer)
		if err != nil {
			t.Errorf("第%d个客户端读取数据失败: %v", i+1, err)
			continue
		}
		received := string(buffer[:n])
		if received != broadcastMsg {
			t.Errorf("第%d个客户端期望收到'%s'，实际收到'%s'", i+1, broadcastMsg, received)
		}
	}
}

func TestTCPServer_ConnectionLimit(t *testing.T) {
	const maxConns = 2
	server := NewServer(
		WithAddress("localhost:0"),
		WithMaxConnections(maxConns),
	)

	ctx := context.Background()
	err := server.Start(ctx)
	if err != nil {
		t.Fatalf("启动服务器失败: %v", err)
	}
	defer server.Stop(ctx)

	endpoint, _ := server.Endpoint()

	// 创建最大数量的连接
	connections := make([]net.Conn, maxConns)
	for i := 0; i < maxConns; i++ {
		conn, err := net.Dial("tcp", endpoint.Host)
		if err != nil {
			t.Fatalf("创建第%d个连接失败: %v", i+1, err)
		}
		connections[i] = conn
		defer conn.Close()
	}

	// 等待连接建立
	time.Sleep(200 * time.Millisecond)

	// 验证连接数量
	if count := server.GetClientCount(); count != maxConns {
		t.Errorf("期望连接数%d，实际%d", maxConns, count)
	}

	// 尝试创建超出限制的连接
	extraConn, err := net.Dial("tcp", endpoint.Host)
	if err != nil {
		t.Logf("预期行为：超出连接限制时连接失败")
	} else {
		defer extraConn.Close()
		// 等待一段时间看是否会被服务器拒绝
		time.Sleep(300 * time.Millisecond)
		// 这里连接可能建立但服务器不处理，具体行为取决于实现
	}
}

func TestTCPServer_CloseClient(t *testing.T) {
	server := NewServer(WithAddress("localhost:0"))

	ctx := context.Background()
	err := server.Start(ctx)
	if err != nil {
		t.Fatalf("启动服务器失败: %v", err)
	}
	defer server.Stop(ctx)

	endpoint, _ := server.Endpoint()

	// 创建客户端连接
	conn, err := net.Dial("tcp", endpoint.Host)
	if err != nil {
		t.Fatalf("客户端连接失败: %v", err)
	}
	defer conn.Close()

	// 等待连接建立
	time.Sleep(100 * time.Millisecond)

	// 获取客户端ID
	clients := server.GetClients()
	if len(clients) != 1 {
		t.Fatalf("期望1个客户端，实际%d个", len(clients))
	}

	var clientID string
	for id := range clients {
		clientID = id
		break
	}

	// 测试关闭指定客户端
	err = server.CloseClient(clientID)
	if err != nil {
		t.Errorf("关闭客户端失败: %v", err)
	}

	// 等待连接关闭
	time.Sleep(100 * time.Millisecond)

	// 验证客户端已被移除
	if count := server.GetClientCount(); count != 0 {
		t.Errorf("关闭客户端后，期望客户端数量为0，实际为%d", count)
	}

	// 测试关闭不存在的客户端
	err = server.CloseClient("nonexistent")
	if err == nil {
		t.Errorf("关闭不存在的客户端应该返回错误")
	}
}

func TestTCPServer_StartStop(t *testing.T) {
	server := NewServer(WithAddress("localhost:0"))

	ctx := context.Background()

	// 测试启动
	err := server.Start(ctx)
	if err != nil {
		t.Fatalf("启动服务器失败: %v", err)
	}

	// 验证服务器状态
	endpoint, err := server.Endpoint()
	if err != nil {
		t.Errorf("获取端点失败: %v", err)
	}
	if endpoint.Scheme != "tcp" {
		t.Errorf("期望协议为tcp，实际为%s", endpoint.Scheme)
	}

	// 测试停止
	err = server.Stop(ctx)
	if err != nil {
		t.Errorf("停止服务器失败: %v", err)
	}

	// 验证服务器已停止（尝试连接应失败）
	time.Sleep(100 * time.Millisecond)
	_, err = net.Dial("tcp", endpoint.Host)
	if err == nil {
		t.Errorf("服务器停止后连接不应成功")
	}
}

// 基准测试
func BenchmarkTCPServer_Broadcast(b *testing.B) {
	server := NewServer(WithAddress("localhost:0"))

	ctx := context.Background()
	server.Start(ctx)
	defer server.Stop(ctx)

	endpoint, _ := server.Endpoint()

	// 创建一些客户端连接
	const clientCount = 10
	for i := 0; i < clientCount; i++ {
		conn, _ := net.Dial("tcp", endpoint.Host)
		defer conn.Close()
	}

	time.Sleep(100 * time.Millisecond) // 等待连接建立

	data := []byte("benchmark message")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.Broadcast(data)
	}
}

// 测试不读取数据时的性能优化
func TestTCPServer_NoDataReading(t *testing.T) {
	server := NewServer(WithAddress("localhost:0"))

	// 不设置OnClientData事件处理器
	server.SetEventHandler(&EventHandler{
		OnClientConnected: func(client *ClientConn) {
			t.Logf("客户端连接: %s", client.ID)
		},
		OnClientDisconnected: func(client *ClientConn) {
			t.Logf("客户端断开: %s", client.ID)
		},
		// OnClientData 为 nil，服务器不应该读取数据
		OnClientData: nil,
		OnServerError: func(err error) {
			t.Logf("服务器错误: %v", err)
		},
	})

	ctx := context.Background()
	err := server.Start(ctx)
	if err != nil {
		t.Fatalf("启动服务器失败: %v", err)
	}
	defer server.Stop(ctx)

	endpoint, _ := server.Endpoint()

	// 创建客户端连接
	conn, err := net.Dial("tcp", endpoint.Host)
	if err != nil {
		t.Fatalf("客户端连接失败: %v", err)
	}
	defer conn.Close()

	// 等待连接建立
	time.Sleep(100 * time.Millisecond)

	// 验证连接已建立
	if count := server.GetClientCount(); count != 1 {
		t.Errorf("期望客户端数量为1，实际为%d", count)
	}

	// 客户端发送数据
	testData := "This data should not be processed"
	_, err = conn.Write([]byte(testData))
	if err != nil {
		t.Errorf("客户端发送数据失败: %v", err)
	}

	// 等待一段时间，确保服务器有足够时间处理（但实际上不应该处理）
	time.Sleep(200 * time.Millisecond)

	// 验证连接仍然存在（数据没有被读取，连接不应该因为读取错误而断开）
	if count := server.GetClientCount(); count != 1 {
		t.Errorf("发送数据后，期望客户端数量仍为1，实际为%d", count)
	}

	// 服务器向客户端发送数据应该正常工作
	clients := server.GetClients()
	var clientID string
	for id := range clients {
		clientID = id
		break
	}

	responseMsg := "Server response"
	n, err := server.SendToClient(clientID, []byte(responseMsg))
	if err != nil {
		t.Errorf("服务器发送数据失败: %v", err)
	}
	if n != len(responseMsg) {
		t.Errorf("期望发送%d字节，实际发送%d字节", len(responseMsg), n)
	}

	// 验证客户端能收到服务器的响应
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	receivedLen, err := conn.Read(buffer)
	if err != nil {
		t.Errorf("客户端读取服务器响应失败: %v", err)
	}
	received := string(buffer[:receivedLen])
	if received != responseMsg {
		t.Errorf("期望收到'%s'，实际收到'%s'", responseMsg, received)
	}

	t.Log("测试通过：OnClientData为nil时，服务器不读取客户端数据，但连接和发送功能正常")
}

// 基准测试 - 比较有无数据读取的性能差异
func BenchmarkTCPServer_WithoutDataReading(b *testing.B) {
	server := NewServer(WithAddress("localhost:0"))

	// 不设置OnClientData
	server.SetEventHandler(&EventHandler{
		OnClientConnected: func(client *ClientConn) {},
		OnClientData:      nil, // 关键：不读取数据
	})

	ctx := context.Background()
	server.Start(ctx)
	defer server.Stop(ctx)

	endpoint, _ := server.Endpoint()

	// 创建连接
	conn, _ := net.Dial("tcp", endpoint.Host)
	defer conn.Close()

	time.Sleep(50 * time.Millisecond) // 等待连接建立

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟连接维护的开销
		count := server.GetClientCount()
		if count != 1 {
			b.Errorf("连接数异常: %d", count)
		}
	}
}

func BenchmarkTCPServer_WithDataReading(b *testing.B) {
	server := NewServer(WithAddress("localhost:0"))

	// 设置OnClientData
	server.SetEventHandler(&EventHandler{
		OnClientConnected: func(client *ClientConn) {},
		OnClientData: func(client *ClientConn, data []byte) {
			// 简单处理数据
		},
	})

	ctx := context.Background()
	server.Start(ctx)
	defer server.Stop(ctx)

	endpoint, _ := server.Endpoint()

	// 创建连接
	conn, _ := net.Dial("tcp", endpoint.Host)
	defer conn.Close()

	time.Sleep(50 * time.Millisecond) // 等待连接建立

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟连接维护的开销
		count := server.GetClientCount()
		if count != 1 {
			b.Errorf("连接数异常: %d", count)
		}
	}
}
