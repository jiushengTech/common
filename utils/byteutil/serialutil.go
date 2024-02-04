package byteutil

import (
	"bytes"
	"encoding/binary"
)

// Marshal 将结构体序列化为字节切片
func Marshal[T any](data T, order binary.ByteOrder) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, order, &data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal 将字节切片反序列化为结构体
func Unmarshal[T any](data []byte, order binary.ByteOrder) (T, error) {
	var t T
	err := binary.Read(bytes.NewReader(data), order, &t)
	return t, err
}
