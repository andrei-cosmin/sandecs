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
	entityLinker         api.EntityLinker
	componentLinkManager api.ComponentLinkManager
	filterRegistry       api.FilterRegistry
}

func New(numEntities, numComponents, poolCapacity uint) *Sandbox {
	entityLinker := internalEntity.NewLinker(numEntities)
	componentLinkManager := internalComponent.NewLinkManager(numEntities, numComponents, poolCapacity, entityLinker)
	filterRegistry := internalFilter.NewRegistry(numEntities, entityLinker, componentLinkManager)
	sandbox := &Sandbox{
		entityLinker:         entityLinker,
		componentLinkManager: componentLinkManager,
		filterRegistry:       filterRegistry,
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

func (s *Sandbox) IsEntityLinked(entityId entity.Id) bool {
	return s.entityLinker.EntityIds().Test(entityId)
}

func (s *Sandbox) Update() {
	s.entityLinker.Update()
	s.componentLinkManager.UpdateLinks(s.entityLinker.GetScheduledRemoves())
	s.filterRegistry.UpdateLinks()
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

	return s.filterRegistry.Register(&filterRules{
		match:   ruleSets[Match],
		exclude: ruleSets[Exclude],
		one:     ruleSets[Union],
	})
}
