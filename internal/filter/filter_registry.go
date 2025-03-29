package filter

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/bits-and-blooms/bitset"
	"strconv"
	"strings"
)

// Registry struct - holds the filter registry
//   - entityLinker entity.MaskView - the entity linker
//   - componentLinkManager api.ComponentLinkRetriever - the component link manager
//   - hashes map[string]int - a map of hashes for the filter rules (used to check if a filter is already registered)
//   - caches []*Cache - a slice of caches for the filters
//   - entitiesBuffer *bitset.Bitset - a bitset buffer (used for recomputing filtered entities)
//   - defaultCacheSize uint - the default size for the caches (bitset sizes)
type Registry struct {
	entityLinker         entity.MaskView
	componentLinkManager api.ComponentLinkRetriever
	hashes               map[string]int
	caches               []*Cache
	entitiesBuffer       *bitset.BitSet
	defaultCacheSize     uint
}

// NewRegistry method - creates a new registry with the given size, entity linker and component link manager
func NewRegistry(size uint, entityLinker entity.MaskView, componentLinkManager api.ComponentLinkManager) *Registry {
	return &Registry{
		entityLinker:         entityLinker,
		componentLinkManager: componentLinkManager,
		hashes:               make(map[string]int),
		caches:               make([]*Cache, 0),
		entitiesBuffer:       bitset.New(size),
		defaultCacheSize:     size,
	}
}

// Register method - registers a filter with the given filter rules and returns a view of the filter
func (r *Registry) Register(filterRules api.FilterRules) entity.View {
	// Get the hash for the filter rules (component ids are sorted before, so that 2 filters with the same component ids have the same hash)
	hash := hashFilter(filterRules)

	// If the filter is already registered, return the view of the existing cache
	if cacheIndex, ok := r.hashes[hash]; ok {
		return r.caches[cacheIndex]
	}

	// Create a new cache for the filter rules and add it to the registry
	filterCache := newCache(r.defaultCacheSize, filterRules)

	// Add the hash to the map, with the index of the new cache
	r.hashes[hash] = len(r.caches)
	r.caches = append(r.caches, filterCache)

	// Return the view of the new cache
	return filterCache
}

// UpdateLinks method - updates the linked entities for the filters
func (r *Registry) UpdateLinks() {
	// Iterate through the caches and update the linked entities
	for _, cache := range r.caches {
		// If no component ids are required or excluded, clear the linked entities buffer
		// If only unions are present, only logical ORs will be performed (in which case the masks present in the cache are sufficient)
		// Performing logical ORs with the empty buffer will not change the result, while having the sandbox entities will give incorrect results
		if len(cache.requiredComponentIds) == 0 && len(cache.excludedComponentIds) == 0 {
			r.entitiesBuffer.ClearAll()
		} else {
			// In case of required or excluded component ids, copy the linked entities from the entity linker into the buffer
			r.entityLinker.EntityMask().CopyFull(r.entitiesBuffer)
		}

		// Perform logical ANDs for all required component ids
		for _, requiredId := range cache.requiredComponentIds {
			var componentResolver = r.componentLinkManager.Get(requiredId)
			componentResolver.EntityMask().Intersection(r.entitiesBuffer)
		}

		// Perform logical XORs for all excluded component ids
		for _, excludedId := range cache.excludedComponentIds {
			var componentResolver = r.componentLinkManager.Get(excludedId)
			componentResolver.EntityMask().Difference(r.entitiesBuffer)
		}

		// Perform logical ORs for all union component ids
		for _, unionId := range cache.unionComponentIds {
			var componentResolver = r.componentLinkManager.Get(unionId)
			componentResolver.EntityMask().Union(r.entitiesBuffer)
		}

		// Check for new changes in the linked entities, and update the cache
		cache.checkForNewChanges(r.entitiesBuffer)
	}
}

// hashFilter method - hashes the filter rules and returns the hash as a string
func hashFilter(rules api.FilterRules) string {
	var stringBuilder strings.Builder

	// Hash the required, excluded and union component ids
	hashFilterComponentIds(&stringBuilder, rules.RequiredComponentIds())
	hashFilterComponentIds(&stringBuilder, rules.ExcludedComponentIds())
	hashFilterComponentIds(&stringBuilder, rules.UnionComponentIds())

	// Return the hash as a string
	return stringBuilder.String()
}

// hashFilterComponentIds method - hashes the component ids and appends them to the string builder
//
// # Filters should be registered at the beginning of the application's lifetime in the initialization phase
//
// NOTE: this method will create a unique hash for each filter, but not necessarily an optimized string
func hashFilterComponentIds(stringBuilder *strings.Builder, ids []component.Id) {
	for _, componentId := range ids {
		stringBuilder.WriteString(strconv.Itoa(int(componentId)))
		stringBuilder.WriteString(",")
	}
	stringBuilder.WriteString("/")
}
