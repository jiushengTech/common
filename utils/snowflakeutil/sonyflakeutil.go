package snowflakeutil

import (
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
		panic(err)
	}
	sf = s
}

// InitSnowflake deprecated plz use idutil.GetId
func GetId() uint64 {
	id, err := sf.NextID()
	if err != nil {
		return uint64(time.Now().UnixNano())
	}
	return id
}
