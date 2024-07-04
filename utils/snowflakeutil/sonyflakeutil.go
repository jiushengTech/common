package snowflakeutil

import (
	"github.com/jiushengTech/common/log"
	"github.com/sony/sonyflake"
	"time"
)

var sf *sonyflake.Sonyflake

// InitSnowflake deprecated plz use idutil.InitSnowflake
func InitSnowflake() {
	settings := sonyflake.Settings{
		StartTime: time.Now(),
	}
	s, err := sonyflake.New(settings)
	if err != nil {
		log.Info("Snowflake 初始化失败")
		panic(err)
	}
	sf = s
	log.Info("Snowflake 初始化成功")
}

// InitSnowflake deprecated plz use idutil.GetId
func GetId() uint64 {
	id, err := sf.NextID()
	if err != nil {
		log.Error("雪花id获取失败，使用当前时间戳代替")
		return uint64(time.Now().UnixNano())
	}
	return id
}
