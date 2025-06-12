package idutil

import (
	"fmt"
	"github.com/sony/sonyflake/v2"
	"math/rand"
	"time"
)

var sf *sonyflake.Sonyflake

// InitSnowflake 初始化 Sonyflake 并应用传递的选项
func InitSnowflake(opts ...Option) (*sonyflake.Sonyflake, error) {
	settings := sonyflake.Settings{}
	// 应用所有传递的选项
	for _, opt := range opts {
		opt(&settings)
	}

	s, err := sonyflake.New(settings)
	if err != nil {
		return nil, err
	}
	sf = s
	return s, nil
}

func GetId() int64 {
	var id int64
	var err error
	for i := 0; i < 3; i++ {
		id, err = sf.NextID()
		if err == nil {
			return id
		}
		fmt.Printf("雪花ID获取失败，重试次数: %d，错误: %v\n", i+1, err)
		time.Sleep(1 * time.Millisecond)
	}

	fmt.Println("雪花ID获取失败，已重试3次，使用随机ID替代")
	// 返回64位正整数
	return rand.Int63()
}
