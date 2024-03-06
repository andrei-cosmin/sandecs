package tests

const (
	numEntities = 500000
	numRemoves  = 50000

	entityNotLinkedMsg   = "Entity not linked %d"
	entityLinkedOrderMsg = "Entity linked out of order %d"

	entityUnlinkedBeforeUpdateMsg = "Entity unlinked before next update %d"
	entityNotUnlinkedMsg          = "Entity not unlinked after update %d"

	entityNotRecycledMsg  = "Entity id not recycled %d"
	entityRecycleOrderMsg = "Entity id was recycled out of order %d"

	componentNotLinkedMsg   = "Component %s not linked for entity %d"
	componentValueMsg       = "Component %s value incorrect for entity %d"
	componentNotUnlinkedMsg = "Component %s not unlinked after update for entity %d"

	filterIncorrectNumEntitiesMsg = "Filter returned incorrect number of entities"

	positionComponent = "POSITION"
	velocityComponent = "VELOCITY"
	healthComponent   = "HEALTH"
	armorComponent    = "ARMOR"
	renderedComponent = "RENDERED"
	nameComponent     = "NAME"
)

type position struct {
	X float64
	Y float64
}

type velocity struct {
	X float64
	Y float64
}

type name struct {
	value string
}

type health struct {
	value float64
}

type rendered struct {
	value bool
}

type armor struct {
	value int
}
