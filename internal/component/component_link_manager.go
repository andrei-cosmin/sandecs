package component

import (
	"github.com/andrei-cosmin/sandata/data"
	"github.com/andrei-cosmin/sandata/flag"
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/bits-and-blooms/bitset"
)

// linkManager struct - manager for component linkers
//   - poolCapacity uint - the capacity of the pool (storing unused specific component instances)
//   - defaultLinkerSize uint - the default size of the component linker (size of pre-allocated slice of component references)
//   - linkedComponents map[string]component.Id - a map storing the relation between component type and component id
//   - entityLinker api.EntityContainer - an entity container (used to retrieve which entities exist in the world at a given time)
//   - componentLinkers data.Array[ api.ComponentLinker ] - an array of component linkers
//   - componentIdCursor component.Id - a cursor for the component id (next available component id)
//   - flag.Flag - a flag to check if updates need to be performed
type linkManager struct {
	poolCapacity      uint
	defaultLinkerSize uint
	linkedComponents  map[string]component.Id
	entityLinker      api.EntityContainer
	componentLinkers  data.Array[api.ComponentLinker]
	componentIdCursor component.Id
	flag.Flag
}

// NewLinkManager method - creates a new link manager with the given parameters
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

// Get method - retrieves a component linker for the given component id
func (l *linkManager) Get(componentId component.Id) api.ComponentLinker {
	return l.componentLinkers.Get(componentId)
}

// UpdateLinks method - updates the links between components and entities (removes component instances for entities)
func (l *linkManager) UpdateLinks(scheduledEntityRemoves *bitset.BitSet) {

	for index := range l.componentIdCursor {
		// retrieve the component linker for the given index
		resolver := l.componentLinkers.Get(index)
		// remove the component instances for the entities that are scheduled for removal
		resolver.Update(scheduledEntityRemoves)
		// refresh the component linker (clear scheduled removals)
		resolver.Refresh()
	}
	// clear link manager flag (all linker have been updated)
	l.Clear()
}

// Accept method - accepts a registration for a component linker
func (l *linkManager) Accept(registration api.Registration) {
	registration.Execute(l)
}
