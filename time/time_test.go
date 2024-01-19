package time

import (
	"fmt"
	"testing"
	"time"
)

func TestLocalTime(t *testing.T) {
	lt := LocalTime(time.Now())
	// 示例：将 LocalTime 转换为 JSON 格式
	jsonData, err := lt.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonData))

	// 示例：将 LocalTime 转换为数据库 Value
	dbValue, err := lt.Value()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("db==>", dbValue)

	// 示例：从数据库 Scan 到 LocalTime
	var ltScan LocalTime
	err = ltScan.Scan(time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("scan==>", ltScan)
}
