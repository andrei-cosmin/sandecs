package api

import (
	"github.com/andrei-cosmin/sandata/bit"
	"github.com/andrei-cosmin/sandecs/entity"
)

// EntityLinker interface - used for linking entities
type EntityLinker interface {
	// Link method - links a new entity with the sandbox, returning the entity id
	Link() entity.Id

	// Unlink method - unlinks the entity id from the sandbox entirely (this will remove the entity, with all its components)
	Unlink(entityId entity.Id)

	// EntityMask method - retrieves the linked entities (as a bitset)
	EntityMask() bit.Mask

	// GetScheduledRemoves method - retrieves the scheduled remove
	GetScheduledRemoves() bit.Mask

	// Update method - updates the linked entities by removing the scheduled removes
	Update()

	// IsClear method - returns true if the entity linker doesn't need to be updated, false otherwise
	IsClear() bool

	// Refresh method - clears the scheduled removes and the flag (marking the entity linker as updated
	Refresh()
}
