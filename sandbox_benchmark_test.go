package sandbox

import (
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/filter"
	"testing"
)

const (
	numPosition         = 9000
	numPositionVelocity = 1000
	numEdits            = 1000
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

	for range b.N {
		for _, entityId := range view.EntityIds() {
			pos := positionHandler.Get(entityId)
			vel := velocityHandler.Get(entityId)
			pos.X += vel.X
			pos.Y += vel.Y
		}
		Update(sandbox)
	}
}

func BenchmarkSandLinkUnlinkComponent(b *testing.B) {
	b.StopTimer()
	sandbox := New(numPosition+numPositionVelocity, 4, 20000)
	positionHandler := ComponentLinker[position](sandbox)
	velocityHandler := ComponentLinker[velocity](sandbox)
	nameHandler := ComponentLinker[name](sandbox)

	for range numPosition {
		id := LinkEntity(sandbox)
		positionHandler.Link(id)
	}
	for range numPositionVelocity {
		id := LinkEntity(sandbox)
		positionHandler.Link(id)
		velocityHandler.Link(id)
	}

	Filter(sandbox, filter.Match2[position, velocity]())
	Filter(sandbox, filter.Match2[position, armor]())
	Filter(sandbox, filter.Match2[position, rendered]())
	Update(sandbox)
	b.StartTimer()

	for range b.N / 2 {
		for entityId := range numEdits {
			nameHandler.Link(entity.Id(entityId))
		}
		Update(sandbox)
		for entityId := range numEdits {
			nameHandler.Unlink(entity.Id(entityId))
		}
		Update(sandbox)
	}
}

func BenchmarkSandLinkUnlinkEntity(b *testing.B) {
	b.StopTimer()
	sandbox := New(numPosition+numPositionVelocity, 4, 20000)
	positionHandler := ComponentLinker[position](sandbox)
	velocityHandler := ComponentLinker[velocity](sandbox)
	nameHandler := ComponentLinker[name](sandbox)

	for range numPosition {
		id := LinkEntity(sandbox)
		positionHandler.Link(id)
	}
	for range numPositionVelocity {
		id := LinkEntity(sandbox)
		positionHandler.Link(id)
		velocityHandler.Link(id)
	}

	Filter(sandbox, filter.Match2[position, velocity]())
	Filter(sandbox, filter.Match2[position, armor]())
	Filter(sandbox, filter.Match2[position, rendered]())
	Update(sandbox)
	b.StartTimer()

	for range b.N / 2 {
		for entityId := range numEdits {
			UnlinkEntity(sandbox, entity.Id(entityId))
		}
		Update(sandbox)
		for entityId := range numEdits {
			LinkEntity(sandbox)
			positionHandler.Link(entity.Id(entityId))
			nameHandler.Link(entity.Id(entityId))
		}
		Update(sandbox)
	}
}
