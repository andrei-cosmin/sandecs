package benchmarks

import (
	"testing"

	sand "github.com/andrei-cosmin/sandecs"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/filter"
	"github.com/andrei-cosmin/sandecs/options"
)

const benchmarkSandboxMode = options.Standard

func Benchmark_Iter_Sand(b *testing.B) {
	b.StopTimer()
	box := sand.New(benchmarkSandboxMode, numPosition+numPositionVelocity+10000, 8, 0)
	positionHandler := sand.ComponentLinker[position](box)
	velocityHandler := sand.ComponentLinker[velocity](box)

	for range numPosition {
		id := sand.LinkEntity(box)
		positionHandler.Link(id)
	}
	for range numPositionVelocity {
		id := sand.LinkEntity(box)
		positionHandler.Link(id)
		velocityHandler.Link(id)
	}
	b.StartTimer()

	view := sand.Filter(box, filter.Match2[position, velocity]())
	sand.Update(box)

	for range b.N {
		for _, entityId := range view.EntityIds() {
			p := positionHandler.Get(entityId)
			v := velocityHandler.Get(entityId)
			p.X += v.X
			p.Y += v.Y
		}
		sand.Update(box)
	}
}

func Benchmark_LinkComponents_Sand(b *testing.B) {
	b.StopTimer()
	box := sand.New(benchmarkSandboxMode, numPosition+numPositionVelocity, 4, 20000)
	positionHandler := sand.ComponentLinker[position](box)
	velocityHandler := sand.ComponentLinker[velocity](box)
	nameHandler := sand.TagLinker(box, "name")

	for range numPosition {
		id := sand.LinkEntity(box)
		positionHandler.Link(id)
	}
	for range numPositionVelocity {
		id := sand.LinkEntity(box)
		positionHandler.Link(id)
		velocityHandler.Link(id)
	}

	sand.Filter(box, filter.Match2[position, velocity]())
	sand.Update(box)
	b.StartTimer()

	for range b.N / 2 {
		for entityId := range numEdits {
			nameHandler.Link(entity.Id(entityId))
		}
		sand.Update(box)
		for entityId := range numEdits {
			nameHandler.Unlink(entity.Id(entityId))
		}
		sand.Update(box)
	}
}

func Benchmark_LinkEntities_Sand(b *testing.B) {
	b.StopTimer()
	box := sand.New(benchmarkSandboxMode, numPosition+numPositionVelocity, 4, 20000)
	positionHandler := sand.ComponentLinker[position](box)
	velocityHandler := sand.ComponentLinker[velocity](box)
	nameHandler := sand.TagLinker(box, "name")

	for range numPosition {
		id := sand.LinkEntity(box)
		positionHandler.Link(id)
	}
	for range numPositionVelocity {
		id := sand.LinkEntity(box)
		positionHandler.Link(id)
		velocityHandler.Link(id)
	}

	sand.Filter(box, filter.Match2[position, velocity]())
	sand.Update(box)
	b.StartTimer()

	for range b.N / 2 {
		for entityId := range numEdits {
			sand.UnlinkEntity(box, entity.Id(entityId))
		}
		sand.Update(box)
		for entityId := range numEdits {
			sand.LinkEntity(box)
			positionHandler.Link(entity.Id(entityId))
			nameHandler.Link(entity.Id(entityId))
		}
		sand.Update(box)
	}
}

func Benchmark_UpdateComponents_Sand(b *testing.B) {
	b.StopTimer()
	box := sand.New(benchmarkSandboxMode, numPosition+numPositionVelocity, 4, 20000)
	positionHandler := sand.ComponentLinker[position](box)
	velocityHandler := sand.ComponentLinker[velocity](box)

	for range numPosition {
		id := sand.LinkEntity(box)
		positionHandler.Link(id)
	}
	for range numPositionVelocity {
		id := sand.LinkEntity(box)
		positionHandler.Link(id)
		velocityHandler.Link(id)
	}

	sand.Filter(box, filter.Match2[position, velocity]())
	sand.Update(box)
	b.StartTimer()

	for range b.N {
		for entityId := range numEdits {
			instance := positionHandler.Get(entity.Id(entityId))
			instance.X += 1
		}
		sand.Update(box)
	}
}
