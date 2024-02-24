package sandbox

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestSandboxSuite(t *testing.T) {
	suite.Run(t, &SandboxTestSuite{pooling: false})
	suite.Run(t, &SandboxTestSuite{pooling: true})
}

type SandboxTestSuite struct {
	sandboxSuite
	pooling        bool
	positionLinker component.Linker[position]
	velocityLinker component.Linker[velocity]
	renderedLinker component.Linker[rendered]
	armorLinker    component.Linker[armor]
	healthLinker   component.Linker[health]
	nameLinker     component.Linker[name]
}

func (suite *SandboxTestSuite) SetupTest() {
	suite.T().Log("Pooling:", suite.pooling)
	if suite.pooling {
		suite.sandbox = NewDefault()
	} else {
		suite.sandbox = New(DefaultNumEntities, DefaultNumComponents, 0)
	}
	suite.positionLinker = ComponentLinker[position](suite.sandbox)
	suite.velocityLinker = ComponentLinker[velocity](suite.sandbox)
	suite.renderedLinker = ComponentLinker[rendered](suite.sandbox)
	suite.armorLinker = ComponentLinker[armor](suite.sandbox)
	suite.healthLinker = ComponentLinker[health](suite.sandbox)
	suite.nameLinker = ComponentLinker[name](suite.sandbox)
}

func (suite *SandboxTestSuite) TestLinkingEntity() {
	for i := 0; i < numEntities; i++ {
		LinkEntity(suite.sandbox)
	}
	Update(suite.sandbox)
	for i := 0; i < numEntities; i++ {
		suite.assertEntity(entity.Id(i), entityNotLinkedMsg, i)
	}
}

func (suite *SandboxTestSuite) TestUnlinkingEntity() {
	for i := 0; i < numEntities; i++ {
		LinkEntity(suite.sandbox)
	}
	Update(suite.sandbox)
	for i := 0; i < numEntities/2; i++ {
		UnlinkEntity(suite.sandbox, entity.Id(i))
	}
	Update(suite.sandbox)
	for i := 0; i < numEntities/2; i++ {
		suite.assertDeletedEntity(entity.Id(i), entityNotUnlinkedMsg, i)
	}
	for i := numEntities / 2; i < numEntities; i++ {
		suite.assertEntity(entity.Id(i), entityNotLinkedMsg, i)
	}
}

func (suite *SandboxTestSuite) TestDuplicateUnlinkingEntity() {
	for i := 0; i < numEntities; i++ {
		LinkEntity(suite.sandbox)
	}
	Update(suite.sandbox)
	for i := 0; i < numEntities/2; i++ {
		UnlinkEntity(suite.sandbox, entity.Id(i))
	}
	for i := 0; i < numEntities/2; i++ {
		UnlinkEntity(suite.sandbox, entity.Id(i))
	}
	for i := 0; i < numEntities/2; i++ {
		suite.assertEntity(entity.Id(i), entityNotLinkedMsg, i)
	}
	Update(suite.sandbox)
	for i := 0; i < numEntities/2; i++ {
		UnlinkEntity(suite.sandbox, entity.Id(i))
	}
	for i := 0; i < numEntities/2; i++ {
		suite.assertDeletedEntity(entity.Id(i), entityNotUnlinkedMsg, i)
	}
}

func (suite *SandboxTestSuite) TestEntityRecycling() {
	for i := 0; i < numEntities; i++ {
		entityId := LinkEntity(suite.sandbox)
		assert.Equal(suite.T(), entity.Id(i), entityId, entityLinkedOrderMsg, i)
		suite.assertEntity(entity.Id(i), entityNotLinkedMsg, i)
		UnlinkEntity(suite.sandbox, entity.Id(i))
		suite.assertEntity(entity.Id(i), entityUnlinkedBeforeUpdateMsg, i)
	}

	Update(suite.sandbox)

	for i := 0; i < numEntities; i++ {
		suite.assertDeletedEntity(entity.Id(i), entityNotUnlinkedMsg, i)
	}

	for i := 0; i < numEntities; i++ {
		entityId := LinkEntity(suite.sandbox)
		assert.Equal(suite.T(), entity.Id(i), entityId, entityNotRecycledMsg, i)
		suite.assertEntity(entity.Id(i), entityNotLinkedMsg, i)
	}
}

