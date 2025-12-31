package sandbox

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/filter"
	"github.com/andrei-cosmin/sandecs/internal/sandbox"
	"github.com/andrei-cosmin/sandecs/options"
)

// Sandbox is the ECS container for entities and components.
type Sandbox struct {
	internal *sandbox.Sandbox
}

// New creates a sandbox with pre-allocated capacity.
func New(mode options.Mode, numEntities, numComponents, poolCapacity uint) *Sandbox {
	return &Sandbox{
		internal: sandbox.New(mode, numEntities, numComponents, poolCapacity),
	}
}

// NewDefault creates a sandbox with default configuration.
func NewDefault() *Sandbox {
	return New(options.Standard, options.DefaultNumEntities, options.DefaultNumComponents, options.DefaultPoolCapacity)
}

// Filter creates a view of entities matching the given filters.
// Register filters during initialization.
func Filter(s *Sandbox, filters ...filter.Filter) entity.View {
	rules := make([]sandbox.Rule, 0)
	for index := range filters {
		rules = append(rules, filters[index].Rules...)
	}
	return sandbox.LinkFilter(s.internal, rules)
}

// LinkEntity creates a new entity and returns its ID.
func LinkEntity(s *Sandbox) entity.Id {
	return s.internal.LinkEntity()
}

// UnlinkEntity removes the entity and all its components.
func UnlinkEntity(s *Sandbox, entityId entity.Id) {
	s.internal.UnlinkEntity(entityId)
}

// IsEntityLinked returns true if the entity exists.
func IsEntityLinked(s *Sandbox, entityId entity.Id) bool {
	return s.internal.IsEntityLinked(entityId)
}

// ComponentLinker returns the linker for component type T.
func ComponentLinker[T component.Component](s *Sandbox) component.Linker[T] {
	registration := sandbox.ComponentRegistration[T]{}
	s.internal.Accept(&registration)
	return registration.GetLinker()
}

// TagLinker returns the linker for the given tag.
func TagLinker(s *Sandbox, tag component.Tag) component.TagLinker {
	registration := sandbox.NewTagRegistration(tag)
	s.internal.Accept(registration)
	return registration.GetLinker()
}

// Update processes all pending changes. Call once per frame.
func Update(s *Sandbox) {
	if s.internal.IsUpdated() {
		return
	}
	s.internal.Update()
}
