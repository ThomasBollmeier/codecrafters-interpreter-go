package main

import (
	"fmt"
	"strings"
)

type ValueType int

const (
	VtNumber ValueType = iota
	VtBoolean
	VtNil
)

type Value interface {
	getType() ValueType
}

type NumValue struct {
	Value float64
}

func NewNumValue(v float64) *NumValue {
	return &NumValue{v}
}

func (n *NumValue) getType() ValueType {
	return VtNumber
}

func (n *NumValue) String() string {
	numStr := strings.TrimRight(fmt.Sprintf("%f", n.Value), "0")
	if numStr[len(numStr)-1] == uint8('.') {
		numStr = numStr + "0"
	}

	return numStr
}

type BooleanValue struct {
	Value bool
}

func NewBooleanValue(v bool) *BooleanValue {
	return &BooleanValue{v}
}

func (b *BooleanValue) getType() ValueType {
	return VtBoolean
}

func (b *BooleanValue) String() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}

type NilValue struct{}

func NewNilValue() *NilValue {
	return &NilValue{}
}

func (n *NilValue) getType() ValueType {
	return VtNil
}

func (n *NilValue) String() string {
	return "nil"
}
