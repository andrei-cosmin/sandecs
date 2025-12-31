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

// linkManager manages all component linkers.
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

// NewLinkManager creates a link manager with pre-allocated capacity.
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

// Get returns the component linker for the given ID.
func (l *linkManager) Get(componentId component.Id) api.ComponentLinker {
	return l.componentLinkers.Get(componentId)
}

// UpdateLinks processes all pending component removals.
func (l *linkManager) UpdateLinks(scheduledSandboxRemoves bit.Mask) {
	for index := range l.componentIdCursor {
		resolver := l.componentLinkers.Get(index)
		resolver.CleanScheduledEntities(scheduledSandboxRemoves)
		resolver.CleanScheduledInstances()
		resolver.Refresh()
	}
	l.Clear()
}

// Accept processes a component registration.
func (l *linkManager) Accept(registration api.Registration) {
	registration.Execute(l)
}
