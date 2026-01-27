package main

type Environment struct {
	store  map[string]Value
	parent *Environment
}

func NewEnvironemnt() *Environment {
	s := make(map[string]Value)
	return &Environment{store: s}
}

func (e *Environment) NewChild() *Environment {
	env := NewEnvironemnt()
	env.parent = e
	return env
}

// returns false if name was not found
func (e *Environment) Get(name string) (Value, bool) {
	val, ok := e.store[name]

	if !ok && e.parent != nil {
		val, ok = e.parent.Get(name)
	}

	return val, ok
}

// returns false if name already exists in current env (not parent)
func (e *Environment) Let(name string, val Value) bool {
	if _, ok := e.store[name]; ok {
		return false
	} else {
		e.store[name] = val
		return true
	}
}

// returns false if name was not found
func (e *Environment) Assign(name string, val Value) bool {
	_, ok := e.store[name]

	if ok {
		e.store[name] = val
		return true
	} else if e.parent != nil {
		return e.parent.Assign(name, val)
	} else {
		return false
	}
}
