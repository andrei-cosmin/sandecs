# sandecs

A fast, type-safe Entity Component System (ECS) library for Go.

## Install

```bash
go get github.com/andrei-cosmin/sandecs
```

## Quick Start

```go
package main

import (
	"github.com/andrei-cosmin/sandecs/sandbox"
	"github.com/andrei-cosmin/sandecs/filter"
	"github.com/andrei-cosmin/sandecs/options"
)

type Position struct{ X, Y float64 }
type Velocity struct{ X, Y float64 }

func main() {
	// Create a sandbox
	sb := sandbox.NewDefault()

	// Get component linkers
	pos := sandbox.ComponentLinker[Position](sb)
	vel := sandbox.ComponentLinker[Velocity](sb)

	// Create entity and attach components
	entity := sandbox.LinkEntity(sb)
	pos.Link(entity).X = 100
	vel.Link(entity).X = 10

	// Apply changes
	sandbox.Update(sb)

	// Query and iterate
	view := sandbox.Filter(sb, filter.Match2[Position, Velocity]())
	for _, id := range view.EntityIds() {
		p, v := pos.Get(id), vel.Get(id)
		p.X += v.X
		p.Y += v.Y
	}

	sandbox.Update(sb)
}
```

## Filters

```go
// Match entities with specific components
filter.Match[Position]()
filter.Match2[Position, Velocity]()
filter.Match3[A, B, C]() // up to Match5

// Exclude entities with components
filter.Exclude[Dead]()

// Match entities with any of these components
filter.Union2[Sprite, Mesh]()

// Tags (zero-storage labels)
rendered := sandbox.TagLinker(sb, "rendered")
rendered.Link(entity)
filter.MatchTags("rendered")
filter.ExcludeTags("disabled")

// Combine filters
view := sandbox.Filter(sb,
filter.Match2[Position, Velocity](),
filter.ExcludeTags("disabled"),
)
```

## Storage Options

```go
// Standard - dynamic arrays (default)
sb := sandbox.New(options.Standard, entityCap, componentCap, 0)

// Pooled - reuses deallocated slots
sb := sandbox.New(options.Pooled, entityCap, componentCap, poolSize)

// Compact - optimized for dense data
sb := sandbox.New(options.Compact, entityCap, componentCap, 0)
```

## Hooks

```go
posLinker := sandbox.ComponentLinker[Position](sb)
posLinker.SetLinkHook(func (p *Position) {
// called when component is linked
})
posLinker.SetUnlinkHook(func (p *Position) {
// called when component is unlinked
})
```

## License

MIT