package id

import (
	"strconv"
	"time"

	"github.com/sony/sonyflake/v2"
)

// Snowflake 对 sonyflake.Sonyflake 进行面向对象封装
type Snowflake struct {
	flake *sonyflake.Sonyflake
}

// NewSnowflake 新建 Sonyflake 并应用传递的选项
func NewSnowflake(opts ...Option) (*Snowflake, error) {
	settings := sonyflake.Settings{}
	// 应用所有传递的选项
	for _, opt := range opts {
		opt(&settings)
	}

	s, err := sonyflake.New(settings)
	if err != nil {
		return nil, err
	}
	return &Snowflake{flake: s}, nil
}

// Uint64 获取一个 uint64 类型的唯一 ID，失败返回0
func (s *Snowflake) Uint64() uint64 {
	return uint64(s.Int64())
}

// Int64 获取一个 int64 类型的唯一 ID，内部重试3次，最终失败返回0
func (s *Snowflake) Int64() int64 {
	if s == nil || s.flake == nil {
		// snowflake未初始化直接返回0
		return 0
	}
	for i := 0; i < 3; i++ {
		id64, e := s.flake.NextID()
		if e == nil {
			return id64
		}
		// 可加短暂延迟再重试
		time.Sleep(10 * time.Millisecond)
	}

	// 最终失败返回0
	return 0
}

// String 获取一个 string 类型的唯一 ID，失败返回空字符串
func (s *Snowflake) String() string {
	return strconv.FormatInt(s.Int64(), 10)
}
