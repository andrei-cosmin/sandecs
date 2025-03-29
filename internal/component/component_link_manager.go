package component

import (
	"github.com/andrei-cosmin/sandata/array"
	"github.com/andrei-cosmin/sandata/bit"
	"github.com/andrei-cosmin/sandata/flag"
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/andrei-cosmin/sandecs/options"
)

// linkManager struct - manager for component linkers
//   - poolCapacity uint - the capacity of the pool (storing unused specific component instances)
//   - defaultLinkerSize uint - the default size of the component linker (size of pre-allocated slice of component references)
//   - linkedComponents map[string]component.Id - a map storing the relation between component type and component id
//   - entityLinker entity.MaskView - an entity container (used to retrieve which entities exist in the world at a given time)
//   - componentLinkers data.Array[ api.ComponentLinker ] - an array of component linkers
//   - componentIdCursor component.Id - a cursor for the component id (next available component id)
//   - flag.Flag - a flag to check if updates need to be performed
type linkManager struct {
	mode              options.Mode
	poolCapacity      uint
	defaultLinkerSize uint
	linkedComponents  map[string]component.Id
	entityLinker      entity.MaskView
	componentLinkers  array.Array[api.ComponentLinker]
	componentIdCursor component.Id
	flag.Flag
}

// NewLinkManager method - creates a new link manager with the given parameters
func NewLinkManager(mode options.Mode, numEntities, numComponents, poolCapacity uint, entityLinker entity.MaskView) api.ComponentLinkManager {
	return &linkManager{
		mode:              mode,
		poolCapacity:      poolCapacity,
		defaultLinkerSize: numEntities,
		linkedComponents:  make(map[string]component.Id),
		entityLinker:      entityLinker,
		componentLinkers:  *array.New[api.ComponentLinker](numComponents),
		componentIdCursor: 0,
		Flag:              flag.New(),
	}
}

// Get method - retrieves a component linker for the given component id
func (l *linkManager) Get(componentId component.Id) api.ComponentLinker {
	return l.componentLinkers.Get(componentId)
}

// UpdateLinks method - updates the links between components and entities (removes component instances for entities)
func (l *linkManager) UpdateLinks(scheduledSandboxRemoves bit.Mask) {

	for index := range l.componentIdCursor {
		// Retrieve the component linker for the given index
		resolver := l.componentLinkers.Get(index)
		// Clean corresponding scheduled entity removal bits from the linked entities bitset
		resolver.CleanScheduledEntities(scheduledSandboxRemoves)
		// Clean corresponding scheduled entity removal instances from the component table
		resolver.CleanScheduledInstances()
		// Refresh the component linker (clear scheduled removals buffer)
		resolver.Refresh()
	}
	// clear link manager flag (all linkers have been updated)
	l.Clear()
}

// Accept method - accepts a registration for a component linker
func (l *linkManager) Accept(registration api.Registration) {
	registration.Execute(l)
}
