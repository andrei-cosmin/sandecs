package filter

import (
	"github.com/andrei-cosmin/sandata/bit"
	"github.com/andrei-cosmin/sandata/flag"
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/bits-and-blooms/bitset"
)

// CacheId type - cache id
type CacheId = uint

// Cache struct - filter cache stores the context for a filter (component types , rules, linked entities)
//   - cacheId CacheId - the id of the filter
//   - requiredComponentIds []component.Id - the required component ids
//   - excludedComponentIds []component.Id - the excluded component ids
//   - unionComponentIds []component.Id - the union component ids
//   - linkMaskBuffer *bitset.Bitset - a bitset buffer
//   - unlinkMaskBuffer *bitset.Bitset - a bitset buffer
//   - filteredEntities *data.BitMask - a bitset storing the entities corresponding to the filter
//   - entityIdsCache []entity.Id -  a cache for the expanded entity ids (pre-allocated buffer for storing the entity ids)
//   - Flag: a flag used to mark that the cache is dirty and the expanded entity ids need to be refreshed
type Cache struct {
	cacheId              CacheId
	requiredComponentIds []component.Id
	excludedComponentIds []component.Id
	unionComponentIds    []component.Id
	linkMaskBuffer       *bitset.BitSet
	unlinkMaskBuffer     *bitset.BitSet
	filteredEntities     *bit.BitMask
	entityIdsCache       []entity.Id
	flag.Flag
}

// newCache method - creates a new cache with the given size for bitsets and filter rules
func newCache(size uint, filterRules api.FilterRules) *Cache {
	return &Cache{
		requiredComponentIds: filterRules.RequiredComponentIds(),
		excludedComponentIds: filterRules.ExcludedComponentIds(),
		unionComponentIds:    filterRules.UnionComponentIds(),
		linkMaskBuffer:       bitset.New(size),
		unlinkMaskBuffer:     bitset.New(size),
		filteredEntities:     bit.NewMask(bitset.New(size)),
		Flag:                 flag.New(),
	}
}

// EntityIds method - retrieves the filtered entities (as a slice of entity ids converted from the bitset)
func (c *Cache) EntityIds() []entity.Id {
	// If the cache is dirty, refresh the entity ids
	if !c.IsClear() {
		// Clear the cache flag
		c.Clear()
		// Clear the entity ids buffer
		c.entityIdsCache = c.entityIdsCache[:0]

		// Iterate through the linked entities bitset and add the entity ids to the buffer
		for entityId, hasNext := c.filteredEntities.NextSet(0); hasNext; entityId, hasNext = c.filteredEntities.NextSet(entityId + 1) {
			c.entityIdsCache = append(c.entityIdsCache, entityId)
		}
	}

	// Return the entity ids buffer
	return c.entityIdsCache
}

// EntityMask method - returns the filtered entities as a bitset
func (c *Cache) EntityMask() bit.Mask {
	return c.filteredEntities
}

// checkForNewAdditions method - checks the entities from the given bitset and adds them to the filtered entities if they are not already included
func (c *Cache) checkForNewAdditions(entities *bitset.BitSet) {
	// If the entities are already included in the filtered entities, return
	if c.filteredEntities.Bits.IsSuperSet(entities) {
		return
	}

	// Mark the cache as dirty
	c.Set()

	// Add the entities to the filtered entities bitset
	c.filteredEntities.Bits.InPlaceUnion(entities)
}

// checkForNewRemovals method -  checks the entities from the given bitset and removes them from the filtered entities if they are included
func (c *Cache) checkForNewRemovals(entities *bitset.BitSet) {
	// If none of the entities are included in the filtered entities, return
	if c.filteredEntities.IntersectionCardinality(entities) == 0 {
		return
	}

	// Mark the cache as dirty
	c.Set()

	// Remove the entities from the filtered entities bitset
	c.filteredEntities.Bits.InPlaceDifference(entities)
}

// checkForNewChanges method - updates the cache with the recomputed filtered entities
func (c *Cache) checkForNewChanges(recomputedFilteredEntities *bitset.BitSet) {
	// Check if the recomputed entities contain new additions
	recomputedFilteredEntities.CopyFull(c.linkMaskBuffer)
	c.linkMaskBuffer.InPlaceDifference(c.filteredEntities.Bits)
	c.checkForNewAdditions(c.linkMaskBuffer)

	// Check if the recomputed entities contain new removals
	c.filteredEntities.CopyFull(c.unlinkMaskBuffer)
	c.unlinkMaskBuffer.InPlaceDifference(recomputedFilteredEntities)
	c.checkForNewRemovals(c.unlinkMaskBuffer)
}
