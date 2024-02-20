package query

import (
	"github.com/andrei-cosmin/hakkt/component"
	"github.com/andrei-cosmin/hakkt/entity"
	"github.com/andrei-cosmin/hakkt/internal/state"
	"github.com/bits-and-blooms/bitset"
)

type CacheId = uint

type Cache struct {
	cacheId              CacheId
	matchedComponentIds  []component.Id
	excludedComponentIds []component.Id
	oneOfComponentIds    []component.Id
	linkMarker           *bitset.BitSet
	removeMarker         *bitset.BitSet
	linkedEntities       *bitset.BitSet
	entityIdsCache       []entity.Id
	state.State
}

func newCache(size uint, cacheId CacheId, matchedIds []component.Id, excludedIds []component.Id, oneOfIds []component.Id) *Cache {
	return &Cache{
		cacheId:              cacheId,
		matchedComponentIds:  matchedIds,
		excludedComponentIds: excludedIds,
		oneOfComponentIds:    oneOfIds,
		linkMarker:           bitset.New(size),
		removeMarker:         bitset.New(size),
		linkedEntities:       bitset.New(size),
		State:                state.New(),
	}
}

func (c *Cache) GetEntities() []entity.Id {
	if !c.IsUpdated() {
		c.Reset()
		c.entityIdsCache = c.entityIdsCache[:0]

		for entityId, hasNext := c.linkedEntities.NextSet(0); hasNext; entityId, hasNext = c.linkedEntities.NextSet(entityId + 1) {
			c.entityIdsCache = append(c.entityIdsCache, entityId)
		}
	}

	return c.entityIdsCache
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
