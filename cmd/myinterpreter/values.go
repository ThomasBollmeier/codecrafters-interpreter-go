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
	VtBuiltinFunc
	VtLambda
	VtClass
	VtInstance
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

type callable interface {
	call(args []Value) (Value, error)
}

type BuiltinFuncValue struct {
	name string
	fn   func(args []Value) (Value, error)
}

func NewBuiltinFuncValue(name string, f func(args []Value) (Value, error)) *BuiltinFuncValue {
	return &BuiltinFuncValue{name, f}
}

func (b *BuiltinFuncValue) getType() ValueType {
	return VtBuiltinFunc
}

func (b *BuiltinFuncValue) isEqualTo(value Value) bool {
	fn, ok := value.(*BuiltinFuncValue)
	if !ok {
		return false
	}
	return b.name == fn.name
}

func (b *BuiltinFuncValue) isTruthy() bool {
	return true
}

func (b *BuiltinFuncValue) String() string {
	return fmt.Sprintf("<builtin-function %s<", b.name)
}

func (b *BuiltinFuncValue) call(args []Value) (Value, error) {
	return b.fn(args)
}

type LambdaValue struct {
	name       string
	parameters []string
	body       Block
	env        Environment
}

func NewLambdaValue(name string, parameters []string, body Block, env Environment) *LambdaValue {
	return &LambdaValue{name, parameters, body, env}
}

func (l *LambdaValue) getType() ValueType {
	return VtLambda
}

func (l *LambdaValue) isEqualTo(Value) bool {
	return false
}

func (l *LambdaValue) isTruthy() bool {
	return true
}

func (l *LambdaValue) String() string {
	return fmt.Sprintf("<fn %s>", l.name)
}

func (l *LambdaValue) call(args []Value) (Value, error) {
	if len(args) != len(l.parameters) {
		return nil, fmt.Errorf("expected %d args, got %d", len(l.parameters), len(args))
	}

	callEnv := NewEnvironment(&l.env)
	for i, param := range l.parameters {
		callEnv.Set(param, args[i])
	}

	interpreter := NewInterpreter(callEnv)

	interpreter.lambdaEvalActive = true
	interpreter.returnOccurred = false

	interpreter.visitBlock(&l.body)

	interpreter.lambdaEvalActive = false
	interpreter.returnOccurred = false

	return interpreter.lastResult, interpreter.lastError
}

type ClassValue struct {
	name    string
	methods []LambdaValue
}

func NewClassValue(name string, methods []LambdaValue) *ClassValue {
	return &ClassValue{name, methods}
}

func (c *ClassValue) getType() ValueType {
	return VtClass
}

func (c *ClassValue) isEqualTo(value Value) bool {
	cls, ok := value.(*ClassValue)
	if !ok {
		return false
	}
	return c.name == cls.name
}

func (c *ClassValue) isTruthy() bool {
	return true
}

func (c *ClassValue) String() string {
	return c.name
}

func (c *ClassValue) call([]Value) (Value, error) {
	return NewInstanceValue(c), nil
}

type InstanceValue struct {
	class      *ClassValue
	properties map[string]Value
}

func NewInstanceValue(class *ClassValue) *InstanceValue {
	return &InstanceValue{
		class:      class,
		properties: make(map[string]Value),
	}
}

func (i *InstanceValue) getType() ValueType {
	return VtInstance
}

func (i *InstanceValue) isEqualTo(value Value) bool {
	other, ok := value.(*InstanceValue)
	if !ok {
		return false
	}
	if !i.class.isEqualTo(other.class) {
		return false
	}

	return true
}

func (i *InstanceValue) isTruthy() bool {
	return true
}

func (i *InstanceValue) String() string {
	return fmt.Sprintf("%s instance", i.class.name)
}

func (i *InstanceValue) getProperty(path []string) (Value, error) {
	currInst := i
	for idx, property := range path {
		value, ok := currInst.properties[property]
		if !ok {
			return nil, fmt.Errorf("property %s not found", property)
		}
		if idx == len(path)-1 {
			return value, nil
		}
		currInst, ok = value.(*InstanceValue)
		if !ok {
			return nil, fmt.Errorf("expected value to be an instance, found %v", value)
		}
	}
	return nil, fmt.Errorf("property path must not be empty")
}

func (i *InstanceValue) setProperty(path []string, value Value) error {
	currInst := i
	for _, property := range path[:len(path)-1] {
		val, ok := currInst.properties[property]
		if !ok {
			return fmt.Errorf("property %s not found", property)
		}
		currInst, ok = val.(*InstanceValue)
		if !ok {
			return fmt.Errorf("expected value to be an instance, found %v", value)
		}
	}
	currInst.properties[path[len(path)-1]] = value
	return nil
}
