package api

import (
	"github.com/andrei-cosmin/hakkt/component"
	"github.com/andrei-cosmin/hakkt/entity"
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
