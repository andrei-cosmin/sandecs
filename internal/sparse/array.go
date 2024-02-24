package sparse

import (
	"github.com/andrei-cosmin/hakkt/internal/util"
	"github.com/bits-and-blooms/bitset"
	"slices"
)

type Array[T any] struct {
	container []T
	empty     T
}

func New[T any](size uint) *Array[T] {
	return &Array[T]{
		container: make([]T, size),
	}
}

func (a *Array[T]) Get(index uint) T {
	return a.container[index]
}

func (a *Array[T]) Set(index uint, value T) {
	a.ensureCapacity(index)
	a.container[index] = value
}

func (a *Array[T]) Size() uint {
	return uint(len(a.container))
}

func (a *Array[T]) ClearAll(set *bitset.BitSet) {
	for index, hasNext := set.NextSet(0); hasNext; index, hasNext = set.NextSet(index + 1) {
		if index >= uint(len(a.container)) {
			return
		}
		a.container[index] = a.empty
	}
}

func (a *Array[T]) Clear(index uint) {
	a.container[index] = a.empty
}

func (a *Array[T]) ensureCapacity(index uint) {
	if index >= uint(len(a.container)) {
		a.container = slices.Grow(a.container, int(util.NextPowerOfTwo(index+1))-len(a.container))
		a.container = a.container[:cap(a.container)]
	}
}
