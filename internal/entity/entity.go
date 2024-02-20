package entity

import (
	"github.com/andrei-cosmin/hakkt-ecs/internal/state"
	"github.com/bits-and-blooms/bitset"
	"github.com/gammazero/deque"
)

type Linker struct {
	pool             *deque.Deque[uint]
	linkedEntities   *bitset.BitSet
	scheduledRemoves *bitset.BitSet
	entityCursor     uint
	state.State
}

func NewLinker(size uint) *Linker {
	return &Linker{
		pool:             deque.New[uint](),
		linkedEntities:   bitset.New(size),
		scheduledRemoves: bitset.New(size),
		entityCursor:     0,
		State:            state.New(),
	}
}

func (l *Linker) CopyLinkedEntitiesInto(set *bitset.BitSet) {
	l.linkedEntities.CopyFull(set)
}

func (l *Linker) Link() uint {
	var entityId uint

	if l.pool.Len() == 0 {
		entityId = l.entityCursor
		l.entityCursor++
	} else {
		entityId = l.pool.PopFront()
	}

	l.linkedEntities.Set(entityId)

	return entityId
}

func (l *Linker) Unlink(entityId uint) {
	if !l.linkedEntities.Test(entityId) {
		return
	}

	l.linkedEntities.Clear(entityId)
	l.scheduledRemoves.Set(entityId)
	l.Mark()
}

func (l *Linker) GetScheduledRemoves() *bitset.BitSet {
	return l.scheduledRemoves
}

func (l *Linker) Refresh() {
	l.Reset()
	l.scheduledRemoves.ClearAll()
}
