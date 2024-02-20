package sandbox

import (
	"github.com/andrei-cosmin/hakkt-ecs/component"
	internalComponent "github.com/andrei-cosmin/hakkt-ecs/internal/component"
	internalEntity "github.com/andrei-cosmin/hakkt-ecs/internal/entity"
	internalQuery "github.com/andrei-cosmin/hakkt-ecs/internal/query"
	"github.com/andrei-cosmin/hakkt-ecs/marker"
	"github.com/andrei-cosmin/hakkt-ecs/query"
)

type Sandbox struct {
	componentLinker *internalComponent.LinkManager
	entityLinker    *internalEntity.Linker
	queryRegistry   *internalQuery.Registry
}

const sandboxDefaultSize = 128

func New() *Sandbox {
	var componentLinker = internalComponent.NewLinkManager(sandboxDefaultSize)
	var entityLinker = internalEntity.NewLinker(sandboxDefaultSize)
	var sandbox = &Sandbox{
		componentLinker: componentLinker,
		entityLinker:    entityLinker,
		queryRegistry:   internalQuery.NewRegistry(sandboxDefaultSize, entityLinker, componentLinker),
	}
	Update(sandbox)
	return sandbox
}

func Filter(s *Sandbox, query query.Query) query.Response {
	return s.queryRegistry.Resolve(query.Get())
}

func LinkEntity(s *Sandbox) uint {
	return s.entityLinker.Link()
}

func UnlinkEntity(s *Sandbox, entity uint) {
	s.entityLinker.Unlink(entity)
}

func GetComponentLinker[T component.Component](s *Sandbox) component.Linker[T] {
	var resolver = s.componentLinker.Link(marker.Symbol[T]())
	return &resolverWrapper[T]{
		internalResolver: resolver,
	}
}

func Update(s *Sandbox) {
	if s.componentLinker.IsUpdated() && s.entityLinker.IsUpdated() {
		return
	}
	s.componentLinker.UpdateLinks(s.entityLinker.GetScheduledRemoves())
	s.queryRegistry.UpdateLinks()
	s.entityLinker.Refresh()
}
