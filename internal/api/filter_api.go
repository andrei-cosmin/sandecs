package api

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
)

// FilterRules defines component requirements for filtering entities.
type FilterRules interface {
	// RequiredComponentIds returns component IDs that must be present.
	RequiredComponentIds() []component.Id
	// ExcludedComponentIds returns component IDs that must be absent.
	ExcludedComponentIds() []component.Id
	// UnionComponentIds returns component IDs where at least one must be present.
	UnionComponentIds() []component.Id
}

// FilterRegistry manages filter registration and cached results.
type FilterRegistry interface {
	// Register creates a filter view from the given rules.
	Register(filter FilterRules) entity.View

	// UpdateLinks refreshes all filter caches.
	UpdateLinks()
}
