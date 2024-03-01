package sandbox

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/api"
)

// Type - the type of the rule
type Type = uint8

// Rule types [ Match, Exclude, Union ]
const (
	Match = iota
	Exclude
	Union
	SetSize
	SetStart = Match
)

// Rule interface - the interface for a rule
//   - RuleType method - returns the type of the rule
//   - Registration method - returns the registration for the rule
//   - ComponentId method - returns the component id for the rule
type Rule interface {
	RuleType() Type
	Registration() api.Registration
	ComponentId() component.Id
}

// ComponentRule struct - holds the rule for a component and its registration
//   - ComponentRegistration[T] - the registration for the component type T
//   - ruleType Type - the type of the rule
type ComponentRule[T component.Component] struct {
	ComponentRegistration[T]
	ruleType Type
}

// NewComponentRule method - creates a new component rule with the given rule type
func NewComponentRule[T component.Component](ruleType Type) *ComponentRule[T] {
	return &ComponentRule[T]{
		ruleType: ruleType,
	}
}

// RuleType method - returns the type of the rule
func (r *ComponentRule[T]) RuleType() Type {
	return r.ruleType
}

// ComponentId method - returns the component id for the rule
func (r *ComponentRule[T]) ComponentId() component.Id {
	return r.GetLinker().ComponentId()
}

// Registration method - returns the registration for the rule
func (r *ComponentRule[T]) Registration() api.Registration {
	return &r.ComponentRegistration
}

// TagRule struct - holds the rule for a tag and its registration
//   - TagRegistration - the registration for the tag
//   - ruleType Type - the type of the rule
type TagRule struct {
	TagRegistration
	ruleType Type
}

// NewTagRule method - creates a new tag rule with the given rule type
func NewTagRule(tag component.Tag, ruleType Type) *TagRule {
	return &TagRule{
		TagRegistration: TagRegistration{
			tag: tag,
		},
		ruleType: ruleType,
	}
}

// RuleType method - returns the type of the rule
func (r *TagRule) RuleType() Type {
	return r.ruleType
}

// ComponentId method - returns the component id for the rule
func (r *TagRule) ComponentId() component.Id {
	return r.GetLinker().ComponentId()
}

// Registration method - returns the registration for the rule
func (r *TagRule) Registration() api.Registration {
	return &r.TagRegistration
}
