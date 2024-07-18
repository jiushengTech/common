package idutil

import "testing"

func TestGetId(t *testing.T) {
	InitSnowflake()
	for _ = range 1000 {
		println(GetId())
	}
}
