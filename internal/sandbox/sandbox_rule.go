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
	Registration() api.ComponentRegistration
	ComponentId() component.Id
}

type BasicRule[T component.Component] struct {
	ComponentRegistration[T]
	ruleType Type
}

func NewRule[T component.Component](ruleType Type) *BasicRule[T] {
	return &BasicRule[T]{
		ruleType: ruleType,
	}
}

func (r *BasicRule[T]) RuleType() Type {
	return r.ruleType
}

func (r *BasicRule[T]) ComponentId() component.Id {
	return r.GetLinker().ComponentId()
}

func (r *BasicRule[T]) Registration() api.ComponentRegistration {
	return &r.ComponentRegistration
}
