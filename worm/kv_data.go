package worm

import "time"

type KeyVal4IntInt struct {
	Key int64
	Val int64
}

type KeyVal4IntFloat struct {
	Key int64
	Val float64
}

type KeyVal4IntString struct {
	Key int64
	Val string
}

type KeyVal4IntTime struct {
	Key int64
	Val time.Time
}

type KeyVal4StringString struct {
	Key string
	Val string
}
