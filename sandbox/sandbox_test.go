package sandbox

import (
	"fmt"
	"github.com/andrei-cosmin/hakkt/filter"
	"github.com/stretchr/testify/assert"
	"testing"
)

const nPos = 9000
const nPosVel = 1000

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

type B struct {
	B string
}

type C struct {
	C float64
}

type D struct {
	D bool
}

type A struct {
	A int
}

func Test(t *testing.T) {

	var sandbox = New()
	a := ComponentLinker[A](sandbox)
	b := ComponentLinker[B](sandbox)
	c := ComponentLinker[C](sandbox)

	for index := 0; index < 100; index++ {
		var entity = LinkEntity(sandbox)
		a.Link(entity)
	}
	for index := 0; index < 100; index++ {
		var entity = LinkEntity(sandbox)
		a.Link(entity)
		b.Link(entity)
	}
	for index := 0; index < 100; index++ {
		var entity = LinkEntity(sandbox)
		a.Link(entity)
		b.Link(entity)
		c.Link(entity)
	}

	var responseA = Filter(sandbox, filter.Match[A]())
	var responseB = Filter(sandbox, filter.Match[B]())
	var responseC = Filter(sandbox, filter.Match[C]())

	var responseAB = Filter(sandbox, filter.Match2[A, B]())
	var responseAC = Filter(sandbox, filter.Match2[A, C]())
	var responseBC = Filter(sandbox, filter.Match2[B, C]())
	var responseABC = Filter(sandbox, filter.Match3[A, B, C]())

	Update(sandbox)

	assert.Equal(t, 300, len(responseA.EntityIds()))
	assert.Equal(t, 200, len(responseB.EntityIds()))
	assert.Equal(t, 100, len(responseC.EntityIds()))
	assert.Equal(t, 200, len(responseAB.EntityIds()))
	assert.Equal(t, 100, len(responseAC.EntityIds()))
	assert.Equal(t, 100, len(responseBC.EntityIds()))
	assert.Equal(t, 100, len(responseABC.EntityIds()))

	UnlinkEntity(sandbox, 0)
	Update(sandbox)

	assert.Equal(t, 299, len(responseA.EntityIds()))
	assert.Equal(t, 200, len(responseB.EntityIds()))
	assert.Equal(t, 100, len(responseC.EntityIds()))
	assert.Equal(t, 200, len(responseAB.EntityIds()))
	assert.Equal(t, 100, len(responseAC.EntityIds()))
	assert.Equal(t, 100, len(responseBC.EntityIds()))
	assert.Equal(t, 100, len(responseABC.EntityIds()))

	UnlinkEntity(sandbox, 299)
	Update(sandbox)

	assert.Equal(t, 298, len(responseA.EntityIds()))
	assert.Equal(t, 199, len(responseB.EntityIds()))
	assert.Equal(t, 99, len(responseC.EntityIds()))
	assert.Equal(t, 199, len(responseAB.EntityIds()))
	assert.Equal(t, 99, len(responseAC.EntityIds()))
	assert.Equal(t, 99, len(responseBC.EntityIds()))
	assert.Equal(t, 99, len(responseABC.EntityIds()))

	cmpA := ComponentLinker[A](sandbox)
	cmpA.Unlink(1)
	Update(sandbox)

	assert.Equal(t, 297, len(responseA.EntityIds()))
	assert.Equal(t, 199, len(responseB.EntityIds()))
	assert.Equal(t, 99, len(responseC.EntityIds()))
	assert.Equal(t, 199, len(responseAB.EntityIds()))
	assert.Equal(t, 99, len(responseAC.EntityIds()))
	assert.Equal(t, 99, len(responseBC.EntityIds()))
	assert.Equal(t, 99, len(responseABC.EntityIds()))
}

func Test3(t *testing.T) {
	var sandbox = New()
	p := ComponentLinker[Position](sandbox)
	v := ComponentLinker[Velocity](sandbox)
	aa := ComponentLinker[A](sandbox)
	bb := ComponentLinker[B](sandbox)

	for i := 0; i < nPos; i++ {
		id := LinkEntity(sandbox)
		p.Link(id)
	}
	for i := 0; i < nPosVel; i++ {
		id := LinkEntity(sandbox)
		p.Link(id)
		v.Link(id)
	}
	var _ = Filter(sandbox, filter.Match[Position](), filter.Match[Velocity]())

	aa.Link(0)
	UnlinkEntity(sandbox, 0)
	fmt.Printf("Unliked entity 0 with component A\n")

	Update(sandbox)
	added := LinkEntity(sandbox)
	aa.Link(added)
	bb.Link(added)
	bb.Get(added).B = "Hello"
	aa.Get(added).A = 42

	fmt.Printf("Linked entity %d with component B\n", added)
	Update(sandbox)

	fmt.Printf("Entity %d has component B %t\n", added, bb.Has(added))
	fmt.Printf("Next available entity id: %d\n", LinkEntity(sandbox))

	println("ss", bb.Get(added).B)
	println("ss", aa.Get(added).A)
	UnlinkEntity(sandbox, added)
	Update(sandbox)
	//added1 := LinkEntity(sandbox)
	//
	//aa.Link(added)
	//println(aa.Get(added).A)

}
