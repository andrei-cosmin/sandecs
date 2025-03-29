package entity

import (
	"github.com/andrei-cosmin/sandata/bit"
	"github.com/andrei-cosmin/sandata/flag"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/bits-and-blooms/bitset"
)

// Linker struct - entity linker (links entities to the sandbox)
//   - linkedEntities *data.BitMask - a bitset storing the linked entities
//   - scheduledRemoves *data.BitMask - a bitset storing the entities that are scheduled for removal
//   - Flag - a flag used to mark the linker for update
type Linker struct {
	linkedEntities   *bit.BitMask
	scheduledRemoves *bit.BitMask
	flag.Flag
}

// NewLinker method - creates a new linker with the given size (pre-allocates the bitsets)
func NewLinker(size uint) *Linker {
	return &Linker{
		linkedEntities:   bit.NewMask(bitset.New(size)),
		scheduledRemoves: bit.NewMask(bitset.New(size)),
		Flag:             flag.New(),
	}
}

// Link method - links a new entity with the sandbox, returning the entity id
func (l *Linker) Link() entity.Id {
	// Find the first clear bit in the linked entities bitset
	entityId, exists := l.linkedEntities.NextClear(0)
	// If the entity id does not exist, set it to the length of the linked entities bitset
	if !exists {
		entityId = l.linkedEntities.Len()
	}
	// Set the corresponding bit in the linked entities bitset
	l.linkedEntities.Bits.Set(entityId)

	// Return the entity id
	return entityId
}

// Unlink method - unlinks the entity id from the sandbox entirely (this effect will propagate to all the component linkers)
func (l *Linker) Unlink(entityId entity.Id) {
	// If the entity id is not part of the sandbox, return
	if !l.linkedEntities.Test(entityId) {
		return
	}

	// Mark the entity id as scheduled for removal
	l.scheduledRemoves.Bits.Set(entityId)

	// Flag the linker for update
	l.Set()
}

// EntityMask method - retrieves the linked entities (as a bitset)
func (l *Linker) EntityMask() bit.Mask {
	return l.linkedEntities
}

// GetScheduledRemoves method - retrieves the scheduled removes
func (l *Linker) GetScheduledRemoves() bit.Mask {
	return l.scheduledRemoves
}

// Update method - updates the linked entities by removing the scheduled removes
func (l *Linker) Update() {
	l.linkedEntities.Bits.InPlaceDifference(l.scheduledRemoves.Bits)
}

// Refresh method - clears the scheduled removes and the flag (marking the entity linker as updated
func (l *Linker) Refresh() {
	l.scheduledRemoves.Bits.ClearAll()
	l.Clear()
}
