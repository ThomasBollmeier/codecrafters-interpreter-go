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
	VtString
)

type Value interface {
	getType() ValueType
	isEqualTo(Value) bool
	isTruthy() bool
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

func (n *NumValue) isEqualTo(other Value) bool {
	if other == nil || other.getType() != n.getType() {
		return false
	}
	return n.Value == other.(*NumValue).Value
}

func (n *NumValue) isTruthy() bool {
	return n.Value != 0
}

func (n *NumValue) String() string {
	numStr := strings.TrimRight(fmt.Sprintf("%f", n.Value), "0")
	lastIdx := len(numStr) - 1
	if numStr[lastIdx] == uint8('.') {
		numStr = numStr[:lastIdx]
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

func (b *BooleanValue) isEqualTo(other Value) bool {
	if other == nil || other.getType() != b.getType() {
		return false
	}
	return b.Value == other.(*BooleanValue).Value
}

func (b *BooleanValue) isTruthy() bool {
	return b.Value
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

func (n *NilValue) isEqualTo(other Value) bool {
	if other == nil || other.getType() != n.getType() {
		return false
	}
	return true
}

func (n *NilValue) isTruthy() bool {
	return false
}

func (n *NilValue) String() string {
	return "nil"
}

type StringValue struct {
	Value string
}

func NewStringValue(v string) *StringValue {
	return &StringValue{v}
}

func (s *StringValue) getType() ValueType {
	return VtString
}

func (s *StringValue) isEqualTo(other Value) bool {
	if other == nil || other.getType() != s.getType() {
		return false
	}
	return s.Value == other.(*StringValue).Value
}

func (s *StringValue) isTruthy() bool {
	return true
}

func (s *StringValue) String() string {
	return s.Value
}