func (suite *SandboxTestSuite) TestEntityRecyclingRandom() {
	randomEntityIds := getRandomIds(numEntities, numRemoves)

	for i := 0; i < numEntities; i++ {
		LinkEntity(suite.sandbox)
	}
	for _, entityId := range randomEntityIds {
		UnlinkEntity(suite.sandbox, entityId)
	}
	Update(suite.sandbox)
	slices.Sort(randomEntityIds)
	for _, entityId := range randomEntityIds {
		suite.assertDeletedEntity(entityId, entityNotUnlinkedMsg, entityId)
		newId := LinkEntity(suite.sandbox)
		assert.Equal(suite.T(), entityId, newId, entityRecycleOrderMsg, newId)
	}
}

func (suite *SandboxTestSuite) TestSimpleLinkingComponents() {
	for i := 0; i < numEntities; i++ {
		LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(i))
		suite.positionLinker.Get(entity.Id(i)).X = float64(i)
		suite.assertComponent(suite.positionLinker, entity.Id(i), componentValueMsg, positionComponent, i)
	}

	Update(suite.sandbox)

	for i := 0; i < numEntities; i++ {
		suite.assertComponent(suite.positionLinker, entity.Id(i), componentNotLinkedMsg, positionComponent, i)
		assert.Equal(suite.T(), float64(i), suite.positionLinker.Get(entity.Id(i)).X, componentValueMsg, positionComponent, i)
	}
}

func (suite *SandboxTestSuite) TestDuplicateLinkingComponent() {
	entityId := LinkEntity(suite.sandbox)
	xValue := 100.5

	suite.positionLinker.Link(entityId)
	suite.assertComponent(suite.positionLinker, entityId, componentNotLinkedMsg, positionComponent, entityId)

	suite.positionLinker.Get(entityId).X = xValue
	suite.positionLinker.Link(entityId)
	assert.Equal(suite.T(), xValue, suite.positionLinker.Get(entityId).X, componentValueMsg, positionComponent, entityId)
}

func (suite *SandboxTestSuite) TestDuplicateUnlinkingComponent() {
	entityId := LinkEntity(suite.sandbox)
	suite.positionLinker.Link(entityId)
	suite.assertComponent(suite.positionLinker, entityId, componentNotLinkedMsg, positionComponent, entityId)

	suite.positionLinker.Unlink(entityId)
	suite.positionLinker.Unlink(entityId)
	suite.assertComponent(suite.positionLinker, entityId, componentNotLinkedMsg, positionComponent, entityId)
	Update(suite.sandbox)

	suite.assertDeletedComponent(suite.positionLinker, entityId, componentNotUnlinkedMsg, positionComponent, entityId)
	suite.positionLinker.Unlink(entityId)
	suite.positionLinker.Unlink(entityId)
	Update(suite.sandbox)
	suite.assertDeletedComponent(suite.positionLinker, entityId, componentNotUnlinkedMsg, positionComponent, entityId)
}

func (suite *SandboxTestSuite) TestComponentForInvalidEntity() {
	suite.positionLinker.Link(10 * numEntities)
	suite.assertDeletedComponent(suite.positionLinker, 10*numEntities, componentNotUnlinkedMsg, positionComponent, 10*numEntities)
	suite.positionLinker.Unlink(10 * numEntities)
	Update(suite.sandbox)
	suite.assertDeletedComponent(suite.positionLinker, 10*numEntities, componentNotLinkedMsg, positionComponent, 10*numEntities)
}

