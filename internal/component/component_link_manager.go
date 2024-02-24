package component

import (
	"github.com/andrei-cosmin/hakkt/component"
	"github.com/andrei-cosmin/hakkt/internal/api"
	"github.com/andrei-cosmin/hakkt/internal/sparse"
	"github.com/andrei-cosmin/hakkt/internal/state"
	"github.com/bits-and-blooms/bitset"
	"reflect"
)

type linkManager struct {
	poolCapacity      uint
	defaultLinkerSize uint
	linkedComponents  map[string]component.Id
	entityLinker      api.EntityContainer
	componentLinkers  *sparse.Array[api.ComponentLinker]
	componentIdCursor component.Id
	state.State
}

func NewLinkManager(numEntities, numComponents, poolCapacity uint, entityLinker api.EntityContainer) api.ComponentLinkManager {
	return &linkManager{
		poolCapacity:      poolCapacity,
		defaultLinkerSize: numEntities,
		linkedComponents:  make(map[string]component.Id),
		entityLinker:      entityLinker,
		componentLinkers:  sparse.New[api.ComponentLinker](numComponents),
		componentIdCursor: 0,
		State:             state.New(),
	}
}

func (l *linkManager) Get(componentId component.Id) api.ComponentLinker {
	return l.componentLinkers.Get(componentId)
}

func (l *linkManager) UpdateLinks(scheduledEntityRemoves *bitset.BitSet) {
	for index := uint(0); index < l.componentIdCursor; index++ {
		resolver := l.componentLinkers.Get(index)
		resolver.Update(scheduledEntityRemoves)
	}
	l.Reset()
}

func (l *linkManager) Accept(registration api.ComponentRegistration) {
	registration.Execute(l)
}

func RegisterLinker[T component.Component](componentLinkManager api.ComponentLinkManager) api.ComponentLinker {
	var componentId component.Id
	var componentType = reflect.TypeFor[T]().String()
	l := componentLinkManager.(*linkManager)

	if id, ok := l.linkedComponents[componentType]; ok {
		componentId = id
	} else {
		componentId = l.componentIdCursor
		l.linkedComponents[componentType] = l.componentIdCursor
		l.componentLinkers.Set(
			l.componentIdCursor,
			newLinker[T](l.defaultLinkerSize, l.poolCapacity, l.componentIdCursor, componentType, l.entityLinker, l.Mark),
		)
		l.componentIdCursor++
	}

	return l.componentLinkers.Get(componentId)
}
