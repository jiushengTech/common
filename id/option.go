package id

import (
	"time"

	"github.com/sony/sonyflake/v2"
)

// Option 是一个函数类型，用于配置 Sonyflake 设置
type Option func(*sonyflake.Settings)

// WithBitsSequence 设置序列号位数（默认 8，最大 30）
func WithBitsSequence(bits int) Option {
	return func(settings *sonyflake.Settings) {
		settings.BitsSequence = bits
	}
}

// WithBitsMachineID 设置机器 ID 位数（默认 16，最大 30）
func WithBitsMachineID(bits int) Option {
	return func(settings *sonyflake.Settings) {
		settings.BitsMachineID = bits
	}
}

// WithTimeUnit 设置时间单位（默认 10ms，最小 1ms）
func WithTimeUnit(unit time.Duration) Option {
	return func(settings *sonyflake.Settings) {
		settings.TimeUnit = unit
	}
}

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
