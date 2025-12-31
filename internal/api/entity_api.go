package api

import (
	"github.com/andrei-cosmin/sandata/bit"
	"github.com/andrei-cosmin/sandecs/entity"
)

// EntityLinker manages entity lifecycle in the sandbox.
type EntityLinker interface {
	// Link creates a new entity and returns its ID.
	Link() entity.Id

	// Unlink schedules entity removal (along with all its components).
	Unlink(entityId entity.Id)

	// EntityMask returns the bitmask of linked entities.
	EntityMask() bit.Mask

	// GetScheduledRemoves returns entities scheduled for removal.
	GetScheduledRemoves() bit.Mask

	// Update processes scheduled entity removals.
	Update()

	// IsCleared returns true if no pending updates exist.
	IsCleared() bool

	// Refresh clears scheduled removes after update.
	Refresh()
}
