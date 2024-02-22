package sandbox

import (
	"github.com/andrei-cosmin/hakkt/component"
)

type filterRules struct {
	match   []component.Id
	exclude []component.Id
	one     []component.Id
}

func (f *filterRules) RequiredComponentIds() []component.Id {
	return f.match
}

func (f *filterRules) ExcludedComponentIds() []component.Id {
	return f.exclude
}

func (f *filterRules) UnionComponentIds() []component.Id {
	return f.one
}
