package worm

import (
	"testing"
)

func NewEngine4Test() (*DbEngine, error) {
	_, err := OpenDb()
	if err != nil {
		return nil, err
	}

	eng, err := NewEngine(&dialectMysql{}, DbConn)
	if err != nil {
		return nil, err
	}
	return eng, nil
}

func InitEngine4Test() error {
	_, err := OpenDb()
	if err != nil {
		return err
	}

	err = InitEngine(&dialectMysql{}, DbConn)
	if err != nil {
		return err
	}
	return nil
}

func TestNewEngine(t *testing.T) {
	_, err := NewEngine4Test()
	if err != nil {
		t.Error(err)
	}
}

func TestInitEngine(t *testing.T)  {
	err := InitEngine4Test()
	if err != nil {
		t.Error(err)
	}
}

