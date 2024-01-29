package time

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type LocalTime time.Time

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	localTime := time.Time(*t)
	return json.Marshal(localTime.Format(time.DateTime))
}

func (t *LocalTime) Value() (driver.Value, error) {
	localTime := time.Time(*t)
	// 如果时间戳为零值，则返回 nil
	if localTime.IsZero() {
		return nil, nil
	}
	return localTime, nil
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

func (t *LocalTime) String() string {
	return time.Time(*t).Format(time.DateTime)
}
