package benchmarks

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"testing"
)

func Benchmark_Iter_Donburi(b *testing.B) {
	b.StopTimer()
	world := donburi.NewWorld()

	var p = donburi.NewComponentType[position]()
	var v = donburi.NewComponentType[velocity]()

	for i := 0; i < numPosition; i++ {
		world.Create(p)
	}
	for i := 0; i < numPositionVelocity; i++ {
		world.Create(p, v)
	}

	query := donburi.NewQuery(filter.Contains(p, v))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		query.Each(world, func(entry *donburi.Entry) {
			pos := p.Get(entry)
			vel := v.Get(entry)

			pos.X += vel.X
			pos.Y += vel.Y
		})
	}
}
