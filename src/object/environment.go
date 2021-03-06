package object

import (
	"errors"
	"fmt"
)

var (
	constError        = errors.New("constant can't not be changed")
	errAlreadyDefined = errors.New("symbol already defined")
	errNotDefined     = errors.New("symbol not defined")
)

func IsConstErr(e error) bool {
	return e == constError
}

type eco struct {
	v        Object
	readonly bool
}

type Environment struct {
	store  map[string]*eco
	parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]*eco),
	}
}

func NewEnclosedEnv(outer *Environment) *Environment {
	env := NewEnvironment()
	env.parent = outer
	return env
}

func (e *Environment) SetParent(env *Environment) {
	e.parent = env
}

func (e *Environment) Parent() *Environment {
	return e.parent
}

func (e *Environment) Print(indent string) {
	for k, v := range e.store {
		fmt.Printf("%s%s = %s\n  %sConst: %t\n", indent, k, v.v.Inspect(), indent, v.readonly)
	}

	if e.parent != nil {
		if e.parent.parent == nil {
			fmt.Printf("\n%sGlobal:\n", indent)
		} else {
			fmt.Printf("\n%sParent:\n", indent)
		}
		e.parent.Print(indent + "  ")
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if ok {
		return obj.v, ok
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, false
}

func (e *Environment) GetLocal(name string) (Object, bool) {
	obj, ok := e.store[name]
	if ok {
		return obj.v, ok
	}
	return nil, false
}

func (e *Environment) IsConst(name string) bool {
	obj, ok := e.store[name]
	if ok {
		return obj.readonly
	}

	if e.parent != nil {
		return e.parent.IsConst(name)
	}
	return false
}

func (e *Environment) IsConstLocal(name string) bool {
	obj, ok := e.store[name]
	if ok {
		return obj.readonly
	}
	return false
}

func (e *Environment) isLocalConst(name string) bool {
	obj, ok := e.store[name]
	if ok {
		return obj.readonly
	}
	return false
}

func (e *Environment) Create(name string, val Object) (Object, error) {
	if _, exists := e.store[name]; exists {
		return nil, errAlreadyDefined
	}

	return e.setLocal(name, val), nil
}

func (e *Environment) CreateConst(name string, val Object) (Object, error) {
	if _, exists := e.store[name]; exists {
		return nil, errAlreadyDefined
	}

	e.store[name] = &eco{
		v:        val,
		readonly: true,
	}
	return val, nil
}

func (e *Environment) Set(name string, val Object) (Object, error) {
	if v, exists := e.store[name]; exists {
		if v.readonly {
			return nil, constError
		}
		return e.setLocal(name, val), nil
	}

	if e.parent != nil {
		return e.parent.Set(name, val)
	}
	return nil, errNotDefined
}

func (e *Environment) SetForce(name string, val Object, readonly bool) {
	e.store[name] = &eco{
		v:        val,
		readonly: readonly,
	}
}

func (e *Environment) setLocal(name string, val Object) Object {
	e.store[name] = &eco{v: val}
	return val
}

func (e *Environment) UnsetLocal(name string) {
	delete(e.store, name)
}

func (e *Environment) Unset(name string) {
	if _, exists := e.store[name]; exists {
		delete(e.store, name)
	}
	if e.parent != nil {
		e.parent.Unset(name)
	}
}
