package component

import (
	"slices"

	"github.com/andrei-cosmin/sandata/array"
	"github.com/andrei-cosmin/sandata/bit"
	"github.com/andrei-cosmin/sandata/pool"
	"github.com/andrei-cosmin/sandecs/component"
)

// table is the storage interface for component instances.
type table[T any] interface {
	set(index uint)
	get(index uint) *T
	clear(mask bit.Mask, hook func(*T))
}

// basicTable stores components in a flat array (Standard mode).
type basicTable[T any] struct {
	content array.Array[*T]
}

func newBasicTable[T component.Component](tableSize uint) *basicTable[T] {
	return &basicTable[T]{content: *array.New[*T](tableSize)}
}

func (b *basicTable[T]) set(index uint) {
	b.content.Set(index, new(T))
}

func (b *basicTable[T]) get(index uint) *T {
	return b.content.Get(index)
}

func (b *basicTable[T]) clear(mask bit.Mask, hook func(*T)) {
	if hook != nil {
		b.content.ClearAllFunc(mask, hook)
	} else {
		b.content.ClearAll(mask)
	}
}

// pooledTable stores components with instance reuse (Pooled mode).
type pooledTable[T any] struct {
	content array.Array[*T]
	pool    pool.Pool[*T]
}

func newPooledTable[T component.Component](tableSize, poolSize uint) *pooledTable[T] {
	return &pooledTable[T]{
		content: *array.New[*T](tableSize),
		pool:    *pool.New[*T](poolSize),
	}
}

func (p *pooledTable[T]) set(index uint) {
	if value, ok := p.pool.Pop(); ok {
		p.content.Set(index, value)
	} else {
		p.content.Set(index, new(T))
	}
}

func (p *pooledTable[T]) get(index uint) *T {
	return p.content.Get(index)
}

func (p *pooledTable[T]) clear(mask bit.Mask, hook func(*T)) {
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

// compactTable stores components densely using sparse set (Compact mode).
type compactTable[T any] struct {
	cursor  uint
	indices array.Array[uint] // entity ID → slot
	reverse array.Array[uint] // slot → entity ID
	content []T
}

func newCompactTable[T any](size uint) *compactTable[T] {
	return &compactTable[T]{
		indices: *array.New[uint](size),
		reverse: *array.New[uint](size),
		content: make([]T, size),
	}
}

func (c *compactTable[T]) set(index uint) {
	c.indices.Set(index, c.cursor)
	c.reverse.Set(c.cursor, index)
	if c.cursor >= uint(len(c.content)) {
		c.content = slices.Grow(c.content, len(c.content)+1)
		c.content = c.content[:cap(c.content)]
	}
	c.cursor++
}

func (c *compactTable[T]) get(index uint) *T {
	sparseIndex := c.indices.Get(index)
	if sparseIndex >= c.cursor {
		return nil
	}
	return &c.content[sparseIndex]
}

func (c *compactTable[T]) clear(mask bit.Mask, hook func(*T)) {
	for index, hasNext := mask.NextSet(0); hasNext && index < c.indices.Size(); index, hasNext = mask.NextSet(index + 1) {
		slotToRemove := c.indices.Get(index)
		if slotToRemove >= c.cursor {
			continue
		}
		c.cursor--
		if slotToRemove != c.cursor {
			lastEntity := c.reverse.Get(c.cursor)
			c.content[slotToRemove] = c.content[c.cursor]
			c.indices.Set(lastEntity, slotToRemove)
			c.reverse.Set(slotToRemove, lastEntity)
		}
		c.indices.Set(index, c.cursor)
	}
}
