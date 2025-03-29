package component

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/andrei-cosmin/sandecs/options"
	"github.com/bits-and-blooms/bitset"
)

// componentLinker struct - manages component instances of type T
//   - baseLinker - the base linker which holds the common linking functionality
//   - poolCapacity int - the capacity of the component pool
//   - components table[T] - a table of component instances
//   - additions *bitset.BitSet - a bitset which holds the additions of component instances between updates
//   - onLink func(*T) - a hook function which is called when a component instance is linked to an entity
//   - onUnlink func(*T) - a hook function which is called when a component instance is unlinked from an entity
type componentLinker[T component.Component] struct {
	baseLinker
	poolCapacity uint
	components   table[T]
	additions    *bitset.BitSet
	onLink       func(*T)
	onUnlink     func(*T)
}

// newComponentLinker method - creates a new component linker with the given parameters
func newComponentLinker[T component.Component](
	mode options.Mode,
	size, poolCapacity uint,
	componentId component.Id, componentType string,
	entityLinker entity.MaskView,
	callback func(),
) api.ComponentLinker {

	if poolCapacity <= 0 {
		poolCapacity = options.DefaultPoolCapacity
	}
	if size <= 0 {
		size = options.DefaultNumComponents
	}

	// Attach the correct type of component table
	var componentTable table[T]
	switch mode {
	default:
	case options.Standard:
		componentTable = newBasicTable[T](size)
	case options.Pooled:
		componentTable = newPooledTable[T](size, poolCapacity)
	case options.Compact:
		componentTable = newCompactTable[T](size)
	}

	// Return the component linker
	return &componentLinker[T]{
		poolCapacity: poolCapacity,
		components:   componentTable,
		additions:    bitset.New(size),
		baseLinker:   *newBaseLinker(size, componentId, componentType, entityLinker, callback),
	}
}

// Get method - retrieves a component instance for the given entity id
func (r *componentLinker[T]) Get(entityId entity.Id) *T {
	return r.components.get(entityId)
}

// Link method - links a component instance to the given entity id and returns the attached component instance
//
// NOTE: if the entity id is already linked with the component type T, the method will return nil
func (r *componentLinker[T]) Link(entityId entity.Id) *T {
	// If the entity id is not already linked with the component type T, link it and return the instance
	if r.baseLinker.Link(entityId) {
		// Set corresponding component instance
		r.components.set(entityId)

		// Set the entity id bit in the additions bitset
		r.additions.Set(entityId)

		// Return the component instance
		return r.components.get(entityId)
	}

	// If the entity id is already linked with the component type T, return nil
	return nil
}

// SetLinkHook method - sets a hook that will be called when an entity is linked to the tag
func (r *componentLinker[T]) SetLinkHook(onLink func(*T)) {
	r.onLink = onLink
}

// SetUnlinkHook method - sets a hook that will be called when an entity is unlinked from the tag
func (r *componentLinker[T]) SetUnlinkHook(onUnlink func(*T)) {
	r.onUnlink = onUnlink
}

// RemoveLinkHook method - removes the link hook
func (r *componentLinker[T]) RemoveLinkHook(onLink func()) {
	r.onLink = nil
}

// RemoveUnlinkHook method - removes the unlink hook
func (r *componentLinker[T]) RemoveUnlinkHook(onUnlink func()) {
	r.onUnlink = nil
}

// CleanScheduledInstances method - clears the instances corresponding to the scheduled entity removals
func (r *componentLinker[T]) CleanScheduledInstances() {
	// If the onLink hook is set, call it for each new linked entity since last update
	if r.onUnlink != nil {
		for addedEntityId, hasNext := r.additions.NextSet(0); hasNext; addedEntityId, hasNext = r.additions.NextSet(addedEntityId + 1) {
			r.onLink(r.components.get(addedEntityId))
		}
	}

	// Clear the instances of the scheduled entity removals
	r.components.clear(r.scheduledRemoves, r.onUnlink)

	// Clear the additions buffer
	r.additions.ClearAll()
}
