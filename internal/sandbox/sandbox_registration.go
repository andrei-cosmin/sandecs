package sandbox

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/api"
	internalComponent "github.com/andrei-cosmin/sandecs/internal/component"
)

// ComponentRegistration struct - holds the registration of the component type T
//   - linker component.Linker[T] - the linker for the component type T
type ComponentRegistration[T component.Component] struct {
	linker component.Linker[T]
}

// Execute method - executes the registration (returns the linker for the component type T and stores it inside the registration)
func (r *ComponentRegistration[T]) Execute(context api.ComponentLinkManager) {
	r.linker = internalComponent.RegisterComponentLinker[T](context).(component.Linker[T])
}

// GetLinker method - returns the linker for the component type T
func (r *ComponentRegistration[T]) GetLinker() component.Linker[T] {
	return r.linker
}

// TagRegistration struct - holds the registration of a tag
//   - tag component.Tag - the tag to register
//   - linker component.TagLinker - the linker for the tag
type TagRegistration struct {
	tag    component.Tag
	linker component.TagLinker
}

// NewTagRegistration method - creates a new tag registration with the given tag
func NewTagRegistration(tag component.Tag) *TagRegistration {
	return &TagRegistration{
		tag: tag,
	}
}

// Execute method - executes the registration (returns the linker for the tag and stores it inside the registration)
func (r *TagRegistration) Execute(context api.ComponentLinkManager) {
	r.linker = internalComponent.RegisterTagLinker(r.tag, context).(component.TagLinker)
}

// GetLinker method - returns the linker for the tag
func (r *TagRegistration) GetLinker() component.TagLinker {
	return r.linker
}
