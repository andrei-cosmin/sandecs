package sandbox

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	internalComponent "github.com/andrei-cosmin/sandecs/internal/component"
	internalEntity "github.com/andrei-cosmin/sandecs/internal/entity"
	internalFilter "github.com/andrei-cosmin/sandecs/internal/filter"
	"slices"
)

// Sandbox struct - holds the internal sandbox for the ECS
//   - entityLinker api.EntityLinker - the entity linker
//   - componentLinkManager api.ComponentLinkManager - the component link manager
//   - filterRegistry api.FilterRegistry - the filter registry
type Sandbox struct {
	entityLinker         api.EntityLinker
	componentLinkManager api.ComponentLinkManager
	filterRegistry       api.FilterRegistry
}

// New method - creates a new sandbox with pre-allocated buffers and memory for number of entities, components and pool capacity
func New(numEntities, numComponents, poolCapacity uint) *Sandbox {
	// Create the entity linker
	entityLinker := internalEntity.NewLinker(numEntities)
	// Create the component link manager
	componentLinkManager := internalComponent.NewLinkManager(numEntities, numComponents, poolCapacity, entityLinker)
	// Create the filter registry
	filterRegistry := internalFilter.NewRegistry(numEntities, entityLinker, componentLinkManager)
	// Create the sandbox
	sandbox := &Sandbox{
		entityLinker:         entityLinker,
		componentLinkManager: componentLinkManager,
		filterRegistry:       filterRegistry,
	}
	// Return the sandbox
	return sandbox
}

// IsUpdated method - checks if the sandbox is updated (checks if the flags on the component link manager and entity linker are clear)
func (s *Sandbox) IsUpdated() bool {
	return s.componentLinkManager.IsClear() && s.entityLinker.IsClear()
}

// LinkEntity method - links an entity with the sandbox
func (s *Sandbox) LinkEntity() entity.Id {
	return s.entityLinker.Link()
}

// UnlinkEntity method - unlinks an entity from the sandbox entirely (this effect will propagate to all the component linkers)
func (s *Sandbox) UnlinkEntity(entityId entity.Id) {
	s.entityLinker.Unlink(entityId)
}

// IsEntityLinked method - checks if the entity id is linked with the sandbox
func (s *Sandbox) IsEntityLinked(entityId entity.Id) bool {
	return s.entityLinker.EntityIds().Test(entityId)
}

// Update method - updates the sandbox (updates the entity linker, component link manager and filter registry)
func (s *Sandbox) Update() {
	s.entityLinker.Update()
	s.componentLinkManager.UpdateLinks(s.entityLinker.GetScheduledRemoves())
	s.filterRegistry.UpdateLinks()
	s.entityLinker.Refresh()
}

// Accept method - accepts a registration, attaches a new linker or returns the existing linker
func (s *Sandbox) Accept(registration api.Registration) {
	s.componentLinkManager.Accept(registration)
}

// LinkFilter method - links a filter with the sandbox
func LinkFilter(s *Sandbox, rules []Rule) entity.SliceView {

	// Create the rule sets buffers
	ruleSets := make([][]component.Id, SetSize)
	for ruleSetIndex := range SetSize {
		ruleSets[ruleSetIndex] = make([]component.Id, 0)
	}

	// For each rule, accept the registration and append the component id to the corresponding rule set
	for _, rule := range rules {
		s.Accept(rule.Registration())
		ruleSets[rule.RuleType()] = append(ruleSets[rule.RuleType()], rule.ComponentId())
	}

	// Sort the rule sets buffers (used for filter deduplication)
	for ruleSetIndex := range SetSize {
		slices.Sort(ruleSets[ruleSetIndex])
	}

	// Register the filter with the filter registry and return the associated view
	return s.filterRegistry.Register(&filterRules{
		match:   ruleSets[Match],
		exclude: ruleSets[Exclude],
		one:     ruleSets[Union],
	})
}
