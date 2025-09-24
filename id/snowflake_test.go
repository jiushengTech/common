package id

import "testing"

func TestGetId(t *testing.T) {
	snowflake, err := NewSnowflake()
	if err != nil {
		panic(err)
	}
	for _ = range 1000 {
		println(snowflake.NextID())
	}
}
