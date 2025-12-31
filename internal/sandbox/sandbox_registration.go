package sandbox

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/api"
	internalComponent "github.com/andrei-cosmin/sandecs/internal/component"
)

// ComponentRegistration holds the linker for component type T.
type ComponentRegistration[T component.Component] struct {
	linker component.Linker[T]
}

// Execute registers the component and stores the linker.
func (r *ComponentRegistration[T]) Execute(context api.ComponentLinkManager) {
	r.linker = internalComponent.RegisterComponentLinker[T](context).(component.Linker[T])
}

// GetLinker returns the linker for component type T.
func (r *ComponentRegistration[T]) GetLinker() component.Linker[T] {
	return r.linker
}

// TagRegistration holds the linker for a tag.
type TagRegistration struct {
	tag    component.Tag
	linker component.TagLinker
}

// NewTagRegistration creates a tag registration.
func NewTagRegistration(tag component.Tag) *TagRegistration {
	return &TagRegistration{tag: tag}
}

// Execute registers the tag and stores the linker.
func (r *TagRegistration) Execute(context api.ComponentLinkManager) {
	r.linker = internalComponent.RegisterTagLinker(r.tag, context).(component.TagLinker)
}

// GetLinker returns the linker for the tag.
func (r *TagRegistration) GetLinker() component.TagLinker {
	return r.linker
}