func (suite *SandboxTestSuite) TestHeavyLinkingComponents() {
	for i := 0; i < numEntities; i++ {
		LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(i))
		suite.positionLinker.Get(entity.Id(i)).X = float64(i)
		suite.assertComponent(suite.positionLinker, entity.Id(i), componentNotLinkedMsg, positionComponent, i)
	}
	for i := 0; i < numEntities/2; i++ {
		suite.velocityLinker.Link(entity.Id(i))
		suite.velocityLinker.Get(entity.Id(i)).Y = float64(i * 2)
		suite.assertComponent(suite.velocityLinker, entity.Id(i), componentNotLinkedMsg, velocityComponent, i)
	}
	for i := numEntities / 2; i < numEntities; i++ {
		suite.healthLinker.Link(entity.Id(i))
		suite.healthLinker.Get(entity.Id(i)).value = float64(i / 2)
		suite.assertComponent(suite.healthLinker, entity.Id(i), componentNotLinkedMsg, healthComponent, i)
	}
	for i := numEntities / 4; i < numEntities*3/4; i++ {
		suite.armorLinker.Link(entity.Id(i))
		suite.armorLinker.Get(entity.Id(i)).value = i / 10
		suite.assertComponent(suite.armorLinker, entity.Id(i), componentNotLinkedMsg, armorComponent, i, i)
	}

	Update(suite.sandbox)

	for i := 0; i < numEntities; i++ {
		suite.assertComponent(suite.positionLinker, entity.Id(i), componentNotLinkedMsg, positionComponent, i)
		assert.Equal(suite.T(), float64(i), suite.positionLinker.Get(entity.Id(i)).X, componentValueMsg, positionComponent, i)
	}
	for i := 0; i < numEntities/2; i++ {
		suite.assertComponent(suite.velocityLinker, entity.Id(i), componentNotLinkedMsg, velocityComponent, i)
		assert.Equal(suite.T(), float64(i*2), suite.velocityLinker.Get(entity.Id(i)).Y, componentValueMsg, velocityComponent, i)
	}
	for i := numEntities / 2; i < numEntities; i++ {
		suite.assertComponent(suite.healthLinker, entity.Id(i), componentNotLinkedMsg, healthComponent, i)
		assert.Equal(suite.T(), float64(i/2), suite.healthLinker.Get(entity.Id(i)).value, componentValueMsg, healthComponent, i)

	}
	for i := numEntities / 4; i < numEntities*3/4; i++ {
		suite.assertComponent(suite.armorLinker, entity.Id(i), componentNotLinkedMsg, armorComponent, i, i)
		assert.Equal(suite.T(), i/10, suite.armorLinker.Get(entity.Id(i)).value, componentValueMsg, armorComponent, i)
	}
}

func (suite *SandboxTestSuite) TestSimpleUnlinkingComponents() {
	removedEntityId := entity.Id(rand.IntN(numEntities))
	removedComponentEntityId := entity.Id(rand.IntN(numEntities))

	for i := 0; i < numEntities; i++ {
		LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(i))
	}

	UnlinkEntity(suite.sandbox, removedEntityId)
	suite.positionLinker.Unlink(removedComponentEntityId)

	Update(suite.sandbox)

	suite.assertDeletedEntity(removedEntityId, entityNotUnlinkedMsg, removedEntityId)
	suite.assertDeletedComponent(suite.positionLinker, removedEntityId, componentNotUnlinkedMsg, positionComponent, removedEntityId)
	suite.assertDeletedComponent(suite.positionLinker, removedComponentEntityId, componentNotUnlinkedMsg, positionComponent, removedComponentEntityId)

	readdedEntityId := LinkEntity(suite.sandbox)
	assert.Equal(suite.T(), removedEntityId, readdedEntityId, entityRecycleOrderMsg, readdedEntityId)
	suite.assertDeletedComponent(suite.positionLinker, readdedEntityId, componentNotUnlinkedMsg, positionComponent, readdedEntityId)

	for i := 0; i < numEntities; i++ {
		if i != int(removedEntityId) && i != int(removedComponentEntityId) {
			suite.assertComponent(suite.positionLinker, entity.Id(i), componentNotLinkedMsg, positionComponent, i)
		}
	}
}

