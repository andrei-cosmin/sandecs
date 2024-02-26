package component

import (
	"github.com/andrei-cosmin/sandata/data"
	"github.com/andrei-cosmin/sandata/flag"
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/bits-and-blooms/bitset"
)

type linkManager struct {
	poolCapacity      uint
	defaultLinkerSize uint
	linkedComponents  map[string]component.Id
	entityLinker      api.EntityContainer
	componentLinkers  data.Array[api.ComponentLinker]
	componentIdCursor component.Id
	flag.Flag
}

func NewLinkManager(numEntities, numComponents, poolCapacity uint, entityLinker api.EntityContainer) api.ComponentLinkManager {
	return &linkManager{
		poolCapacity:      poolCapacity,
		defaultLinkerSize: numEntities,
		linkedComponents:  make(map[string]component.Id),
		entityLinker:      entityLinker,
		componentLinkers:  *data.NewArray[api.ComponentLinker](numComponents),
		componentIdCursor: 0,
		Flag:              flag.New(),
	}
}

func (l *linkManager) Get(componentId component.Id) api.ComponentLinker {
	return l.componentLinkers.Get(componentId)
}

func (l *linkManager) UpdateLinks(scheduledEntityRemoves *bitset.BitSet) {

	for index := range l.componentIdCursor {
		resolver := l.componentLinkers.Get(index)
		resolver.Update(scheduledEntityRemoves)
		resolver.Refresh()
	}
	l.Clear()
}

func (l *linkManager) Accept(registration api.Registration) {
	registration.Execute(l)
}
