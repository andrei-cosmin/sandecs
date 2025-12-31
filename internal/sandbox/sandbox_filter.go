package sandbox

import (
	"github.com/andrei-cosmin/sandecs/component"
)

// filterRules holds component ID lists for filtering.
type filterRules struct {
	match   []component.Id
	exclude []component.Id
	union   []component.Id
}

// RequiredComponentIds returns required component IDs.
func (f *filterRules) RequiredComponentIds() []component.Id {
	return f.match
}

// ExcludedComponentIds returns excluded component IDs.
func (f *filterRules) ExcludedComponentIds() []component.Id {
	return f.exclude
}

// UnionComponentIds returns union component IDs.
func (f *filterRules) UnionComponentIds() []component.Id {
	return f.union
}
