package sandbox

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/api"
)

// Type represents a filter rule type.
type Type = uint8

// Rule types.
const (
	Match = iota
	Exclude
	Union
	SetSize
	SetStart = Match
)

// Rule defines a filter rule.
type Rule interface {
	RuleType() Type
	Registration() api.Registration
	ComponentId() component.Id
}

// ComponentRule is a filter rule for a component type.
type ComponentRule[T component.Component] struct {
	ComponentRegistration[T]
	ruleType Type
}

// NewComponentRule creates a component rule.
func NewComponentRule[T component.Component](ruleType Type) *ComponentRule[T] {
	return &ComponentRule[T]{ruleType: ruleType}
}

// RuleType returns the rule type.
func (r *ComponentRule[T]) RuleType() Type {
	return r.ruleType
}

// ComponentId returns the component ID.
func (r *ComponentRule[T]) ComponentId() component.Id {
	return r.GetLinker().ComponentId()
}

// Registration returns the component registration.
func (r *ComponentRule[T]) Registration() api.Registration {
	return &r.ComponentRegistration
}

// TagRule is a filter rule for a tag.
type TagRule struct {
	TagRegistration
	ruleType Type
}

// NewTagRule creates a tag rule.
func NewTagRule(tag component.Tag, ruleType Type) *TagRule {
	return &TagRule{
		TagRegistration: TagRegistration{tag: tag},
		ruleType:        ruleType,
	}
}

// RuleType returns the rule type.
func (r *TagRule) RuleType() Type {
	return r.ruleType
}

// ComponentId returns the tag's component ID.
func (r *TagRule) ComponentId() component.Id {
	return r.GetLinker().ComponentId()
}

// Registration returns the tag registration.
func (r *TagRule) Registration() api.Registration {
	return &r.TagRegistration
}
