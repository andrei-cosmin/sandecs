package filter

import (
	"github.com/andrei-cosmin/sandecs/component"
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
	hash := hashFilter(filterRules)

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

		for _, requiredId := range cache.requiredComponentIds {
			var componentResolver = r.componentLinkManager.Get(requiredId)
			r.linkedEntitiesBuffer.InPlaceIntersection(componentResolver.EntityIds())
		}

		for _, excludedId := range cache.excludedComponentIds {
			var componentResolver = r.componentLinkManager.Get(excludedId)
			r.linkedEntitiesBuffer.InPlaceDifference(componentResolver.EntityIds())
		}

		for _, unionId := range cache.unionComponentIds {
			var componentResolver = r.componentLinkManager.Get(unionId)
			r.linkedEntitiesBuffer.InPlaceUnion(componentResolver.EntityIds())
		}

		cache.updateWith(r.linkedEntitiesBuffer)
	}
}

func hashFilter(rules api.FilterRules) string {
	var stringBuilder strings.Builder

	hashFilterComponentIds(&stringBuilder, rules.RequiredComponentIds())
	hashFilterComponentIds(&stringBuilder, rules.ExcludedComponentIds())
	hashFilterComponentIds(&stringBuilder, rules.UnionComponentIds())

	return stringBuilder.String()
}

func hashFilterComponentIds(stringBuilder *strings.Builder, ids []component.Id) {
	for _, componentId := range ids {
		stringBuilder.WriteString(strconv.Itoa(int(componentId)))
		stringBuilder.WriteString(",")
	}
	stringBuilder.WriteString("/")
}
