package component

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/bits-and-blooms/bitset"
)

// componentLinker struct - manages component instances of type T
//   - poolCapacity int - the capacity of the component pool
//   - components table[T] - a table of component instances
//   - linker - internal linker
type componentLinker[T component.Component] struct {
	poolCapacity uint
	components   table[T]
	linker
}

// newComponentLinker method - creates a new component linker with the given parameters
func newComponentLinker[T component.Component](
	size, poolCapacity uint,
	componentId component.Id, componentType string,
	entityLinker api.EntityContainer,
	callback func(),
) api.ComponentLinker {

	// Attach the correct type of component table
	var componentTable table[T]
	if poolCapacity > 0 {
		// Create a pooled component table
		componentTable = newPooledTable[T](size, poolCapacity)
	} else {
		// Create a basic component table (no pooling)
		componentTable = newBasicTable[T](size)
	}

	// Return the component linker
	return &componentLinker[T]{
		poolCapacity: poolCapacity,
		components:   componentTable,
		linker:       *newTagLinker(size, componentId, componentType, entityLinker, callback),
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
	if r.linker.Link(entityId) {
		// Set corresponding bit in the entity id bitset, marking that the entityId has the component type T
		r.components.set(entityId)

		// Return the component instance
		return r.components.get(entityId)
	}

	// If the entity id is already linked with the component type T, return nil
	return nil
}

// Update method - updates the component instances for the entities that are scheduled for removal
//
// NOTE: the method will overwrite the update method from the internal linker
func (r *componentLinker[T]) Update(scheduledEntityRemoves *bitset.BitSet) {
	// Remove the component instances for the entities that are scheduled for removal
	r.linker.Update(scheduledEntityRemoves)

	// Clear the scheduled removals
	r.components.clear(r.scheduledRemoves)
}
