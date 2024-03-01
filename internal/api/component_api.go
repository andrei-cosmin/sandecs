package api

import (
	"github.com/andrei-cosmin/sandata/data"
	"github.com/andrei-cosmin/sandecs/component"
)

// ComponentLinkRetriever interface - stores a component linker
type ComponentLinkRetriever interface {
	Get(componentId component.Id) ComponentLinker
}

// ComponentLinkManager interface - manager for component linkers
type ComponentLinkManager interface {
	Get(componentId component.Id) ComponentLinker
	UpdateLinks(scheduledSandboxRemoves data.Mask)
	Accept(registration Registration)
	IsClear() bool
}

// ComponentLinker interface - internal manager for component instances (of a single type)
type ComponentLinker interface {
	ComponentId() component.Id
	EntityMask() data.Mask
	CleanScheduledEntities(scheduledSandboxRemoves data.Mask)
	CleanScheduledInstances()
	Refresh()
}

// Registration interface - used for registering a component linker
type Registration interface {
	Execute(ctx ComponentLinkManager)
}
