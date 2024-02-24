package api

import (
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/bits-and-blooms/bitset"
)

type EntityContainer interface {
	EntityIds() *bitset.BitSet
}

type EntityLinker interface {
	Link() entity.Id
	Unlink(entityId entity.Id)
	EntityIds() *bitset.BitSet
	GetScheduledRemoves() *bitset.BitSet
	Update()
	IsClear() bool
	Refresh()
}
