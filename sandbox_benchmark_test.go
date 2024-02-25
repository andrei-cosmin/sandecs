package sandbox

import (
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/filter"
	"testing"
)

const (
	numPosition         = 9000
	numPositionVelocity = 1000
)

func BenchmarkSandIter(b *testing.B) {
	b.StopTimer()
	sandbox := NewDefault()
	positionHandler := ComponentLinker[position](sandbox)
	velocityHandler := ComponentLinker[velocity](sandbox)

	for range numPosition {
		id := LinkEntity(sandbox)
		positionHandler.Link(id)
	}
	for range numPositionVelocity {
		id := LinkEntity(sandbox)
		positionHandler.Link(id)
		velocityHandler.Link(id)
	}

	view := Filter(sandbox, filter.Match2[position, velocity]())
	Update(sandbox)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, entityId := range view.EntityIds() {
			pos := positionHandler.Get(entityId)
			vel := velocityHandler.Get(entityId)
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}

func BenchmarkSandLinkUnlink(b *testing.B) {
	b.StopTimer()
	sandbox := NewDefault()
	positionHandler := ComponentLinker[position](sandbox)
	velocityHandler := ComponentLinker[velocity](sandbox)
	nameHandler := ComponentLinker[name](sandbox)

	Filter(sandbox, filter.Match[position]())
	Filter(sandbox, filter.Match2[position, velocity]())
	Filter(sandbox, filter.Match2[position, name]())

	for range numPosition {
		id := LinkEntity(sandbox)
		positionHandler.Link(id)
	}
	for range numPositionVelocity {
		id := LinkEntity(sandbox)
		positionHandler.Link(id)
		velocityHandler.Link(id)
	}

	Update(sandbox)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for entityId := range numPosition {
			nameHandler.Link(entity.Id(entityId))
		}
		Update(sandbox)
		for entityId := range numPosition {
			nameHandler.Unlink(entity.Id(entityId))
		}
		Update(sandbox)
	}
}
