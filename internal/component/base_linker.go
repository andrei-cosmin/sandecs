package component

import (
	"github.com/andrei-cosmin/sandata/bit"
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/bits-and-blooms/bitset"
)

// baseLinker provides shared linking logic for components and tags.
type baseLinker struct {
	componentId      component.Id
	componentType    string
	entityLinker     entity.MaskView
	callback         func()
	scheduledRemoves *bit.BitMask
	linkedEntities   *bit.BitMask
}

func newBaseLinker(size uint, componentId component.Id, componentType string, entityLinker entity.MaskView, callback func()) *baseLinker {
	return &baseLinker{
		componentId:      componentId,
		componentType:    componentType,
		entityLinker:     entityLinker,
		callback:         callback,
		scheduledRemoves: bit.NewMask(bitset.New(size)),
		linkedEntities:   bit.NewMask(bitset.New(size)),
	}
}

// Link associates the entity with this component. Returns false if already linked or entity doesn't exist.
func (r *baseLinker) Link(entityId entity.Id) bool {
	if !r.entityLinker.EntityMask().Test(entityId) || r.Has(entityId) {
		return false
	}
	r.linkedEntities.Bits().Set(entityId)
	r.callback()
	return true
}

// Has returns true if the entity has this component.
func (r *baseLinker) Has(entityId entity.Id) bool {
	return r.linkedEntities.Test(entityId)
}

// Unlink schedules removal of the component from the entity.
func (r *baseLinker) Unlink(entityId entity.Id) bool {
	if !r.entityLinker.EntityMask().Test(entityId) || !r.Has(entityId) {
		return false
	}
	if r.scheduledRemoves.Test(entityId) {
		return false
	}
	r.scheduledRemoves.Bits().Set(entityId)
	r.callback()
	return true
}

// ComponentId returns the component ID.
func (r *baseLinker) ComponentId() component.Id {
	return r.componentId
}

// EntityMask returns the bitmask of linked entities.
func (r *baseLinker) EntityMask() bit.Mask {
	return r.linkedEntities
}

// CleanScheduledEntities removes scheduled entities from the linked set.
func (r *baseLinker) CleanScheduledEntities(scheduledSandboxRemoves bit.Mask) {
	scheduledSandboxRemoves.Union(r.scheduledRemoves.Bits())
	r.scheduledRemoves.Bits().InPlaceIntersection(r.linkedEntities.Bits())
	r.linkedEntities.Bits().InPlaceDifference(r.scheduledRemoves.Bits())
}

// Refresh clears scheduled removals.
func (r *baseLinker) Refresh() {
	r.scheduledRemoves.Bits().ClearAll()
}
