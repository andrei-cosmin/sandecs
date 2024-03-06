package component

import (
	"github.com/andrei-cosmin/sandata/bit"
	"github.com/andrei-cosmin/sandecs/entity"
)

// Id type - unique identifier for a component
type Id = uint

// Component type - used for components that can be linked to an entity
type Component = interface{}

// Tag type - label that can be linked to an entity
type Tag = string

// BasicLinker interface - used for linking components to entities
//   - Has method - checks if the entity has the component
//   - Unlink method - unlinks the component / tag from the entity
//   - ComponentId method - returns the component id
type BasicLinker interface {

	// Has method - checks if the entity has the component / tag
	Has(entity entity.Id) bool

	// Unlink method - unlinks the component / tag from the entity
	//
	// WARNING: unlinking the same component / tag twice will have no additional effect (no panic will be thrown, no side effects will occur)
	Unlink(entity entity.Id) bool

	// ComponentId method - returns the component id
	ComponentId() Id
}

// Linker interface - used for linking components to entities, as well as retrieving specific instances of components
type Linker[T Component] interface {

	// Link method - links the component to the entity
	//
	// WARNING: linking the same component twice will return a nil pointer (no panic will be thrown, no side effects will occur)
	Link(entity entity.Id) *T

	// Get method - retrieves the component instance linked to the entity (if it doesn't exist, it will return nil)
	Get(entity entity.Id) *T

	// Has method - checks if the entity has the component
	Has(entity entity.Id) bool

	// Unlink method - unlinks the component from the entity
	//
	// WARNING: unlinking the same component instance twice will have no additional effect (no panic will be thrown, no side effects will occur)
	Unlink(entity entity.Id) bool

	// SetLinkHook method - sets a hook that will be called when an entity is linked to the component
	SetLinkHook(onLink func(*T))

	// SetUnlinkHook method - sets a hook that will be called when an entity is unlinked from the component
	SetUnlinkHook(onUnlink func(*T))

	// RemoveLinkHook method - removes the link hook
	RemoveLinkHook(onLink func())

	// RemoveUnlinkHook method - removes the unlink hook
	RemoveUnlinkHook(onUnlink func())

	// EntityMask method - returns the entity mask corresponding to the component type T
	EntityMask() bit.Mask

	// ComponentId method - returns the component id
	ComponentId() Id
}

// TagLinker interface - used for linking tags to entities
type TagLinker interface {

	// Link method - links the tag to the entity (will return false if linking same tag twice)
	//
	// WARNING: linking the same tag twice will have no additional effect (no panic will be thrown, no side effects will occur)
	Link(entity entity.Id) bool

	// Has method - checks if the entity has the tag
	Has(entity entity.Id) bool

	// Unlink method - unlinks the tag from the entity
	//
	// WARNING: unlinking the same tag twice will have no additional effect (no panic will be thrown, no side effects will occur)
	Unlink(entity entity.Id) bool

	// SetLinkHook method - sets a hook that will be called when an entity is linked to the tag
	SetLinkHook(onLink func())

	// SetUnlinkHook method - sets a hook that will be called when an entity is unlinked from the tag
	SetUnlinkHook(onUnlink func())

	// RemoveLinkHook method - removes the link hook
	RemoveLinkHook(onLink func())

	// RemoveUnlinkHook method - removes the unlink hook
	RemoveUnlinkHook(onUnlink func())

	// EntityMask method - returns the entity mask corresponding to the tag
	EntityMask() bit.Mask

	// ComponentId method - returns the component id
	ComponentId() Id
}
