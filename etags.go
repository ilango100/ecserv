package main

import (
	"time"
)

func tTag(t time.Time) string {
	buf := make([]byte, 3, 3)
	buf[0] = 64 + byte(t.Day())
	buf[1] = 97 + byte(t.Hour())
	buf[2] = 65 + byte(t.Minute())
	return string(buf)
}
