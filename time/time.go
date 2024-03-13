package time

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
	"unicode"
)

type LocalTime time.Time

// MarshalJSON 重写 MarshaJSON 方法，在此方法中实现自定义格式的转换
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
	return localTime, nil
}

// Scan 实现 Scan 方法，读取数据库时会调用该方法将时间数据转换成自定义时间类型；
func (t *LocalTime) Scan(v interface{}) error {
	switch value := v.(type) {
	case time.Time:
		*t = LocalTime(value)
		return nil
	case nil:
		*t = LocalTime(time.Time{})
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
	parseTime, err := time.Parse(time.DateTime, cleanedStr)
	if err != nil {
		return err
	}
	*t = LocalTime(parseTime)
	return nil
}
