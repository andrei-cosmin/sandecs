package filter

import (
	"github.com/andrei-cosmin/hakkt/entity"
	"github.com/andrei-cosmin/hakkt/internal/api"
	"github.com/bits-and-blooms/bitset"
)

type Registry struct {
	entityLinker         api.EntityContainer
	componentLinkManager api.ComponentLinkRetriever
	registry             []*Cache
	linkedEntitiesBuffer *bitset.BitSet
	defaultCacheSize     uint
}

func NewRegistry(size uint, entityLinker api.EntityContainer, componentLinkManager api.ComponentLinkManager) *Registry {
	return &Registry{
		entityLinker:         entityLinker,
		componentLinkManager: componentLinkManager,
		registry:             make([]*Cache, 0),
		linkedEntitiesBuffer: bitset.New(size),
		defaultCacheSize:     size,
	}
}

func (r *Registry) Register(filterRules api.FilterRules) entity.View {
	filterCache := newCache(r.defaultCacheSize, filterRules)
	r.registry = append(r.registry, filterCache)
	return filterCache
}

func (r *Registry) UpdateLinks() {
	for _, cache := range r.registry {
		if len(cache.requiredComponentIds) == 0 && len(cache.excludedComponentIds) == 0 {
			r.linkedEntitiesBuffer.ClearAll()
		} else {
			r.entityLinker.EntityIds().CopyFull(r.linkedEntitiesBuffer)
		}

		for index := range cache.requiredComponentIds {
			var componentResolver = r.componentLinkManager.Get(cache.requiredComponentIds[index])
			r.linkedEntitiesBuffer.InPlaceIntersection(componentResolver.EntityIds())
		}

		for index := range cache.excludedComponentIds {
			var componentResolver = r.componentLinkManager.Get(cache.excludedComponentIds[index])
			r.linkedEntitiesBuffer.InPlaceDifference(componentResolver.EntityIds())
		}

		for index := range cache.unionComponentIds {
			var componentResolver = r.componentLinkManager.Get(cache.unionComponentIds[index])
			r.linkedEntitiesBuffer.InPlaceUnion(componentResolver.EntityIds())
		}

		cache.updateWith(r.linkedEntitiesBuffer)
	}
}
