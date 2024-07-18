package idutil

import (
	"github.com/jiushengTech/common/log"
	"github.com/sony/sonyflake"
	"time"
)

var sf *sonyflake.Sonyflake

// InitSnowflake 初始化 Sonyflake 并应用传递的选项
func InitSnowflake(opts ...Option) *sonyflake.Sonyflake {
	settings := sonyflake.Settings{}
	// 应用所有传递的选项
	for _, opt := range opts {
		opt(&settings)
	}

	s, err := sonyflake.New(settings)
	if err != nil {
		log.Fatal("Sonyflake 初始化失败: ", err)
	}
	log.Info("Sonyflake 初始化成功")
	return s
}

func GetId() int64 {
	var id uint64
	var err error
	for i := 0; i < 3; i++ {
		id, err = sf.NextID()
		if err == nil {
			return int64(id)
		}
		log.Errorf("雪花id获取失败，重试次数: %d", i+1)
		time.Sleep(1 * time.Millisecond)
	}
	log.Error("雪花id获取失败，已重试3次")
	return -1 // 失败时返回一个默认值，视情况而定
}

func GetUId() uint64 {
	var id uint64
	var err error
	for i := 0; i < 3; i++ {
		id, err = sf.NextID()
		if err == nil {
			return id
		}
		log.Errorf("雪花id获取失败，重试次数: %d", i+1)
		time.Sleep(1 * time.Millisecond)
	}
	log.Error("雪花id获取失败，已重试3次")
	return 0 // 失败时返回一个默认值，视情况而定
}
