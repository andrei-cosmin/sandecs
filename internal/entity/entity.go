package entity

import (
	"github.com/andrei-cosmin/hakkt/entity"
	"github.com/andrei-cosmin/hakkt/internal/state"
	"github.com/bits-and-blooms/bitset"
)

type Linker struct {
	linkedEntities   *bitset.BitSet
	scheduledRemoves *bitset.BitSet
	state.State
}

func NewLinker(size uint) *Linker {
	return &Linker{
		linkedEntities:   bitset.New(size),
		scheduledRemoves: bitset.New(size),
		State:            state.New(),
	}
}

func (l *Linker) EntityIds() *bitset.BitSet {
	return l.linkedEntities
}

func (l *Linker) Link() entity.Id {
	entityId, exists := l.linkedEntities.NextClear(0)
	if !exists {
		entityId = l.linkedEntities.Len()
	}
	l.linkedEntities.Set(entityId)

	return entityId
}

func (l *Linker) Unlink(entityId entity.Id) {
	if !l.linkedEntities.Test(entityId) {
		return
	}

	l.scheduledRemoves.Set(entityId)
	l.Mark()
}

func (l *Linker) GetScheduledRemoves() *bitset.BitSet {
	return l.scheduledRemoves
}

func (l *Linker) Update() {
	l.linkedEntities.InPlaceDifference(l.scheduledRemoves)
}

func (l *Linker) Refresh() {
	l.scheduledRemoves.ClearAll()
	l.Reset()
}
