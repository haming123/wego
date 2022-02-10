package wego

import (
	"errors"
	"time"
)

var errNotFind = errors.New("not find")

type ParamItem struct {
	Key   	string
	Val 	string
}

type ValidString struct {
	Value 	string
	Error  	error
}

type ValidBool struct {
	Value 	bool
	Error  	error
}

type ValidInt struct {
	Value 	int
	Error  	error
}

type ValidInt32 struct {
	Value 	int32
	Error  	error
}

type ValidInt64 struct {
	Value 	int64
	Error  	error
}

type ValidFloat struct {
	Value 	float64
	Error  	error
}

type ValidTime struct {
	Value 	time.Time
	Error  	error
}
