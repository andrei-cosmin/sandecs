package sandbox

import (
	"github.com/andrei-cosmin/hakkt/component"
	"github.com/andrei-cosmin/hakkt/entity"
	"github.com/andrei-cosmin/hakkt/internal/api"
	internalComponent "github.com/andrei-cosmin/hakkt/internal/component"
	internalEntity "github.com/andrei-cosmin/hakkt/internal/entity"
	internalFilter "github.com/andrei-cosmin/hakkt/internal/filter"
)

type Sandbox struct {
	componentLinkManager api.ComponentLinkManager
	entityLinker         api.EntityLinker
	queryRegistry        api.FilterRegistry
}

func New(size uint) *Sandbox {
	componentLinkManager := internalComponent.NewLinkManager(size)
	entityLinker := internalEntity.NewLinker(size)
	sandbox := &Sandbox{
		componentLinkManager: componentLinkManager,
		entityLinker:         entityLinker,
		queryRegistry:        internalFilter.NewRegistry(size, entityLinker, componentLinkManager),
	}
	return sandbox
}

func (s *Sandbox) IsUpdated() bool {
	return s.componentLinkManager.IsUpdated() && s.entityLinker.IsUpdated()
}

func (s *Sandbox) LinkEntity() entity.Id {
	return s.entityLinker.Link()
}

func (s *Sandbox) UnlinkEntity(entityId entity.Id) {
	s.entityLinker.Unlink(entityId)
}

func (s *Sandbox) Update() {
	s.componentLinkManager.UpdateLinks(s.entityLinker.GetScheduledRemoves())
	s.queryRegistry.UpdateLinks()
	s.entityLinker.Refresh()
}

func (s *Sandbox) Accept(registration api.ComponentRegistration) {
	s.componentLinkManager.Accept(registration)
}

func LinkFilter(s *Sandbox, rules []Rule) entity.View {
	ruleSets := make([][]component.Id, SetSize)

	for ruleSetIndex := SetStart; ruleSetIndex < SetSize; ruleSetIndex++ {
		ruleSets[ruleSetIndex] = make([]component.Id, 0)
	}

	for _, rule := range rules {
		s.Accept(rule.Registration())
		ruleSets[rule.RuleType()] = append(ruleSets[rule.RuleType()], rule.ComponentId())
	}

	return s.queryRegistry.Register(&filterRules{
		match:   ruleSets[Match],
		exclude: ruleSets[Exclude],
		one:     ruleSets[Union],
	})
}
