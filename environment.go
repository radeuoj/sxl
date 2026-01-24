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

func (e *Environment) Let(name string, val Value) bool {
	if _, ok := e.store[name]; ok {
		return false
	} else {
		e.store[name] = val
		return true
	}
}

func (e *Environment) Assign(name string, val Value) bool {
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return true
	} else {
		return false
	}
}

func (e *Environment) set(name string, val Value) {
	e.store[name] = val
}
