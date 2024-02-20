package sandbox

import (
	"github.com/andrei-cosmin/hakkt/marker"
	"github.com/andrei-cosmin/hakkt/query"
	"github.com/stretchr/testify/assert"
	"testing"
)

const nPos = 9000
const nPosVel = 1000

type HealthComponent struct {
	Health int
}

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

type A struct {
}

type B struct {
}

type C struct {
}

type D struct {
}

func Test(t *testing.T) {

	var sandbox = New()
	a := GetComponentLinker[A](sandbox)
	b := GetComponentLinker[B](sandbox)
	c := GetComponentLinker[C](sandbox)

	for index := 0; index < 100; index++ {
		var entity = LinkEntity(sandbox)
		a.Link(entity, A{})
	}
	for index := 0; index < 100; index++ {
		var entity = LinkEntity(sandbox)
		a.Link(entity, A{})
		b.Link(entity, B{})
	}
	for index := 0; index < 100; index++ {
		var entity = LinkEntity(sandbox)
		a.Link(entity, A{})
		b.Link(entity, B{})
		c.Link(entity, C{})
	}

	var responseA = Filter(sandbox, query.New().MatchAll(marker.Symbol[A]()))
	var responseB = Filter(sandbox, query.New().MatchAll(marker.Symbol[B]()))
	var responseC = Filter(sandbox, query.New().MatchAll(marker.Symbol[C]()))

	var responseAB = Filter(sandbox, query.New().MatchAll(
		marker.Symbol[A](),
		marker.Symbol[B](),
	))
	var responseAC = Filter(sandbox, query.New().MatchAll(
		marker.Symbol[A](),
		marker.Symbol[C](),
	))
	var responseBC = Filter(sandbox, query.New().MatchAll(
		marker.Symbol[B](),
		marker.Symbol[C](),
	))
	var responseABC = Filter(sandbox, query.New().MatchAll(
		marker.Symbol[A](),
		marker.Symbol[B](),
		marker.Symbol[C](),
	))

	Update(sandbox)

	assert.Equal(t, 300, len(responseA.GetEntities()))
	assert.Equal(t, 200, len(responseB.GetEntities()))
	assert.Equal(t, 100, len(responseC.GetEntities()))
	assert.Equal(t, 200, len(responseAB.GetEntities()))
	assert.Equal(t, 100, len(responseAC.GetEntities()))
	assert.Equal(t, 100, len(responseBC.GetEntities()))
	assert.Equal(t, 100, len(responseABC.GetEntities()))

	UnlinkEntity(sandbox, 0)
	Update(sandbox)

	assert.Equal(t, 299, len(responseA.GetEntities()))
	assert.Equal(t, 200, len(responseB.GetEntities()))
	assert.Equal(t, 100, len(responseC.GetEntities()))
	assert.Equal(t, 200, len(responseAB.GetEntities()))
	assert.Equal(t, 100, len(responseAC.GetEntities()))
	assert.Equal(t, 100, len(responseBC.GetEntities()))
	assert.Equal(t, 100, len(responseABC.GetEntities()))

	UnlinkEntity(sandbox, 299)
	Update(sandbox)

	assert.Equal(t, 298, len(responseA.GetEntities()))
	assert.Equal(t, 199, len(responseB.GetEntities()))
	assert.Equal(t, 99, len(responseC.GetEntities()))
	assert.Equal(t, 199, len(responseAB.GetEntities()))
	assert.Equal(t, 99, len(responseAC.GetEntities()))
	assert.Equal(t, 99, len(responseBC.GetEntities()))
	assert.Equal(t, 99, len(responseABC.GetEntities()))

	cmpA := GetComponentLinker[A](sandbox)
	cmpA.Remove(1)
	Update(sandbox)

	assert.Equal(t, 297, len(responseA.GetEntities()))
	assert.Equal(t, 199, len(responseB.GetEntities()))
	assert.Equal(t, 99, len(responseC.GetEntities()))
	assert.Equal(t, 199, len(responseAB.GetEntities()))
	assert.Equal(t, 99, len(responseAC.GetEntities()))
	assert.Equal(t, 99, len(responseBC.GetEntities()))
	assert.Equal(t, 99, len(responseABC.GetEntities()))
}

func Test3(b *testing.T) {
	var sandbox = New()
	p := GetComponentLinker[Position](sandbox)
	v := GetComponentLinker[Velocity](sandbox)

	pos := Position{X: 0, Y: 0}
	vel := Velocity{X: 0, Y: 0}

	for i := 0; i < nPos; i++ {
		id := LinkEntity(sandbox)
		p.Link(id, pos)
	}
	for i := 0; i < nPosVel; i++ {
		id := LinkEntity(sandbox)
		p.Link(id, pos)
		v.Link(id, vel)
	}

	var _ = Filter(sandbox, query.New().MatchAll(
		marker.Symbol[Position](),
		marker.Symbol[Velocity](),
	))
	Update(sandbox)
}
