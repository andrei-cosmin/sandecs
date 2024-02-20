package query

import (
	"github.com/andrei-cosmin/hakkt/internal/state"
	"github.com/bits-and-blooms/bitset"
)

type Cache struct {
	id                   uint
	matchedComponentIds  []uint
	excludedComponentIds []uint
	oneOfComponentIds    []uint
	linkMarker           *bitset.BitSet
	removeMarker         *bitset.BitSet
	linkedEntities       *bitset.BitSet
	entityCache          []uint
	state.State
}

func newCache(size uint, id uint, matchedIds []uint, excludedIds []uint, oneOfIds []uint) *Cache {
	return &Cache{
		id:                   id,
		matchedComponentIds:  matchedIds,
		excludedComponentIds: excludedIds,
		oneOfComponentIds:    oneOfIds,
		linkMarker:           bitset.New(size),
		removeMarker:         bitset.New(size),
		linkedEntities:       bitset.New(size),
		State:                state.New(),
	}
}

func (c *Cache) GetEntities() []uint {
	if !c.IsUpdated() {
		c.Reset()
		c.entityCache = c.entityCache[:0]

		index := uint(0)
		for index, hasNext := c.linkedEntities.NextSet(index); hasNext; index, hasNext = c.linkedEntities.NextSet(index + 1) {
			c.entityCache = append(c.entityCache, index)
		}
	}

	return c.entityCache
}

func (c *Cache) linkAll(entities *bitset.BitSet) {
	if c.linkedEntities.IsSuperSet(entities) {
		return
	}

	c.Mark()
	c.linkedEntities.InPlaceUnion(entities)
}

func (c *Cache) removeAll(entities *bitset.BitSet) {
	if c.linkedEntities.IntersectionCardinality(entities) == 0 {
		return
	}

	c.Mark()
	c.linkedEntities.InPlaceDifference(entities)
}

func (c *Cache) updateWith(after *bitset.BitSet) {
	after.CopyFull(c.linkMarker)
	c.linkMarker.InPlaceDifference(c.linkedEntities)
	c.linkAll(c.linkMarker)

	c.linkedEntities.CopyFull(c.removeMarker)
	c.removeMarker.InPlaceDifference(after)
	c.removeAll(c.removeMarker)
}
