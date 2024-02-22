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
	linkedComponents  map[string]component.Id
	componentLinkers  *sparse.Array[api.ComponentLinker]
	componentIdCursor component.Id
	state.State
}

func NewLinkManager(size uint) api.ComponentLinkManager {
	return &linkManager{
		linkedComponents:  make(map[string]component.Id),
		componentLinkers:  sparse.New[api.ComponentLinker](size),
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
			newLinker[T](l.componentLinkers.Size(), l.componentIdCursor, componentType, l.Mark),
		)
		l.componentIdCursor++
	}

	return l.componentLinkers.Get(componentId)
}
