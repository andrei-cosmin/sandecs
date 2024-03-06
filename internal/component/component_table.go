package component

import (
	"github.com/andrei-cosmin/sandata/array"
	"github.com/andrei-cosmin/sandata/bit"
	"github.com/andrei-cosmin/sandata/pool"
	"github.com/andrei-cosmin/sandecs/component"
	"slices"
)

// table interface - used for a table of components (component instances storage)
type table[T any] interface {
	set(index uint)
	get(index uint) *T
	clear(mask bit.Mask, hook func(*T))
}

// basicTable struct - basic table of components (flat component instances storage)
//   - content data.Array[*T] - an array of component instances
type basicTable[T any] struct {
	content array.Array[*T]
}

// newBasicTable method - creates a new basic table with the given size
func newBasicTable[T component.Component](tableSize uint) *basicTable[T] {
	return &basicTable[T]{
		content: *array.New[*T](tableSize),
	}
}

// set method - sets a component instance for the given index
func (b *basicTable[T]) set(index uint) {
	b.content.Set(index, new(T))
}

// get method - retrieves a component instance for the given index
func (b *basicTable[T]) get(index uint) *T {
	return b.content.Get(index)
}

// clear method - clears the component instances for the given set of indices
func (b *basicTable[T]) clear(mask bit.Mask, hook func(*T)) {
	if hook != nil {
		b.content.ClearAllFunc(mask, hook)
	} else {
		b.content.ClearAll(mask)
	}
}

// pooledTable struct - pooled table of components (component instances storage with pooling)
//   - content data.Array[*T] - an array of component instances
//   - pool data.Pool[*T] - a pool of component instances
type pooledTable[T any] struct {
	content array.Array[*T]
	pool    pool.Pool[*T]
}

// newPooledTable method - creates a new pooled table with the given size and pool size
func newPooledTable[T component.Component](tableSize, poolSize uint) *pooledTable[T] {
	return &pooledTable[T]{
		content: *array.New[*T](tableSize),
		pool:    *pool.New[*T](poolSize),
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
func (p *pooledTable[T]) clear(mask bit.Mask, hook func(*T)) {
	// In the case of the pooled table push the instance to the pool, before clearing the table
	if hook != nil {
		p.content.ClearAllFunc(mask, func(instance *T) {
			hook(instance)
			p.pool.Push(instance)
		})
	} else {
		p.content.ClearAllFunc(mask, func(instance *T) {
			p.pool.Push(instance)
		})
	}
}

// compactTable struct - compact table of components
type compactTable[T any] struct {
	cursor  uint
	indices array.Array[uint]
	content []T
}

// newCompactTable method - creates a new compact table with the given size
func newCompactTable[T any](size uint) *compactTable[T] {
	return &compactTable[T]{
		indices: *array.New[uint](size),
		content: make([]T, size),
	}
}

// set method - sets a component instance for the given index
func (c *compactTable[T]) set(index uint) {
	// The item at the specified index will always be inserted at the cursor position
	// We associate the index with the current cursor position
	c.indices.Set(index, c.cursor)
	// If the position of the cursor is outside the bounds of the slice, grow, the slice
	if c.cursor >= uint(len(c.content)) {
		c.content = slices.Grow(c.content, len(c.content)+1)
		c.content = c.content[:cap(c.content)]
	}
	// Update the cursor position
	c.cursor++
}

// get method - retrieves a component instance for the given index
func (c *compactTable[T]) get(index uint) *T {
	// Retrieve the real position at which the item with the index value was saved
	sparseIndex := c.indices.Get(index)
	// If the item is outside the bounds, return nil
	if sparseIndex >= c.cursor {
		return nil
	}
	// Return the item at the specified position
	return &c.content[sparseIndex]
}

// clear method - clears the component instances for the given set of indices
func (c *compactTable[T]) clear(mask bit.Mask, hook func(*T)) {
	for index, hasNext := mask.NextSet(0); hasNext && index < c.indices.Size(); index, hasNext = mask.NextSet(index + 1) {
		// Update the cursor for each remove item
		c.cursor--
		// Retrieve the real position for item at the index value
		sparseIndex := c.indices.Get(index)
		// Switch the remove item with the last item in the table in order to keep compactness
		c.indices.Set(index, c.indices.Get(c.cursor))
		c.indices.Set(c.cursor, sparseIndex)

		// Retrieve the instance
		tempInstance := c.content[sparseIndex]
		// Switch the instance of the item with the last instance in the table in order to keep compactness
		// Also keep the instances correctly associated with the indices
		c.content[sparseIndex] = c.content[c.cursor]
		c.content[c.cursor] = tempInstance
	}
}
