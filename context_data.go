package wego

import (
	"errors"
	"time"
)

var errWrongType = errors.New("incorrect type")

type ContextData struct {
	data	map[string]interface{}
}

func (c *ContextData) reset() {
	c.data = nil
}

func (c *ContextData) Set(key string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = value
}

func (c *ContextData) Get(key string) (interface{}, bool) {
	value, exists := c.data[key]
	return value, exists
}

func (c *ContextData) GetString(key string) (string, error) {
	data, ok := c.Get(key)
	if !ok {
		return "", errNotFind
	}
	val, ok := data.(string)
	if !ok {
		return "", errWrongType
	}
	return val, nil
}

func (c *ContextData) GetBool(key string) (bool, error) {
	data, ok := c.Get(key)
	if !ok {
		return false, errNotFind
	}
	val, ok := data.(bool)
	if !ok {
		return false, errWrongType
	}
	return val, nil
}

func (c *ContextData) GetTime(key string) (time.Time, error) {
	data, ok := c.Get(key)
	if !ok {
		return time.Time{}, errNotFind
	}
	val, ok := data.(time.Time)
	if !ok {
		return time.Time{}, errWrongType
	}
	return val, nil
}

func (s *ContextData) GetInt64(key string) (int64, error) {
	data, ok := s.Get(key)
	if !ok {
		return 0, errNotFind
	}
	switch val := data.(type) {
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return int64(val), nil
	case uint:
		return int64(val), nil
	case uint8:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint64:
		return int64(val), nil
	}
	return 0, errWrongType
}

func (s *ContextData) GetFloat64(key string) (float64, error) {
	data, ok := s.Get(key)
	if !ok {
		return 0, errNotFind
	}
	switch val := data.(type) {
	case float32:
		return float64(val), nil
	case float64:
		return float64(val), nil
	}
	return 0, errWrongType
}