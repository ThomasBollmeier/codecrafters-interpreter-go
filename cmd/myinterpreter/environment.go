package main

import "errors"

type Environment struct {
	parent *Environment
	values map[string]Value
}

func NewEnvironment(parent *Environment) *Environment {
	values := make(map[string]Value)

	if parent == nil {
		initBuiltins(values)
	}

	return &Environment{
		parent: parent,
		values: values,
	}
}

func initBuiltins(values map[string]Value) {
	values["clock"] = NewBuiltinFuncValue("clock", clock)
}

func (env *Environment) Get(name string) (Value, error) {
	value, ok := env.values[name]
	if ok {
		return value, nil
	}
	if env.parent != nil {
		return env.parent.Get(name)
	} else {
		return nil, errors.New("unknown identifier " + name)
	}
}

func (env *Environment) GetDefiningEnv(name string) (*Environment, error) {
	_, ok := env.values[name]
	if ok {
		return env, nil
	}
	if env.parent != nil {
		return env.parent.GetDefiningEnv(name)
	} else {
		return nil, errors.New("unknown identifier " + name)
	}
}

func (env *Environment) Set(name string, value Value) {
	env.values[name] = value
}
