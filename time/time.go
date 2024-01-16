package main

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type LocalTime time.Time

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func (t *LocalTime) Value() (driver.Value, error) {
	tlt := time.Time(*t)
	// 使用 IsZero 方法判断时间是否为零时
	if tlt.IsZero() {
		return nil, nil
	}
	return tlt, nil
}

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
