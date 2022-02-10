package worm

import "testing"

func TestSimpleLogger (t *testing.T) {
	log := NewSimpleLogger()
	log.Debug("hello debug message")
	log.Info("hello info message")
}
