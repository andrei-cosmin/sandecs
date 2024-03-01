package api

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
)

// FilterRules interface - used for defining rules for filtering entities
type FilterRules interface {
	// RequiredComponentIds returns a list of component ids that are required for an entity to be included in the filter
	RequiredComponentIds() []component.Id
	// ExcludedComponentIds returns a list of component ids that are excluded for an entity to be included in the filter
	ExcludedComponentIds() []component.Id
	// UnionComponentIds returns a list of component ids that are used for union filtering
	UnionComponentIds() []component.Id
}

// FilterRegistry interface - used for registering filter rules and updating the caches storing the filtered entities
type FilterRegistry interface {
	Register(filter FilterRules) entity.SliceView
	UpdateLinks()
}
