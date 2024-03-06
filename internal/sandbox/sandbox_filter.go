package sandbox

import (
	"github.com/andrei-cosmin/sandecs/component"
)

// filterRules struct - holds the rules for filtering components / tags
//   - match []component.Id - the required component ids
//   - exclude []component.Id - the excluded component ids
//   - union []component.Id - the union component ids
type filterRules struct {
	match   []component.Id
	exclude []component.Id
	union   []component.Id
}

// RequiredComponentIds method - returns the required component ids
func (f *filterRules) RequiredComponentIds() []component.Id {
	return f.match
}

// ExcludedComponentIds method - returns the excluded component ids
func (f *filterRules) ExcludedComponentIds() []component.Id {
	return f.exclude
}

// UnionComponentIds method - returns the union component ids
func (f *filterRules) UnionComponentIds() []component.Id {
	return f.union
}
