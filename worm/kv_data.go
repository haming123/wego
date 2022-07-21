package worm

import "time"

//KeyVal use int64 as key
type KeyVal4IntInt struct {
	Key int64
	Val int64
}

type KeyVal4IntFloat struct {
	Key int64
	Val float64
}

type KeyVal4IntTime struct {
	Key int64
	Val time.Time
}

type KeyVal4IntString struct {
	Key int64
	Val string
}

//KeyVal use string as key
type KeyVal4StringInt struct {
	Key string
	Val int64
}

type KeyVal4StringFloat struct {
	Key string
	Val float64
}

type KeyVal4StringString struct {
	Key string
	Val string
}

type KeyVal4StringTime struct {
	Key string
	Val time.Time
}
