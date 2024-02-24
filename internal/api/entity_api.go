package api

import (
	"github.com/andrei-cosmin/hakkt/entity"
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
	IsUpdated() bool
	Refresh()
}
