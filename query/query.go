package query

import (
	"github.com/andrei-cosmin/hakkt/internal/query"
	"github.com/andrei-cosmin/hakkt/internal/state"
	"github.com/hashicorp/go-set"
)

type Query interface {
	Match(marker string) Query
	MatchAll(markers ...string) Query
	Exclude(marker string) Query
	ExcludeAll(markers ...string) Query
	OneOf(markers ...string) Query
	Get() *query.Info
}

type builder struct {
	match   *set.Set[string]
	exclude *set.Set[string]
	one     *set.Set[string]
	status  state.State
	info    *query.Info
}

const queryDefaultSize = 8

func New() Query {
	return &builder{
		match:   set.New[string](queryDefaultSize),
		exclude: set.New[string](queryDefaultSize),
		one:     set.New[string](queryDefaultSize),
		status:  state.New(),
	}
}

func (b *builder) Match(marker string) Query {
	b.match.Insert(marker)
	b.exclude.Remove(marker)
	b.one.Remove(marker)
	b.status.Mark()
	return b
}

func (b *builder) MatchAll(markers ...string) Query {
	for _, marker := range markers {
		b.match.Insert(marker)
		b.exclude.Remove(marker)
		b.one.Remove(marker)
	}
	b.status.Mark()
	return b
}

func (b *builder) Exclude(marker string) Query {
	b.exclude.Insert(marker)
	b.match.Remove(marker)
	b.one.Remove(marker)
	b.status.Mark()
	return b
}

func (b *builder) ExcludeAll(markers ...string) Query {
	for _, marker := range markers {
		b.exclude.Insert(marker)
		b.match.Remove(marker)
		b.one.Remove(marker)
	}
	b.status.Mark()
	return b
}

func (b *builder) OneOf(markers ...string) Query {
	for _, marker := range markers {
		b.one.Insert(marker)
		b.match.Remove(marker)
		b.exclude.Remove(marker)
	}
	b.status.Mark()
	return b
}

func (b *builder) Get() *query.Info {
	if !b.status.IsUpdated() {
		b.info = query.NewInfo(b.match, b.exclude, b.one)
		b.status.Reset()
	}
	return b.info
}
