package component

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
)

// tagLinker struct - manages the linking of entities with a specific component type
//   - baseLinker - the base linker which holds the common linking functionality
//   - onLink func() - a hook function which is called when a component instance is linked to an entity
//   - onUnlink func() - a hook function which is called when a component instance is unlinked from an entity
type tagLinker struct {
	baseLinker
	onLink   func()
	onUnlink func()
}

// newTagLinker method - creates a new tagLinker with the given parameters
func newTagLinker(size uint, componentId component.Id, componentType string, entityLinker entity.MaskView, callback func()) *tagLinker {
	return &tagLinker{
		baseLinker: *newBaseLinker(size, componentId, componentType, entityLinker, callback),
	}
}

// Link method - links an entity to the tag
func (r *tagLinker) Link(entityId entity.Id) bool {
	// Link the entity to the tag
	successfullyLinked := r.baseLinker.Link(entityId)

	// Call the onLink hook if the entity was successfully linked
	if successfullyLinked && r.onLink != nil {
		r.onLink()
	}

	// Return the result of the linking operation
	return successfullyLinked
}

// SetLinkHook method - sets a hook that will be called when an entity is linked to the tag
func (r *tagLinker) SetLinkHook(onLink func()) {
	r.onLink = onLink
}

// SetUnlinkHook method - sets a hook that will be called when an entity is unlinked from the tag
func (r *tagLinker) SetUnlinkHook(onUnlink func()) {
	r.onUnlink = onUnlink
}

// RemoveLinkHook method - removes the link hook
func (r *tagLinker) RemoveLinkHook(onLink func()) {
	r.onLink = nil
}

// RemoveUnlinkHook method - removes the unlink hook
func (r *tagLinker) RemoveUnlinkHook(onUnlink func()) {
	r.onUnlink = nil
}

// CleanScheduledInstances method - clears the instances corresponding to the scheduled entity removals
//
// NOTE: In the case of the tagLinker, it will only call the listener OnRemove hook
func (r *tagLinker) CleanScheduledInstances() {
	if r.onUnlink != nil {
		for range r.scheduledRemoves.Bits.Count() {
			r.onUnlink()
		}
	}
}
