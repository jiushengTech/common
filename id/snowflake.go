package id

import (
	"github.com/sony/sonyflake/v2"
)

// NewSnowflake 新建 Sonyflake 并应用传递的选项
func NewSnowflake(opts ...Option) (*sonyflake.Sonyflake, error) {
	settings := sonyflake.Settings{}
	// 应用所有传递的选项
	for _, opt := range opts {
		opt(&settings)
	}

	s, err := sonyflake.New(settings)
	if err != nil {
		return nil, err
	}
	return s, nil
}
