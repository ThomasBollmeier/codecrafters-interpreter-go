package main

import "errors"

type Environment struct {
	parent *Environment
	values map[string]Value
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent: parent,
		values: make(map[string]Value),
	}
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
