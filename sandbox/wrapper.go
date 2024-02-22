package sandbox

import (
	"github.com/andrei-cosmin/hakkt/component"
	internalComponent "github.com/andrei-cosmin/hakkt/internal/component"
	"unsafe"
)

type resolverWrapper[T component.Component] struct {
	internalResolver *internalComponent.Linker
}

func (r *resolverWrapper[T]) Link(entityId uint, instance T) {
	r.internalResolver.Link(entityId, (*component.Component)(unsafe.Pointer(&instance)))
}

func (r *resolverWrapper[T]) Get(entityId uint) *T {
	return (*T)(unsafe.Pointer(r.internalResolver.Get(entityId)))
}

func (r *resolverWrapper[T]) Has(entityId uint) bool {
	return r.internalResolver.Has(entityId)
}

func (r *resolverWrapper[T]) Remove(entityId uint) {
	r.internalResolver.Remove(entityId)
}
