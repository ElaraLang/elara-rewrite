package interpreter

import (
	"fmt"
	"sync"
)

var contextPool = sync.Pool{
	New: func() interface{} {
		return &Context{
			variables:   map[uint64][]*Variable{},
			parameters:  []*Value{},
			namespace:   "",
			name:        "",
			contextPath: map[string][]*Context{},
			types:       map[string]Type{},
			parent:      nil,
			function:    nil,
			extensions:  map[Type]map[string]*Variable{},
		}
	},
}

func (c *Context) Clone() *Context {
	var parentClone *Context = nil
	if c.parent != nil {
		parentClone = c.parent.Clone()
	}
	fromPool := contextPool.Get().(*Context)
	fromPool.variables = c.variables
	fromPool.parameters = c.parameters
	fromPool.namespace = c.namespace
	fromPool.name = c.name
	fromPool.contextPath = c.contextPath
	fromPool.types = c.types
	fromPool.parent = parentClone
	fromPool.function = c.function
	fromPool.extensions = c.extensions
	return fromPool
}

func (c *Context) Cleanup() {
	c.function = nil

	for s := range c.variables {
		delete(c.variables, s)
	}
	c.variables = map[uint64][]*Variable{}
	c.parameters = []*Value{}

	c.namespace = ""
	c.name = ""
	c.contextPath = map[string][]*Context{}
	c.types = map[string]Type{}
	c.extensions = map[Type]map[string]*Variable{}
	c.parent = nil
	contextPool.Put(c)
}

func NewContext(init bool) *Context {
	c := contextPool.Get().(*Context)
	if !init {
		return c
	}
	c.DefineVariable(&Variable{
		Name:    "stdout",
		Mutable: false,
		Type:    OutputType,
		Value: &Value{
			Type:  OutputType,
			Value: nil,
		},
	})
	inputFunctionName := "input"
	inputFunction := Function{
		name: &inputFunctionName,
		Signature: Signature{
			Parameters: []Parameter{},
			ReturnType: StringType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				panic(err)
			}

			return NonReturningValue(&Value{Value: input, Type: StringType})
		}),
	}

	inputContract := NewFunctionType(&inputFunction)

	c.DefineVariable(&Variable{
		Name:    inputFunctionName,
		Mutable: false,
		Type:    inputContract,
		Value: &Value{
			Type:  inputContract,
			Value: inputFunction,
		},
	})

	emptyName := "empty"
	emptyFun := &Function{
		name: &emptyName,
		Signature: Signature{
			Parameters: []Parameter{},
			ReturnType: NewCollectionTypeOf(AnyType),
		},
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {
			return NonReturningValue(&Value{
				Type: NewCollectionTypeOf(AnyType),
				Value: &Collection{
					ElementType: AnyType,
					Elements:    []*Value{},
				},
			})
		}),
	}
	emptyContract := NewFunctionType(emptyFun)

	c.DefineVariable(&Variable{
		Name:    emptyName,
		Mutable: false,
		Type:    emptyContract,
		Value: &Value{
			Type:  emptyContract,
			Value: emptyFun,
		},
	})

	Init(c)
	return c
}
