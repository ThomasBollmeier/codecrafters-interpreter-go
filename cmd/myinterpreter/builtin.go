package main

import (
	"errors"
	"time"
)

func clock(args []Value) (Value, error) {
	if len(args) != 0 {
		return nil, errors.New("clock() expects no arguments")
	}
	seconds := time.Now().Unix()
	return NewNumValue(float64(seconds)), nil
}
