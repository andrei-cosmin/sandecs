package sandbox

import (
	"github.com/andrei-cosmin/hakkt/component"
	"github.com/andrei-cosmin/hakkt/internal/api"
	internalComponent "github.com/andrei-cosmin/hakkt/internal/component"
)

type ComponentRegistration[T component.Component] struct {
	linker component.Linker[T]
}

func (r *ComponentRegistration[T]) Execute(context api.ComponentLinkManager) {
	r.linker = internalComponent.RegisterLinker[T](context).(component.Linker[T])
}

func (r *ComponentRegistration[T]) GetLinker() component.Linker[T] {
	return r.linker
}
