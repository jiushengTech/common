package socket

import (
	"fmt"
	"time"
)

// 工厂方法统一创建服务器实例

// CreateServer 根据网络类型创建适当的服务器
// network: 网络类型，支持"tcp"或"udp"
// address: 服务器地址
// opts: 额外的配置选项
func CreateServer(network, address string, opts ...Option) (*Server, error) {
	// 创建基本选项
	baseOpts := []Option{
		WithNetwork(network),
		WithAddress(address),
	}

	// 合并选项
	allOpts := append(baseOpts, opts...)

	// 根据网络类型创建服务器
	srv := NewServer(allOpts...)

	if srv.address == "" {
		return nil, fmt.Errorf("无效的地址")
	}

	return srv, nil
}

// 以下是使用示例

// CreateTCPServer 创建TCP服务器的便捷方法
func CreateTCPServer(address string, opts ...Option) (*Server, error) {
	return CreateServer("tcp", address, opts...)
}

// CreateUDPServer 创建UDP服务器的便捷方法
func CreateUDPServer(address string, opts ...Option) (*Server, error) {
	return CreateServer("udp", address, opts...)
}

// Example 创建服务器的示例函数
func Example() {
	// 创建TCP服务器
	tcpServer, _ := CreateServer("tcp", "127.0.0.1:8080",
		WithTimeout(time.Second*10),
		WithMaxConns(200),
	)

	// 创建UDP服务器
	udpServer, _ := CreateServer("udp", "127.0.0.1:8081",
		WithReadDeadline(time.Second*5),
		WithTargetAddr([]string{"192.168.1.100:9000", "192.168.1.101:9000"}),
	)

	// 使用服务器
	_ = tcpServer
	_ = udpServer
}
