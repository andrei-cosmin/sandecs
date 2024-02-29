package api

import (
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/bits-and-blooms/bitset"
)

// EntityContainer interface - used for retrieving entity ids (stored as bitsets)
type EntityContainer interface {
	EntityIds() *bitset.BitSet
}

// EntityLinker interface - used for linking entities
type EntityLinker interface {
	Link() entity.Id
	Unlink(entityId entity.Id)
	EntityIds() *bitset.BitSet
	GetScheduledRemoves() *bitset.BitSet
	Update()
	IsClear() bool
	Refresh()
}
