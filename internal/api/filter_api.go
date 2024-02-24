package api

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
)

type FilterRules interface {
	RequiredComponentIds() []component.Id
	ExcludedComponentIds() []component.Id
	UnionComponentIds() []component.Id
}

type FilterRegistry interface {
	Register(filter FilterRules) entity.View
	UpdateLinks()
}
