package filter

import (
	"github.com/andrei-cosmin/sandata/flag"
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/bits-and-blooms/bitset"
)

type CacheId = uint

type Cache struct {
	cacheId              CacheId
	requiredComponentIds []component.Id
	excludedComponentIds []component.Id
	unionComponentIds    []component.Id
	linkMarker           *bitset.BitSet
	removeMarker         *bitset.BitSet
	linkedEntities       *bitset.BitSet
	entityIdsCache       []entity.Id
	flag.Flag
}

func newCache(size uint, filterRules api.FilterRules) *Cache {
	return &Cache{
		requiredComponentIds: filterRules.RequiredComponentIds(),
		excludedComponentIds: filterRules.ExcludedComponentIds(),
		unionComponentIds:    filterRules.UnionComponentIds(),
		linkMarker:           bitset.New(size),
		removeMarker:         bitset.New(size),
		linkedEntities:       bitset.New(size),
		Flag:                 flag.New(),
	}
}

func (c *Cache) EntityIds() []entity.Id {
	if !c.IsClear() {
		c.Clear()
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

	c.Set()
	c.linkedEntities.InPlaceUnion(entities)
}

func (c *Cache) removeAll(entities *bitset.BitSet) {
	if c.linkedEntities.IntersectionCardinality(entities) == 0 {
		return
	}

	c.Set()
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
