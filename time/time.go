package time

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jiushengTech/common/log"
	"time"
	"unicode"
)

type LocalTime time.Time

// MarshalJSON implements the json.Marshaler interface.
func (t *LocalTime) MarshalJSON() ([]byte, error) {
	localTime := time.Time(*t)
	return json.Marshal(localTime.Format(time.DateTime))
}

// Value 实现 Value 方法，写入数据库时会调用该方法将自定义时间类型转换并写入数据库
func (t LocalTime) Value() (driver.Value, error) {
	localTime := time.Time(t)
	// 如果时间戳为零值，则返回 nil
	if localTime.IsZero() {
		return nil, nil
	}
	format := localTime.Format(time.DateTime)
	parse, err := time.Parse(time.DateTime, format)
	if err != nil {
		return nil, err
	}
	return parse, nil
}

// Scan 实现 Scan 方法，读取数据库时会调用该方法将时间数据转换成自定义时间类型；
func (t *LocalTime) Scan(v interface{}) error {
	switch value := v.(type) {
	case time.Time:
		parse, err := time.Parse(time.DateTime, value.Format(time.DateTime))
		if err != nil {
			return err
		}
		*t = LocalTime(parse)
		return nil
	case nil:
		*t = LocalTime(time.Time{})
		return nil
	case []uint8:
		parse, err := time.Parse(time.DateTime, string(value))
		if err != nil {
			return err
		}
		*t = LocalTime(parse)
		return nil
	default:
		return fmt.Errorf("can not convert %v to timestamp", v)
	}
}

func (t *LocalTime) String() string {
	return time.Time(*t).Format(time.DateTime)
}

func (t *LocalTime) ConvertTime() time.Time {
	return time.Time(*t)
}

func Now() LocalTime {
	format := time.Now().Format(time.DateTime)
	parse, err := time.Parse(time.DateTime, format)
	if err != nil {
		panic(err)
	}
	return LocalTime(parse)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	// 去除字母并用空格替换
	var cleanedStr string
	for _, char := range timeStr {
		if unicode.IsDigit(char) || unicode.IsPunct(char) {
			cleanedStr += string(char)
		} else {
			cleanedStr += " "
		}
	}
	// 截取固定长度的字符串
	const layoutLength = 19
	if len(cleanedStr) < layoutLength {
		// 打印错误信息并将时间设置为空
		log.Info("Invalid time string format")
		*t = LocalTime(time.Time{})
		return nil
	}
	cleanedStr = cleanedStr[:layoutLength]
	parseTime, err := time.Parse(time.DateTime, cleanedStr)
	if err != nil {
		return err
	}
	*t = LocalTime(parseTime)
	return nil
}
