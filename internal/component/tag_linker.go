package component

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/bits-and-blooms/bitset"
)

// linker struct - manages the linking of entities with a specific component type
//   - componentId component.Id - the id of the component type
//   - componentType string - the type of the component
//   - entityLinker api.EntityContainer - an entity container (used to retrieve which entities exist in the world at a given time)
//   - callback func() - this callback will mark the link manager for update (will be executed when changes in  entities / instances are performed)
//   - scheduledRemoves *bitset.Bitset - a bitset storing the entities that are scheduled for removal from the component type
//   - linkedEntities *bitset.Bitset - a bitset storing the entities that are linked with the component type
type linker struct {
	componentId      component.Id
	componentType    string
	entityLinker     api.EntityContainer
	callback         func()
	scheduledRemoves *bitset.BitSet
	linkedEntities   *bitset.BitSet
}

// newTagLinker method - creates a new linker with the given parameters
func newTagLinker(size uint, componentId component.Id, componentType string, entityLinker api.EntityContainer, callback func()) *linker {
	return &linker{
		componentId:      componentId,
		componentType:    componentType,
		entityLinker:     entityLinker,
		callback:         callback,
		scheduledRemoves: bitset.New(size),
		linkedEntities:   bitset.New(size),
	}
}

// Link method - links the entity id with the component type
func (r *linker) Link(entityId entity.Id) bool {
	// If the entity id is not part of the world, or it is already linked, return false (linking failed)
	if !r.entityLinker.EntityIds().Test(entityId) || r.linkedEntities.Test(entityId) {
		return false
	}

	// Set the corresponding bit in the linked entities bitset
	r.linkedEntities.Set(entityId)

	// Flag the link manager for update
	r.callback()

	// Return true, as the linking was successful
	return true
}

// Has method - checks if the entity id is linked with the component type
func (r *linker) Has(entityId entity.Id) bool {
	return r.linkedEntities.Test(entityId)
}

// Unlink method - unlinks the entity id from the component type
func (r *linker) Unlink(entityId entity.Id) {
	// If the entity id is not linked with the component type, return
	if !r.Has(entityId) {
		return
	}

	// If the entity id is already scheduled for removal, return
	if r.scheduledRemoves.Test(entityId) {
		return
	}

	// Set the corresponding bit in the scheduled removes bitset
	r.scheduledRemoves.Set(entityId)

	// Flag the link manager for update
	r.callback()
}

// ComponentId method - returns the component id
func (r *linker) ComponentId() component.Id {
	return r.componentId
}

// EntityIds method - returns the linked entities as a bitset
func (r *linker) EntityIds() *bitset.BitSet {
	return r.linkedEntities
}

// CleanScheduledEntities  method - updates the linked entities (bitsets)
func (r *linker) CleanScheduledEntities(scheduledSandboxRemoves *bitset.BitSet) {
	// Perform logical OR (Union) between:
	// - the scheduled entity removes of the world
	// - the scheduled entity removes of the component
	r.scheduledRemoves.InPlaceUnion(scheduledSandboxRemoves)

	// Perform logical difference between:
	// - the linked entities of the component
	// - the total scheduled entity removes (world + component)
	// NOTE: This will clear all the bits that are scheduled for removal
	r.linkedEntities.InPlaceDifference(r.scheduledRemoves)
}

// CleanScheduledInstances method - clears the instances corresponding to the scheduled entity removals
//
// NOTE: In the case of the tag linker, it will only call the listener OnRemove hook
func (r *linker) CleanScheduledInstances() {
	// No-op
}

// Refresh method - clears the scheduled removals
func (r *linker) Refresh() {
	r.scheduledRemoves.ClearAll()
}
