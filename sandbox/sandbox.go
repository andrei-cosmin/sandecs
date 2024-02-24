package sandbox

import (
	"github.com/andrei-cosmin/hakkt/component"
	"github.com/andrei-cosmin/hakkt/entity"
	"github.com/andrei-cosmin/hakkt/filter"
	"github.com/andrei-cosmin/hakkt/internal/sandbox"
)

const DefaultNumEntities = 128
const DefaultNumComponents = 16
const DefaultPoolCapacity = 1024

type Sandbox struct {
	internal *sandbox.Sandbox
}

func New(numEntities, numComponents, poolCapacity uint) *Sandbox {
	box := &Sandbox{
		internal: sandbox.New(numEntities, numComponents, poolCapacity),
	}
	Update(box)
	return box
}

func NewDefault() *Sandbox {
	return New(DefaultNumEntities, DefaultNumComponents, DefaultPoolCapacity)
}

func Filter(s *Sandbox, filters ...filter.Filter) entity.View {
	rules := make([]sandbox.Rule, 0)
	for index := range filters {
		rules = append(rules, filters[index].Rules...)
	}
	return sandbox.LinkFilter(s.internal, rules)
}

func LinkEntity(s *Sandbox) entity.Id {
	return s.internal.LinkEntity()
}

func UnlinkEntity(s *Sandbox, entityId entity.Id) {
	s.internal.UnlinkEntity(entityId)
}

func IsEntityLinked(s *Sandbox, entityId entity.Id) bool {
	return s.internal.IsEntityLinked(entityId)
}

func ComponentLinker[T component.Component](s *Sandbox) component.Linker[T] {
	registration := sandbox.ComponentRegistration[T]{}
	s.internal.Accept(&registration)
	return registration.GetLinker()
}

func Update(s *Sandbox) {
	if s.internal.IsUpdated() {
		return
	}
	s.internal.Update()
}
