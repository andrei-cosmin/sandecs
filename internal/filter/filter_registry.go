package filter

import (
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/bits-and-blooms/bitset"
	"strconv"
	"strings"
)

type Registry struct {
	entityLinker         api.EntityContainer
	componentLinkManager api.ComponentLinkRetriever
	hashes               map[string]int
	caches               []*Cache
	linkedEntitiesBuffer *bitset.BitSet
	defaultCacheSize     uint
	stringBuilder        strings.Builder
}

func NewRegistry(size uint, entityLinker api.EntityContainer, componentLinkManager api.ComponentLinkManager) *Registry {
	return &Registry{
		entityLinker:         entityLinker,
		componentLinkManager: componentLinkManager,
		hashes:               make(map[string]int),
		caches:               make([]*Cache, 0),
		linkedEntitiesBuffer: bitset.New(size),
		defaultCacheSize:     size,
	}
}

func (r *Registry) Register(filterRules api.FilterRules) entity.View {
	hash := r.hashFilter(filterRules)

	if cacheIndex, ok := r.hashes[hash]; ok {
		return r.caches[cacheIndex]
	}

	filterCache := newCache(r.defaultCacheSize, filterRules)
	r.hashes[hash] = len(r.caches)
	r.caches = append(r.caches, filterCache)
	return filterCache
}

func (r *Registry) UpdateLinks() {
	for _, cache := range r.caches {
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

func (r *Registry) hashFilter(rules api.FilterRules) string {
	r.stringBuilder.Reset()

	for _, componentId := range rules.RequiredComponentIds() {
		r.stringBuilder.WriteString(strconv.Itoa(int(componentId)))
		r.stringBuilder.WriteString(",")
	}
	r.stringBuilder.WriteString("/")

	for _, componentId := range rules.ExcludedComponentIds() {
		r.stringBuilder.WriteString(strconv.Itoa(int(componentId)))
		r.stringBuilder.WriteString(",")
	}
	r.stringBuilder.WriteString("/")

	for _, componentId := range rules.UnionComponentIds() {
		r.stringBuilder.WriteString(strconv.Itoa(int(componentId)))
		r.stringBuilder.WriteString(",")
	}
	r.stringBuilder.WriteString("/")

	return r.stringBuilder.String()
}
