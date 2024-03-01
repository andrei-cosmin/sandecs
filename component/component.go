package component

import (
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
//   - Unlink method - unlinks the component from the entity
//   - ComponentId method - returns the component id
type BasicLinker interface {

	// Has method - checks if the entity has the component / tag
	Has(entity entity.Id) bool

	// Unlink method - unlinks the component / tag from the entity
	//
	// WARNING: unlinking the same component / tag twice will have no additional effect (no panic will be thrown, no side effects will occur)
	Unlink(entity entity.Id)

	// ComponentId method - returns the component id
	ComponentId() Id
}

// Linker interface - used for linking components to entities, as well as retrieving specific instances of components
type Linker[T Component] interface {

	// Link method - links the component to the entity
	//
	// WARNING: linking the same component instance twice will return a nil pointer (no panic will be thrown, no side effects will occur)
	Link(entity entity.Id) *T

	// Get method - retrieves the component instance linked to the entity (if it doesn't exist, it will return nil)
	Get(entity entity.Id) *T

	// Has method - checks if the entity has the component
	Has(entity entity.Id) bool

	// Unlink method - unlinks the component from the entity
	//
	// WARNING: unlinking the same component instance twice will have no additional effect (no panic will be thrown, no side effects will occur)
	Unlink(entity entity.Id)

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
	Unlink(entity entity.Id)

	// ComponentId method - returns the component id
	ComponentId() Id
}
