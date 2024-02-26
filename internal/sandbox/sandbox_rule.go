package sandbox

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/api"
)

type Type = uint8

const (
	Match = iota
	Exclude
	Union
	SetSize
	SetStart = Match
)

type Rule interface {
	RuleType() Type
	Registration() api.Registration
	ComponentId() component.Id
}

type ComponentRule[T component.Component] struct {
	ComponentRegistration[T]
	ruleType Type
}

func NewComponentRule[T component.Component](ruleType Type) *ComponentRule[T] {
	return &ComponentRule[T]{
		ruleType: ruleType,
	}
}

func (r *ComponentRule[T]) RuleType() Type {
	return r.ruleType
}

func (r *ComponentRule[T]) ComponentId() component.Id {
	return r.GetLinker().ComponentId()
}

func (r *ComponentRule[T]) Registration() api.Registration {
	return &r.ComponentRegistration
}

type TagRule struct {
	TagRegistration
	ruleType Type
}

func NewTagRule(tag component.Tag, ruleType Type) *TagRule {
	return &TagRule{
		TagRegistration: TagRegistration{
			tag: tag,
		},
		ruleType: ruleType,
	}
}

func (r *TagRule) RuleType() Type {
	return r.ruleType
}

func (r *TagRule) ComponentId() component.Id {
	return r.GetLinker().ComponentId()
}

func (r *TagRule) Registration() api.Registration {
	return &r.TagRegistration
}
