package api

import (
	"github.com/andrei-cosmin/sandata/bit"
	"github.com/andrei-cosmin/sandecs/component"
)

// ComponentLinkRetriever retrieves component linkers by ID.
type ComponentLinkRetriever interface {
	Get(componentId component.Id) ComponentLinker
}

// ComponentLinkManager manages all component linkers.
type ComponentLinkManager interface {
	Get(componentId component.Id) ComponentLinker
	UpdateLinks(scheduledSandboxRemoves bit.Mask)
	Accept(registration Registration)
	IsCleared() bool
}

// ComponentLinker manages instances of a single component type.
type ComponentLinker interface {
	ComponentId() component.Id
	EntityMask() bit.Mask
	CleanScheduledEntities(scheduledSandboxRemoves bit.Mask)
	CleanScheduledInstances()
	Refresh()
}

// Registration registers a component linker with the manager.
type Registration interface {
	Execute(ctx ComponentLinkManager)
}
