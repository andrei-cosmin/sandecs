package component

import (
	"github.com/andrei-cosmin/sandecs/entity"
)

type Id = uint

type Component = interface{}
type Tag = string

type BasicLinker interface {
	Has(entity entity.Id) bool
	Unlink(entity entity.Id)
	ComponentId() Id
}

type Linker[T Component] interface {
	Link(entity entity.Id) *T
	Get(entity entity.Id) *T
	Has(entity entity.Id) bool
	Unlink(entity entity.Id)
	ComponentId() Id
}
