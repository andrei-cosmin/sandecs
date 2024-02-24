package component

import (
	"github.com/andrei-cosmin/sandata/data"
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/bits-and-blooms/bitset"
)

type table[T any] interface {
	set(index uint)
	get(index uint) *T
	clear(set *bitset.BitSet)
}

type basicTable[T any] struct {
	content data.Array[*T]
}

func newBasicTable[T component.Component](tableSize uint) *basicTable[T] {
	return &basicTable[T]{
		content: *data.NewArray[*T](tableSize),
	}
}

func (t *basicTable[T]) set(index uint) {
	t.content.Set(index, new(T))
}

func (t *basicTable[T]) get(index uint) *T {
	return t.content.Get(index)
}

func (t *basicTable[T]) clear(set *bitset.BitSet) {
	t.content.ClearAll(set)
}

type pooledTable[T any] struct {
	content data.Array[*T]
	pool    data.Pool[*T]
}

func newPooledTable[T component.Component](tableSize, poolSize uint) *pooledTable[T] {
	return &pooledTable[T]{
		content: *data.NewArray[*T](tableSize),
		pool:    *data.NewPool[*T](poolSize),
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

func (p *pooledTable[T]) clear(set *bitset.BitSet) {
	for index, hasNext := set.NextSet(0); hasNext; index, hasNext = set.NextSet(index + 1) {
		if index >= p.content.Size() {
			return
		}

		p.pool.Push(p.content.Get(index))
		p.content.Clear(index)
	}
}
