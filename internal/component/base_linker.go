package component

import (
	"github.com/andrei-cosmin/sandata/bit"
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/bits-and-blooms/bitset"
)

// baseLinker struct - manages the linking of entities with a specific component type / tag (contains shared logic for all linkers)
//   - componentId component.Id - the id of the component type
//   - componentType string - the type of the component
//   - entityLinker entity.MaskView - an entity container (used to retrieve which entities exist in the world at a given time)
//   - callback func() - this callback will mark the link manager for update (will be executed when changes in  entities / instances are performed)
//   - scheduledRemoves *bit.BitMask - a bitset storing the entities that are scheduled for removal from the component type
//   - linkedEntities *bit.BitMask - a bitset storing the entities that are linked with the component type
type baseLinker struct {
	componentId      component.Id
	componentType    string
	entityLinker     entity.MaskView
	callback         func()
	scheduledRemoves *bit.BitMask
	linkedEntities   *bit.BitMask
}

// newBaseLinker method - creates a new base linker with the given parameters
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

// Link method - links the entity id with the component type
func (r *baseLinker) Link(entityId entity.Id) bool {
	// If the entity id is not part of the world, or it is already linked, return false (linking failed)
	if !r.entityLinker.EntityMask().Test(entityId) || r.Has(entityId) {
		return false
	}

	// Set the corresponding bit in the linked entities bitset
	r.linkedEntities.Bits.Set(entityId)

	// Flag the link manager for update
	r.callback()

	// Return true, as the linking was successful
	return true
}

// Has method - checks if the entity id is linked with the component type
func (r *baseLinker) Has(entityId entity.Id) bool {
	return r.linkedEntities.Test(entityId)
}

// Unlink method - unlinks the entity id from the component type
func (r *baseLinker) Unlink(entityId entity.Id) bool {
	// If the entity id is not part of the world, or the entity id is not linked with the component type, return
	if !r.entityLinker.EntityMask().Test(entityId) || !r.Has(entityId) {
		return false
	}

	// If the entity id is already scheduled for removal, return
	if r.scheduledRemoves.Test(entityId) {
		return false
	}

	// Set the corresponding bit in the scheduled removes bitset
	r.scheduledRemoves.Bits.Set(entityId)

	// Flag the link manager for update
	r.callback()
	return true
}

// ComponentId method - returns the component id
func (r *baseLinker) ComponentId() component.Id {
	return r.componentId
}

// EntityMask method - returns the linked entities as a bitset
func (r *baseLinker) EntityMask() bit.Mask {
	return r.linkedEntities
}

// CleanScheduledEntities  method - updates the linked entities (bitsets)
func (r *baseLinker) CleanScheduledEntities(scheduledSandboxRemoves bit.Mask) {
	// Perform logical OR (Union) between:
	// - the scheduled entity removes of the world
	// - the scheduled entity removes of the component
	scheduledSandboxRemoves.Union(r.scheduledRemoves.Bits)

	// Perform logical intersection between:
	// - the linked entities of the component
	// - the total scheduled entity removes (world + component)
	r.scheduledRemoves.Bits.InPlaceIntersection(r.linkedEntities.Bits)

	// Perform logical difference between:
	// - the linked entities of the component
	// - the intersection containing the total scheduled applicable entity removes (world + component)
	r.linkedEntities.Bits.InPlaceDifference(r.scheduledRemoves.Bits)
}

// Refresh method - clears the scheduled removals
func (r *baseLinker) Refresh() {
	r.scheduledRemoves.Bits.ClearAll()
}
