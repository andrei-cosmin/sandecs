package sandbox

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/api"
	internalComponent "github.com/andrei-cosmin/sandecs/internal/component"
)

type ComponentRegistration[T component.Component] struct {
	linker component.Linker[T]
}

func (r *ComponentRegistration[T]) Execute(context api.ComponentLinkManager) {
	r.linker = internalComponent.RegisterComponentLinker[T](context).(component.Linker[T])
}

func (r *ComponentRegistration[T]) GetLinker() component.Linker[T] {
	return r.linker
}

type TagRegistration struct {
	tag    component.Tag
	linker component.TagLinker
}

func NewTagRegistration(tag component.Tag) *TagRegistration {
	return &TagRegistration{
		tag: tag,
	}
}

func (r *TagRegistration) Execute(context api.ComponentLinkManager) {
	r.linker = internalComponent.RegisterTagLinker(r.tag, context).(component.TagLinker)
}

func (r *TagRegistration) GetLinker() component.TagLinker {
	return r.linker
}
