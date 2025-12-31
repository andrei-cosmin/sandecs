package benchmarks

import (
	"testing"

	ecs "github.com/marioolofo/go-gameengine-ecs"
)

const (
	positionComponentID ecs.ComponentID = iota
	velocityComponentID
)

func Benchmark_Iter_GGEcs(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld(1024)
	world.Register(ecs.NewComponentRegistry[position](positionComponentID))
	world.Register(ecs.NewComponentRegistry[velocity](velocityComponentID))

	for i := 0; i < numPosition; i++ {
		_ = world.NewEntity(positionComponentID)
	}
	for i := 0; i < numPositionVelocity; i++ {
		_ = world.NewEntity(positionComponentID, velocityComponentID)
	}

	mask := ecs.MakeComponentMask(positionComponentID, velocityComponentID)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := world.Query(mask)
		for query.Next() {
			pos := (*position)(query.Component(positionComponentID))
			vel := (*velocity)(query.Component(velocityComponentID))
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
