package sparse

import "github.com/bits-and-blooms/bitset"

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

func (a *Array[T]) GetPointer(index uint) *T {
	return &a.container[index]
}

func (a *Array[T]) Set(index uint, value T) {
	a.EnsureCapacity(index)
	a.container[index] = value
}

func (a *Array[T]) EnsureCapacity(index uint) {
	if index >= uint(len(a.container)) {
		a.resize(nextPowerOfTwo(index + 1))
	}
}

func (a *Array[T]) Size() uint {
	return uint(len(a.container))
}

func (a *Array[T]) Clear(set *bitset.BitSet) {
	for index, hasNext := set.NextSet(0); hasNext; index, hasNext = set.NextSet(index + 1) {
		a.container[index] = a.empty
	}
}

func (a *Array[T]) Empty() bool {
	return len(a.container) == 0
}

func (a *Array[T]) resize(capacity uint) {
	newContainer := make([]T, capacity)
	copy(newContainer, a.container)
	a.container = newContainer
}

func nextPowerOfTwo(value uint) uint {
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	value++

	return value
}
