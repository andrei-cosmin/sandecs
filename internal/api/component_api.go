package api

import (
	"github.com/andrei-cosmin/hakkt/component"
	"github.com/bits-and-blooms/bitset"
)

type ComponentLinkRetriever interface {
	Get(componentId component.Id) ComponentLinker
}

type ComponentLinkManager interface {
	Get(componentId component.Id) ComponentLinker
	UpdateLinks(scheduledEntityRemoves *bitset.BitSet)
	Accept(registration ComponentRegistration)
	IsUpdated() bool
}

type ComponentLinker interface {
	ComponentId() component.Id
	EntityIds() *bitset.BitSet
	Update(scheduledEntityRemoves *bitset.BitSet)
}

type ComponentRegistration interface {
	Execute(ctx ComponentLinkManager)
}
