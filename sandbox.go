package sandbox

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/filter"
	"github.com/andrei-cosmin/sandecs/internal/sandbox"
	"github.com/andrei-cosmin/sandecs/options"
)

// Sandbox struct - the context in which the entities and components are linked
type Sandbox struct {
	internal *sandbox.Sandbox
}

// New method - creates a new sandbox with pre-allocated memory for numEntities entities and numComponents components
//   - numEntities uint - the number of entities to pre-allocate memory for
//   - numComponents uint - the number of components to pre-allocate memory for
//   - poolCapacity uint - the capacity of the component pools (component instances will be stored in memory and reused instead of being garbage collected)
//
// NOTE: pool capacity set to 0 will deactivate pooling altogether (basic array tables will be used to store component instances)
func New(mode options.Mode, numEntities, numComponents, poolCapacity uint) *Sandbox {
	return &Sandbox{
		internal: sandbox.New(mode, numEntities, numComponents, poolCapacity),
	}
}

// NewDefault method - creates a new sandbox with default values for numEntities, numComponents and poolCapacity
func NewDefault() *Sandbox {
	return New(options.Standard, options.DefaultNumEntities, options.DefaultNumComponents, options.DefaultPoolCapacity)
}

// Filter method - creates a view of the entities that pass the given filters
//
// WARNING: it is strongly recommended that all filters be registered at the beginning of the application's lifetime (in the initialization phase)
func Filter(s *Sandbox, filters ...filter.Filter) entity.View {
	rules := make([]sandbox.Rule, 0)
	for index := range filters {
		rules = append(rules, filters[index].Rules...)
	}
	return sandbox.LinkFilter(s.internal, rules)
}

// LinkEntity method - links a new entity to the sandbox and returns its id
func LinkEntity(s *Sandbox) entity.Id {
	return s.internal.LinkEntity()
}

// UnlinkEntity method - unlinks the entity with the given id from the sandbox (removing all its components)
//
// WARNING: unlinking the same component instance twice will have no additional effect (no panic will be thrown)
func UnlinkEntity(s *Sandbox, entityId entity.Id) {
	s.internal.UnlinkEntity(entityId)
}

// IsEntityLinked method - checks if the entity with the given id is linked to the sandbox
func IsEntityLinked(s *Sandbox, entityId entity.Id) bool {
	return s.internal.IsEntityLinked(entityId)
}

// ComponentLinker method - retrieves the linker for the given component type T
func ComponentLinker[T component.Component](s *Sandbox) component.Linker[T] {
	registration := sandbox.ComponentRegistration[T]{}
	s.internal.Accept(&registration)
	return registration.GetLinker()
}

// TagLinker method - retrieves the linker for the given tag
func TagLinker(s *Sandbox, tag component.Tag) component.TagLinker {
	registration := sandbox.NewTagRegistration(tag)
	s.internal.Accept(registration)
	return registration.GetLinker()
}

// Update method - updates the sandbox (updates the entity linker, component linkers and filter registry)
//
// NOTE: this method should be called at the end of the frame or before the draw calls (depending on the game loop implementation)
//
// WARNING: this method should be called only once per frame
//
// WARNING: no changes will be reflected upon the entities / components until this method is called (all updates are batched and performed at once)
func Update(s *Sandbox) {
	if s.internal.IsUpdated() {
		return
	}
	s.internal.Update()
}
