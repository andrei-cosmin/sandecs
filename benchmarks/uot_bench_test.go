package benchmarks

import (
	"github.com/unitoftime/ecs"
	"testing"
)

func Benchmark_Iter_Uot(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld()

	for i := 0; i < numPosition; i++ {
		id := world.NewId()
		ecs.Write(world, id,
			ecs.C(position{0, 0}),
		)
	}
	for i := 0; i < numPositionVelocity; i++ {
		id := world.NewId()
		ecs.Write(world, id,
			ecs.C(position{0, 0}),
			ecs.C(velocity{0, 0}),
		)
	}
	query := ecs.Query2[position, velocity](world)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query.MapId(func(id ecs.Id, pos *position, vel *velocity) {
			pos.X += vel.X
			pos.Y += vel.Y
		})
	}
}
