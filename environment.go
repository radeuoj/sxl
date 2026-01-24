package main

type Environment struct {
	store map[string]Value
}

func NewEnvironemnt() *Environment {
	s := make(map[string]Value)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (Value, bool) {
	val, ok := e.store[name]
	return val, ok
}

func (e *Environment) Set(name string, val Value) {
	e.store[name] = val
}
