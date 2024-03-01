package component

import (
	"github.com/andrei-cosmin/sandata/data"
	"github.com/andrei-cosmin/sandecs/component"
)

// table interface - used for a table of components (component instances storage)
type table[T any] interface {
	set(index uint)
	get(index uint) *T
	clear(set data.Mask)
}

// basicTable struct - basic table of components (flat component instances storage)
//   - content data.Array[*T] - an array of component instances
type basicTable[T any] struct {
	content data.Array[*T]
}

// newBasicTable method - creates a new basic table with the given size
func newBasicTable[T component.Component](tableSize uint) *basicTable[T] {
	return &basicTable[T]{
		content: *data.NewArray[*T](tableSize),
	}
}

// set method - sets a component instance for the given index
func (t *basicTable[T]) set(index uint) {
	t.content.Set(index, new(T))
}

// get method - retrieves a component instance for the given index
func (t *basicTable[T]) get(index uint) *T {
	return t.content.Get(index)
}

// clear method - clears the component instances for the given set of indices
func (t *basicTable[T]) clear(mask data.Mask) {
	t.content.ClearAll(mask)
}

// pooledTable struct - pooled table of components (component instances storage with pooling)
//   - content data.Array[*T] - an array of component instances
//   - pool data.Pool[*T] - a pool of component instances
type pooledTable[T any] struct {
	content data.Array[*T]
	pool    data.Pool[*T]
}

// newPooledTable method - creates a new pooled table with the given size and pool size
func newPooledTable[T component.Component](tableSize, poolSize uint) *pooledTable[T] {
	return &pooledTable[T]{
		content: *data.NewArray[*T](tableSize),
		pool:    *data.NewPool[*T](poolSize),
	}
}

// set method - sets a component instance for the given index
func (p *pooledTable[T]) set(index uint) {
	// if a component instance is available in the pool, pop it and set it for the given index
	if value, ok := p.pool.Pop(); ok {
		p.content.Set(index, value)
	} else {
		// if no component instance is available in the pool, create a new one and set it for the given index
		p.content.Set(index, new(T))
	}
}

// get method - retrieves a component instance for the given index
func (p *pooledTable[T]) get(index uint) *T {
	return p.content.Get(index)
}

// clear method - clears the component instances for the given set of indices
func (p *pooledTable[T]) clear(mask data.Mask) {
	// In the case of the pooled table push the instance to the pool, before clearing the table
	p.content.ClearAllFunc(mask, func(instance *T) {
		p.pool.Push(instance)
	})
}
