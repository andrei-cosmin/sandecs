package component

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
)

// tagLinker manages tag associations for entities (no data storage).
type tagLinker struct {
	baseLinker
	onLink   func()
	onUnlink func()
}

func newTagLinker(size uint, componentId component.Id, componentType string, entityLinker entity.MaskView, callback func()) *tagLinker {
	return &tagLinker{
		baseLinker: *newBaseLinker(size, componentId, componentType, entityLinker, callback),
	}
}

// Link attaches the tag to the entity. Returns false if already linked.
func (r *tagLinker) Link(entityId entity.Id) bool {
	if r.baseLinker.Link(entityId) {
		if r.onLink != nil {
			r.onLink()
		}
		return true
	}
	return false
}

// SetLinkHook sets a callback invoked when a tag is linked.
func (r *tagLinker) SetLinkHook(onLink func()) {
	r.onLink = onLink
}

// SetUnlinkHook sets a callback invoked when a tag is unlinked.
func (r *tagLinker) SetUnlinkHook(onUnlink func()) {
	r.onUnlink = onUnlink
}

// RemoveLinkHook clears the link hook.
func (r *tagLinker) RemoveLinkHook(onLink func()) {
	r.onLink = nil
}

// RemoveUnlinkHook clears the unlink hook.
func (r *tagLinker) RemoveUnlinkHook(onUnlink func()) {
	r.onUnlink = nil
}

// CleanScheduledInstances triggers unlink hooks for scheduled removals.
func (r *tagLinker) CleanScheduledInstances() {
	if r.onUnlink != nil {
		for range r.scheduledRemoves.Bits().Count() {
			r.onUnlink()
		}
	}
}
