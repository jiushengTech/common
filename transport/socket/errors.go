package socket

import (
	"errors"
	"fmt"
)

var (
	// ErrServerClosed 服务已关闭
	ErrServerClosed = errors.New("socket: 服务已关闭")
	// ErrInvalidAddress 无效地址
	ErrInvalidAddress = errors.New("socket: 无效地址")
	// ErrUnsupportedNetwork 不支持的网络类型
	ErrUnsupportedNetwork = errors.New("socket: 不支持的网络类型")
	// ErrConnectionFailed 连接失败
	ErrConnectionFailed = errors.New("socket: 连接失败")
	// ErrWriteFailed 写入失败
	ErrWriteFailed = errors.New("socket: 写入失败")
	// ErrReadFailed 读取失败
	ErrReadFailed = errors.New("socket: 读取失败")
	// ErrPoolExhausted 连接池耗尽
	ErrPoolExhausted = errors.New("socket: 连接池已耗尽")
)

// NetworkError 网络错误
type NetworkError struct {
	Op  string // 操作
	Err error  // 原始错误
}

func (e *NetworkError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("socket: %s: %v", e.Op, e.Err)
	}
	return fmt.Sprintf("socket: %s", e.Op)
}

// Unwrap 解包错误
func (e *NetworkError) Unwrap() error {
	return e.Err
}
