package entity

import (
	"github.com/andrei-cosmin/hakkt/entity"
	"github.com/andrei-cosmin/hakkt/internal/state"
	"github.com/bits-and-blooms/bitset"
	"github.com/gammazero/deque"
)

type Linker struct {
	pool             *deque.Deque[entity.Id]
	linkedEntities   *bitset.BitSet
	scheduledRemoves *bitset.BitSet
	entityIdCursor   entity.Id
	state.State
}

func NewLinker(size uint) *Linker {
	return &Linker{
		pool:             deque.New[entity.Id](),
		linkedEntities:   bitset.New(size),
		scheduledRemoves: bitset.New(size),
		entityIdCursor:   0,
		State:            state.New(),
	}
}

func (l *Linker) EntityIds() *bitset.BitSet {
	return l.linkedEntities
}

func (l *Linker) Link() entity.Id {
	var entityId entity.Id
	if l.pool.Len() == 0 {
		entityId = l.entityIdCursor
		l.entityIdCursor++
	} else {
		entityId = l.pool.PopFront()
	}

	l.linkedEntities.Set(entityId)

	return entityId
}

func (l *Linker) Unlink(entityId entity.Id) {
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
	l.scheduledRemoves.ClearAll()
	l.Reset()
}
