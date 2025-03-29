package benchmarks

import (
	"github.com/wfranczyk/ento"
	"testing"
)

type PosVelSystem struct {
	Pos *position `ento:"required"`
	Vel *velocity `ento:"required"`
}

func (s *PosVelSystem) Update(entity *ento.Entity) {
	s.Pos.X += s.Vel.X
	s.Pos.Y += s.Vel.Y
}

func Benchmark_Iter_Ento(b *testing.B) {
	b.StopTimer()
	world := ento.NewWorldBuilder().
		WithDenseComponents(position{}).
		WithSparseComponents(velocity{}).
		Build(1024)

	system := PosVelSystem{}
	world.AddSystems(&system)

	for i := 0; i < numPosition; i++ {
		world.AddEntity(position{})
	}
	for i := 0; i < numPositionVelocity; i++ {
		world.AddEntity(position{}, velocity{})
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		world.Update()
	}
}
