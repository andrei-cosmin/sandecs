package entity

import (
	"github.com/andrei-cosmin/sandata/flag"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/bits-and-blooms/bitset"
)

// Linker struct - entity linker (links entities to the sandbox)
//   - linkedEntities *bitset.Bitset - a bitset storing the linked entities
//   - scheduledRemoves *bitset.Bitset - a bitset storing the entities that are scheduled for removal
//   - Flag - a flag used to mark the linker for update
type Linker struct {
	linkedEntities   *bitset.BitSet
	scheduledRemoves *bitset.BitSet
	flag.Flag
}

// NewLinker method - creates a new linker with the given size (pre-allocates the bitsets)
func NewLinker(size uint) *Linker {
	return &Linker{
		linkedEntities:   bitset.New(size),
		scheduledRemoves: bitset.New(size),
		Flag:             flag.New(),
	}
}

// EntityIds method - retrieves the linked entities (as a bitset)
func (l *Linker) EntityIds() *bitset.BitSet {
	return l.linkedEntities
}

// Link method - links the entity id with the sandbox
func (l *Linker) Link() entity.Id {
	// Find the first clear bit in the linked entities bitset
	entityId, exists := l.linkedEntities.NextClear(0)
	// If the entity id does not exist, set it to the length of the linked entities bitset
	if !exists {
		entityId = l.linkedEntities.Len()
	}
	// Set the corresponding bit in the linked entities bitset
	l.linkedEntities.Set(entityId)

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
	l.scheduledRemoves.Set(entityId)

	// Flag the linker for update
	l.Set()
}

// GetScheduledRemoves method - retrieves the scheduled removes
func (l *Linker) GetScheduledRemoves() *bitset.BitSet {
	return l.scheduledRemoves
}

// Update method - updates the linked entities by removing the scheduled removes
func (l *Linker) Update() {
	l.linkedEntities.InPlaceDifference(l.scheduledRemoves)
}

// Refresh method - clears the scheduled removes and the flag (marking the entity linker as updated
func (l *Linker) Refresh() {
	l.scheduledRemoves.ClearAll()
	l.Clear()
}
