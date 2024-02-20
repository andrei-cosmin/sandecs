package query

import (
	"github.com/andrei-cosmin/hakkt/internal/component"
	"github.com/andrei-cosmin/hakkt/internal/entity"
	"github.com/andrei-cosmin/hakkt/internal/sparse"
	"github.com/bits-and-blooms/bitset"
)

type Registry struct {
	entityLinker         *entity.Linker
	componentLinker      *component.LinkManager
	registry             map[string]uint
	caches               *sparse.Array[*Cache]
	linkedEntitiesBuffer *bitset.BitSet
	defaultCacheSize     uint
}

func NewRegistry(size uint, entityLinker *entity.Linker, componentLinker *component.LinkManager) *Registry {
	return &Registry{
		entityLinker:         entityLinker,
		componentLinker:      componentLinker,
		registry:             make(map[string]uint),
		caches:               sparse.New[*Cache](size),
		linkedEntitiesBuffer: bitset.New(size),
		defaultCacheSize:     size,
	}
}

func (r *Registry) Resolve(query *Info) *Cache {
	var queryId uint

	if id, ok := r.registry[query.hash]; ok {
		queryId = id
	} else {
		queryId = r.caches.Size()
		r.entityLinker.CopyLinkedEntitiesInto(r.linkedEntitiesBuffer)
		var matchedIds = make([]uint, len(query.match))
		var excludedIds = make([]uint, len(query.exclude))
		var oneOfIds = make([]uint, len(query.one))

		for index, componentType := range query.match {
			var componentResolver = r.componentLinker.Link(componentType)
			matchedIds[index] = componentResolver.GetComponentId()
		}
		for index, componentType := range query.exclude {
			var componentResolver = r.componentLinker.Link(componentType)
			excludedIds[index] = componentResolver.GetComponentId()
		}
		for index, componentType := range query.one {
			var componentResolver = r.componentLinker.Link(componentType)
			oneOfIds[index] = componentResolver.GetComponentId()
		}
		var queryCache = newCache(r.defaultCacheSize, queryId, matchedIds, excludedIds, oneOfIds)

		r.caches.Set(queryId, queryCache)
		r.registry[query.hash] = queryId
	}

	return r.caches.Get(queryId)
}

func (r *Registry) UpdateLinks() {
	for _, index := range r.registry {
		var cache = r.caches.Get(index)
		if cache != nil {
			r.entityLinker.CopyLinkedEntitiesInto(r.linkedEntitiesBuffer)

			for index := range cache.matchedComponentIds {
				var componentResolver = r.componentLinker.Find(cache.matchedComponentIds[index])
				r.linkedEntitiesBuffer.InPlaceIntersection(componentResolver.GetEntities())
			}

			for index := range cache.excludedComponentIds {
				var componentResolver = r.componentLinker.Find(cache.excludedComponentIds[index])
				r.linkedEntitiesBuffer.InPlaceDifference(componentResolver.GetEntities())
			}

			for index := range cache.oneOfComponentIds {
				var componentResolver = r.componentLinker.Find(cache.oneOfComponentIds[index])
				r.linkedEntitiesBuffer.InPlaceUnion(componentResolver.GetEntities())
			}

			cache.updateWith(r.linkedEntitiesBuffer)
		}
	}
}
