package api

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/bits-and-blooms/bitset"
)

// ComponentLinkRetriever interface - stores a component linker
type ComponentLinkRetriever interface {
	Get(componentId component.Id) ComponentLinker
}

// ComponentLinkManager interface - manager for component linkers
type ComponentLinkManager interface {
	Get(componentId component.Id) ComponentLinker
	UpdateLinks(scheduledEntityRemoves *bitset.BitSet)
	Accept(registration Registration)
	IsClear() bool
}

// ComponentLinker interface - internal manager for component instances (of a single type)
type ComponentLinker interface {
	ComponentId() component.Id
	EntityIds() *bitset.BitSet
	Update(scheduledEntityRemoves *bitset.BitSet)
	Refresh()
}

// Registration interface - used for registering a component linker
type Registration interface {
	Execute(ctx ComponentLinkManager)
}
