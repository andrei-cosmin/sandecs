package component

import (
	"github.com/andrei-cosmin/sandata/bit"
	"github.com/andrei-cosmin/sandecs/entity"
)

// Id is a unique identifier for a component.
type Id = uint

// Component is any type that can be linked to an entity.
type Component = interface{}

// Tag is a label that can be linked to an entity (no backing storage).
type Tag = string

// BasicLinker provides basic linking operations for components and tags.
type BasicLinker interface {
	// Has returns true if the entity has this component/tag.
	Has(entity entity.Id) bool

	// Unlink removes the component/tag from the entity. Returns false if not linked.
	Unlink(entity entity.Id) bool

	// ComponentId returns the unique identifier for this component type.
	ComponentId() Id
}

// Linker manages component instances of type T for entities.
type Linker[T Component] interface {
	// Link attaches a component to the entity and returns a pointer to it.
	// Returns nil if already linked.
	Link(entity entity.Id) *T

	// Get returns the component for the entity, or nil if not linked.
	Get(entity entity.Id) *T

	// Has returns true if the entity has this component.
	Has(entity entity.Id) bool

	// Unlink removes the component from the entity. Returns false if not linked.
	Unlink(entity entity.Id) bool

	// SetLinkHook sets a callback invoked when a component is linked.
	SetLinkHook(onLink func(*T))

	// SetUnlinkHook sets a callback invoked when a component is unlinked.
	SetUnlinkHook(onUnlink func(*T))

	// RemoveLinkHook clears the link hook.
	RemoveLinkHook(onLink func())

	// RemoveUnlinkHook clears the unlink hook.
	RemoveUnlinkHook(onUnlink func())

	// EntityMask returns a bitmask of entities with this component.
	EntityMask() bit.Mask

	// ComponentId returns the unique identifier for this component type.
	ComponentId() Id
}

// TagLinker manages tag associations for entities (no data storage).
type TagLinker interface {
	// Link attaches the tag to the entity. Returns false if already linked.
	Link(entity entity.Id) bool

	// Has returns true if the entity has this tag.
	Has(entity entity.Id) bool

	// Unlink removes the tag from the entity. Returns false if not linked.
	Unlink(entity entity.Id) bool

	// SetLinkHook sets a callback invoked when a tag is linked.
	SetLinkHook(onLink func())

	// SetUnlinkHook sets a callback invoked when a tag is unlinked.
	SetUnlinkHook(onUnlink func())

	// RemoveLinkHook clears the link hook.
	RemoveLinkHook(onLink func())

	// RemoveUnlinkHook clears the unlink hook.
	RemoveUnlinkHook(onUnlink func())

	// EntityMask returns a bitmask of entities with this tag.
	EntityMask() bit.Mask

	// ComponentId returns the unique identifier for this tag.
	ComponentId() Id
}
