package utils

import "time"

func Sleep(second int) {
	time.Sleep(time.Duration(second) * time.Second)
}

func SleepShort(millisecond int) {
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
}
