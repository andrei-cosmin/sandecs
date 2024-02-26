package api

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/bits-and-blooms/bitset"
)

type ComponentLinkRetriever interface {
	Get(componentId component.Id) ComponentLinker
}

type ComponentLinkManager interface {
	Get(componentId component.Id) ComponentLinker
	UpdateLinks(scheduledEntityRemoves *bitset.BitSet)
	Accept(registration Registration)
	IsClear() bool
}

type ComponentLinker interface {
	ComponentId() component.Id
	EntityIds() *bitset.BitSet
	Update(scheduledEntityRemoves *bitset.BitSet)
	Refresh()
}

type Registration interface {
	Execute(ctx ComponentLinkManager)
}
