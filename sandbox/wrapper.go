package sandbox

import (
	"github.com/andrei-cosmin/hakkt-ecs/component"
	internalComponent "github.com/andrei-cosmin/hakkt-ecs/internal/component"
	"unsafe"
)

type resolverWrapper[T component.Component] struct {
	internalResolver *internalComponent.Linker
}

func (r *resolverWrapper[T]) Link(entity uint, instance T) {
	r.internalResolver.Link(entity, (*component.Component)(unsafe.Pointer(&instance)))
}

func (r *resolverWrapper[T]) Get(entity uint) *T {
	return (*T)(unsafe.Pointer(r.internalResolver.Get(entity)))
}

func (r *resolverWrapper[T]) Has(entity uint) bool {
	return r.internalResolver.Has(entity)
}

func (r *resolverWrapper[T]) Remove(entity uint) {
	r.internalResolver.Remove(entity)
}
