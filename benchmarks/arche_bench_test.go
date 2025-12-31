package benchmarks

import (
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func Benchmark_Iter_ArcheGeneric(b *testing.B) {
	b.StopTimer()
	world := ecs.NewWorld()

	posID := ecs.ComponentID[position](&world)
	positionVelocityMapper := generic.NewMap2[velocity, position](&world)
	for range numPosition {
		ecs.NewBuilder(&world, posID).New()
	}
	for range numPositionVelocity {
		positionVelocityMapper.New()
	}

	filter := generic.NewFilter2[velocity, position]()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query := filter.Query(&world)
		for query.Next() {
			v, p := query.Get()
			p.X += v.X
			p.Y += v.Y
		}
	}
}