func (suite *SandboxTestSuite) TestHeavyUnlinkingComponents() {
	for i := 0; i < numEntities; i++ {
		LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(i))
	}
	for i := 0; i < numEntities/2; i++ {
		suite.velocityLinker.Link(entity.Id(i))
	}
	for i := numEntities / 2; i < numEntities; i++ {
		suite.healthLinker.Link(entity.Id(i))
	}
	for i := numEntities / 4; i < numEntities*3/4; i++ {
		suite.armorLinker.Link(entity.Id(i))
	}

	removedEntities := getRandomIds(numEntities, numRemoves)
	removedPosition := getRandomIds(numEntities, numRemoves)
	removedVelocity := getRandomIds(numEntities/2, numRemoves)
	removedHealth := getRandomIds(numEntities/4, numRemoves)
	removedArmor := getRandomIds(numEntities/3, numRemoves)

	for _, entityId := range removedEntities {
		UnlinkEntity(suite.sandbox, entityId)
	}
	for _, entityId := range removedPosition {
		suite.positionLinker.Unlink(entityId)
	}
	for _, entityId := range removedVelocity {
		suite.velocityLinker.Unlink(entityId)
	}
	for _, entityId := range removedHealth {
		suite.healthLinker.Unlink(entityId)
	}
	for _, entityId := range removedArmor {
		suite.armorLinker.Unlink(entityId)
	}

	Update(suite.sandbox)

	for _, entityId := range removedEntities {
		suite.assertDeletedEntity(entityId, entityNotUnlinkedMsg, entityId)
		suite.assertDeletedComponent(suite.positionLinker, entityId, componentNotUnlinkedMsg, positionComponent, entityId)
		suite.assertDeletedComponent(suite.velocityLinker, entityId, componentNotUnlinkedMsg, velocityComponent, entityId)
		suite.assertDeletedComponent(suite.healthLinker, entityId, componentNotUnlinkedMsg, healthComponent, entityId)
		suite.assertDeletedComponent(suite.armorLinker, entityId, componentNotUnlinkedMsg, armorComponent, entityId)
	}
	for _, entityId := range removedPosition {
		suite.assertDeletedComponent(suite.positionLinker, entityId, componentNotUnlinkedMsg, positionComponent, entityId)
	}
	for _, entityId := range removedVelocity {
		suite.assertDeletedComponent(suite.velocityLinker, entityId, componentNotUnlinkedMsg, velocityComponent, entityId)
	}
	for _, entityId := range removedHealth {
		suite.assertDeletedComponent(suite.healthLinker, entityId, componentNotUnlinkedMsg, healthComponent, entityId)
	}
	for _, entityId := range removedArmor {
		suite.assertDeletedComponent(suite.armorLinker, entityId, componentNotUnlinkedMsg, armorComponent, entityId)
	}
}

func (suite *SandboxTestSuite) TestSimplePooling() {
	for i := 0; i < numEntities; i++ {
		LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(i))
		suite.positionLinker.Get(entity.Id(i)).X = float64(i + numEntities)
	}
	Update(suite.sandbox)
	for i := 0; i < numEntities; i++ {
		suite.positionLinker.Unlink(entity.Id(i))
	}
	Update(suite.sandbox)
	for i := 0; i < numEntities; i++ {
		suite.positionLinker.Link(entity.Id(i))
	}
	Update(suite.sandbox)

	poolSize := 0
	if suite.pooling {
		poolSize = DefaultPoolCapacity
	}
	for i := 0; i < poolSize; i++ {
		assert.NotEqual(suite.T(), float64(0), suite.positionLinker.Get(entity.Id(i)).X, componentValueMsg, positionComponent, i)
	}
	for i := poolSize; i < numEntities; i++ {
		assert.Equal(suite.T(), float64(0), suite.positionLinker.Get(entity.Id(i)).X, componentValueMsg, positionComponent, i)
	}

}

