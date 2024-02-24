package component

import (
	"github.com/andrei-cosmin/hakkt/entity"
)

type Id = uint

type Component = interface{}

type BasicLinker interface {
	Link(entity entity.Id)
	Has(entity entity.Id) bool
	Unlink(entity entity.Id)
	ComponentId() Id
}

type Linker[T Component] interface {
	Link(entity entity.Id)
	Get(entity entity.Id) *T
	Has(entity entity.Id) bool
	Unlink(entity entity.Id)
	ComponentId() Id
}
