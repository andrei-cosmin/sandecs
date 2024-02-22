package sandbox

import (
	"github.com/andrei-cosmin/hakkt/component"
	"github.com/andrei-cosmin/hakkt/entity"
	"github.com/andrei-cosmin/hakkt/filter"
	"github.com/andrei-cosmin/hakkt/internal/sandbox"
)

const sandboxDefaultSize = 128

type Sandbox struct {
	internal *sandbox.Sandbox
}

func New() *Sandbox {
	box := &Sandbox{
		internal: sandbox.New(sandboxDefaultSize),
	}
	Update(box)
	return box
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
