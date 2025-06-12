package idutil

import (
	"github.com/sony/sonyflake/v2"
	"time"
)

// Option 是一个函数类型，用于配置 Sonyflake 设置
type Option func(*sonyflake.Settings)

// WithStartTime 设置 StartTime 选项
func WithStartTime(startTime time.Time) Option {
	return func(settings *sonyflake.Settings) {
		settings.StartTime = startTime
	}
}

// WithMachineID 设置 MachineID 选项
func WithMachineID(machineID func() (int, error)) Option {
	return func(settings *sonyflake.Settings) {
		settings.MachineID = machineID
	}
}

// WithCheckMachineID 设置 CheckMachineID 选项
func WithCheckMachineID(checkMachineID func(int) bool) Option {
	return func(settings *sonyflake.Settings) {
		settings.CheckMachineID = checkMachineID
	}
}