func (suite *SandboxTestSuite) TestSimpleFilter() {
	positionFilter := Filter(suite.sandbox, filter.Match[position]())
	duplicateFilter := Filter(suite.sandbox, filter.Match[position]())
	assert.Equal(suite.T(), positionFilter, duplicateFilter)
	Filter(suite.sandbox, filter.Match[position]())

	assert.Len(suite.T(), positionFilter.EntityIds(), 0, filterIncorrectNumEntitiesMsg)
	for i := 0; i < numEntities; i++ {
		LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(i))
	}
	assert.Len(suite.T(), positionFilter.EntityIds(), 0, filterIncorrectNumEntitiesMsg)
	Update(suite.sandbox)
	assert.Len(suite.T(), positionFilter.EntityIds(), numEntities, filterIncorrectNumEntitiesMsg)
}

func (suite *SandboxTestSuite) TestHeavyFilter() {
	filter1 := Filter(suite.sandbox, filter.Match2[position, velocity](), filter.Exclude[armor]())
	filter2 := Filter(suite.sandbox, filter.Match3[position, velocity, armor]())
	filter3 := Filter(suite.sandbox, filter.Union2[velocity, armor]())
	filter4 := Filter(suite.sandbox, filter.Union3[velocity, armor, position]())

	filterCount1 := 0
	filterCount2 := 0
	filterCount3 := 0
	filterCount4 := 0

	for i := 0; i < numEntities; i++ {
		LinkEntity(suite.sandbox)
	}
	for i := 0; i < numEntities/5; i++ {
		suite.positionLinker.Link(entity.Id(i))
	}
	for i := numEntities / 4; i < numEntities*3/4; i++ {
		suite.velocityLinker.Link(entity.Id(i))
	}
	for i := numEntities / 2; i < numEntities; i++ {
		suite.armorLinker.Link(entity.Id(i))
	}

	assert.Len(suite.T(), filter1.EntityIds(), filterCount1, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter2.EntityIds(), filterCount2, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter3.EntityIds(), filterCount3, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter4.EntityIds(), filterCount4, filterIncorrectNumEntitiesMsg)

	Update(suite.sandbox)

	filterCount3 += numEntities - numEntities/4
	filterCount4 += numEntities - numEntities/4 + numEntities/5

	assert.Len(suite.T(), filter1.EntityIds(), filterCount1, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter2.EntityIds(), filterCount2, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter3.EntityIds(), filterCount3, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter4.EntityIds(), filterCount4, filterIncorrectNumEntitiesMsg)

	for i := 0; i < numEntities/10; i++ {
		suite.velocityLinker.Link(entity.Id(i))
	}
	Update(suite.sandbox)

	filterCount1 += numEntities / 10
	filterCount3 += numEntities / 10

	assert.Len(suite.T(), filter1.EntityIds(), filterCount1, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter2.EntityIds(), filterCount2, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter3.EntityIds(), filterCount3, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter4.EntityIds(), filterCount4, filterIncorrectNumEntitiesMsg)

	for i := 0; i < numEntities/14; i++ {
		suite.armorLinker.Link(entity.Id(i))
	}
	Update(suite.sandbox)

	filterCount1 -= numEntities / 14
	filterCount2 += numEntities / 14

	assert.Len(suite.T(), filter1.EntityIds(), filterCount1, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter2.EntityIds(), filterCount2, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter3.EntityIds(), filterCount3, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter4.EntityIds(), filterCount4, filterIncorrectNumEntitiesMsg)
}
