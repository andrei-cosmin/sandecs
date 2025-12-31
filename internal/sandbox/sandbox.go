package sandbox

import (
	"slices"

	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	internalComponent "github.com/andrei-cosmin/sandecs/internal/component"
	internalEntity "github.com/andrei-cosmin/sandecs/internal/entity"
	internalFilter "github.com/andrei-cosmin/sandecs/internal/filter"
	"github.com/andrei-cosmin/sandecs/options"
)

// Sandbox is the internal ECS container.
type Sandbox struct {
	entityLinker         api.EntityLinker
	componentLinkManager api.ComponentLinkManager
	filterRegistry       api.FilterRegistry
}

// New creates a sandbox with pre-allocated capacity.
func New(mode options.Mode, numEntities, numComponents, poolCapacity uint) *Sandbox {
	entityLinker := internalEntity.NewLinker(numEntities)
	componentLinkManager := internalComponent.NewLinkManager(mode, numEntities, numComponents, poolCapacity, entityLinker)
	filterRegistry := internalFilter.NewRegistry(numEntities, entityLinker, componentLinkManager)
	return &Sandbox{
		entityLinker:         entityLinker,
		componentLinkManager: componentLinkManager,
		filterRegistry:       filterRegistry,
	}
}

// IsUpdated returns true if no pending updates exist.
func (s *Sandbox) IsUpdated() bool {
	return s.componentLinkManager.IsCleared() && s.entityLinker.IsCleared()
}

// LinkEntity creates a new entity and returns its ID.
func (s *Sandbox) LinkEntity() entity.Id {
	return s.entityLinker.Link()
}

// UnlinkEntity schedules entity removal.
func (s *Sandbox) UnlinkEntity(entityId entity.Id) {
	s.entityLinker.Unlink(entityId)
}

// IsEntityLinked returns true if the entity exists.
func (s *Sandbox) IsEntityLinked(entityId entity.Id) bool {
	return s.entityLinker.EntityMask().Test(entityId)
}

// Update processes all pending changes.
func (s *Sandbox) Update() {
	s.entityLinker.Update()
	s.componentLinkManager.UpdateLinks(s.entityLinker.GetScheduledRemoves())
	s.filterRegistry.UpdateLinks()
	s.entityLinker.Refresh()
}

// Accept processes a component registration.
func (s *Sandbox) Accept(registration api.Registration) {
	s.componentLinkManager.Accept(registration)
}

// LinkFilter creates a filter view from the given rules.
func LinkFilter(s *Sandbox, rules []Rule) entity.View {
	ruleSets := make([][]component.Id, SetSize)
	for ruleSetIndex := range SetSize {
		ruleSets[ruleSetIndex] = make([]component.Id, 0)
	}

	for _, rule := range rules {
		s.Accept(rule.Registration())
		ruleSets[rule.RuleType()] = append(ruleSets[rule.RuleType()], rule.ComponentId())
	}

	for ruleSetIndex := range SetSize {
		slices.Sort(ruleSets[ruleSetIndex])
	}

	return s.filterRegistry.Register(&filterRules{
		match:   ruleSets[Match],
		exclude: ruleSets[Exclude],
		union:   ruleSets[Union],
	})
}
