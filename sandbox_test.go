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
	for range numEntities {
		LinkEntity(suite.sandbox)
	}
	Update(suite.sandbox)
	for index := range numEntities {
		suite.assertEntity(entity.Id(index), entityNotLinkedMsg, index)
	}
}

func (suite *SandboxTestSuite) TestUnlinkingEntity() {
	for range numEntities {
		LinkEntity(suite.sandbox)
	}
	Update(suite.sandbox)
	for index := range numEntities / 2 {
		UnlinkEntity(suite.sandbox, entity.Id(index))
	}
	Update(suite.sandbox)
	for index := range numEntities / 2 {
		suite.assertDeletedEntity(entity.Id(index), entityNotUnlinkedMsg, index)
	}
	for index := numEntities / 2; index < numEntities; index++ {
		suite.assertEntity(entity.Id(index), entityNotLinkedMsg, index)
	}
}

func (suite *SandboxTestSuite) TestDuplicateUnlinkingEntity() {
	for range numEntities {
		LinkEntity(suite.sandbox)
	}
	Update(suite.sandbox)
	for index := range numEntities / 2 {
		UnlinkEntity(suite.sandbox, entity.Id(index))
	}
	for index := range numEntities / 2 {
		UnlinkEntity(suite.sandbox, entity.Id(index))
	}
	for index := range numEntities / 2 {
		suite.assertEntity(entity.Id(index), entityNotLinkedMsg, index)
	}
	Update(suite.sandbox)
	for index := range numEntities / 2 {
		UnlinkEntity(suite.sandbox, entity.Id(index))
	}
	for index := range numEntities / 2 {
		suite.assertDeletedEntity(entity.Id(index), entityNotUnlinkedMsg, index)
	}
}

func (suite *SandboxTestSuite) TestEntityRecycling() {
	for index := range numEntities {
		entityId := LinkEntity(suite.sandbox)
		assert.Equal(suite.T(), entity.Id(index), entityId, entityLinkedOrderMsg, index)
		suite.assertEntity(entity.Id(index), entityNotLinkedMsg, index)
		UnlinkEntity(suite.sandbox, entity.Id(index))
		suite.assertEntity(entity.Id(index), entityUnlinkedBeforeUpdateMsg, index)
	}

	Update(suite.sandbox)

	for index := range numEntities {
		suite.assertDeletedEntity(entity.Id(index), entityNotUnlinkedMsg, index)
	}

	for index := range numEntities {
		entityId := LinkEntity(suite.sandbox)
		assert.Equal(suite.T(), entity.Id(index), entityId, entityNotRecycledMsg, index)
		suite.assertEntity(entity.Id(index), entityNotLinkedMsg, index)
	}
}

func (suite *SandboxTestSuite) TestEntityRecyclingRandom() {
	randomEntityIds := getRandomIds(numEntities, numRemoves)

	for range numEntities {
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
	for index := range numEntities {
		LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(index)).X = float64(index)
		suite.assertComponent(suite.positionLinker, entity.Id(index), componentValueMsg, positionComponent, index)
	}

	Update(suite.sandbox)

	for index := range numEntities {
		suite.assertComponent(suite.positionLinker, entity.Id(index), componentNotLinkedMsg, positionComponent, index)
		assert.Equal(suite.T(), float64(index), suite.positionLinker.Get(entity.Id(index)).X, componentValueMsg, positionComponent, index)
	}
}

func (suite *SandboxTestSuite) TestDuplicateLinkingComponent() {
	entityId := LinkEntity(suite.sandbox)
	xValue := 100.5

	suite.positionLinker.Link(entityId)
	suite.assertComponent(suite.positionLinker, entityId, componentNotLinkedMsg, positionComponent, entityId)

	suite.positionLinker.Get(entityId).X = xValue
	assert.Nil(suite.T(), suite.positionLinker.Link(entityId))
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
	for index := range numEntities {
		LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(index)).X = float64(index)
		suite.assertComponent(suite.positionLinker, entity.Id(index), componentNotLinkedMsg, positionComponent, index)
	}
	for index := range numEntities / 2 {
		suite.velocityLinker.Link(entity.Id(index)).Y = float64(index * 2)
		suite.assertComponent(suite.velocityLinker, entity.Id(index), componentNotLinkedMsg, velocityComponent, index)
	}
	for index := numEntities / 2; index < numEntities; index++ {
		suite.healthLinker.Link(entity.Id(index)).value = float64(index / 2)
		suite.assertComponent(suite.healthLinker, entity.Id(index), componentNotLinkedMsg, healthComponent, index)
	}
	for index := numEntities / 4; index < numEntities*3/4; index++ {
		suite.armorLinker.Link(entity.Id(index)).value = index / 10
		suite.assertComponent(suite.armorLinker, entity.Id(index), componentNotLinkedMsg, armorComponent, index, index)
	}

	Update(suite.sandbox)

	for index := range numEntities {
		suite.assertComponent(suite.positionLinker, entity.Id(index), componentNotLinkedMsg, positionComponent, index)
		assert.Equal(suite.T(), float64(index), suite.positionLinker.Get(entity.Id(index)).X, componentValueMsg, positionComponent, index)
	}
	for index := range numEntities / 2 {
		suite.assertComponent(suite.velocityLinker, entity.Id(index), componentNotLinkedMsg, velocityComponent, index)
		assert.Equal(suite.T(), float64(index*2), suite.velocityLinker.Get(entity.Id(index)).Y, componentValueMsg, velocityComponent, index)
	}
	for index := numEntities / 2; index < numEntities; index++ {
		suite.assertComponent(suite.healthLinker, entity.Id(index), componentNotLinkedMsg, healthComponent, index)
		assert.Equal(suite.T(), float64(index/2), suite.healthLinker.Get(entity.Id(index)).value, componentValueMsg, healthComponent, index)

	}
	for index := numEntities / 4; index < numEntities*3/4; index++ {
		suite.assertComponent(suite.armorLinker, entity.Id(index), componentNotLinkedMsg, armorComponent, index, index)
		assert.Equal(suite.T(), index/10, suite.armorLinker.Get(entity.Id(index)).value, componentValueMsg, armorComponent, index)
	}
}

func (suite *SandboxTestSuite) TestSimpleUnlinkingComponents() {
	removedEntityId := entity.Id(rand.IntN(numEntities))
	removedComponentEntityId := entity.Id(rand.IntN(numEntities))

	for index := range numEntities {
		LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(index))
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

	for index := range numEntities {
		if index != int(removedEntityId) && index != int(removedComponentEntityId) {
			suite.assertComponent(suite.positionLinker, entity.Id(index), componentNotLinkedMsg, positionComponent, index)
		}
	}
}

func (suite *SandboxTestSuite) TestHeavyUnlinkingComponents() {
	for index := range numEntities {
		LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(index))
	}
	for index := range numEntities / 2 {
		suite.velocityLinker.Link(entity.Id(index))
	}
	for index := numEntities / 2; index < numEntities; index++ {
		suite.healthLinker.Link(entity.Id(index))
	}
	for index := numEntities / 4; index < numEntities*3/4; index++ {
		suite.armorLinker.Link(entity.Id(index))
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

func (suite *SandboxTestSuite) TestHeavyTags() {
	tagA := "A"
	tagB := "B"
	tagC := "C"
	tagD := "D"
	aHandler := TagLinker(suite.sandbox, tagA)
	bHandler := TagLinker(suite.sandbox, tagB)
	cHandler := TagLinker(suite.sandbox, tagC)
	dHandler := TagLinker(suite.sandbox, tagD)

	for index := range numEntities {
		LinkEntity(suite.sandbox)
		aHandler.Link(entity.Id(index))
	}
	for index := range numEntities / 2 {
		bHandler.Link(entity.Id(index))
	}
	for index := numEntities / 2; index < numEntities; index++ {
		cHandler.Link(entity.Id(index))
	}
	for index := numEntities / 4; index < numEntities*3/4; index++ {
		dHandler.Link(entity.Id(index))
	}

	removedEntities := getRandomIds(numEntities, numRemoves)
	removedA := getRandomIds(numEntities, numRemoves)
	removedB := getRandomIds(numEntities/2, numRemoves)
	removedC := getRandomIds(numEntities/4, numRemoves)
	removedD := getRandomIds(numEntities/3, numRemoves)

	for _, entityId := range removedEntities {
		UnlinkEntity(suite.sandbox, entityId)
	}
	for _, entityId := range removedA {
		aHandler.Unlink(entityId)
	}
	for _, entityId := range removedB {
		bHandler.Unlink(entityId)
	}
	for _, entityId := range removedC {
		cHandler.Unlink(entityId)
	}
	for _, entityId := range removedD {
		dHandler.Unlink(entityId)
	}

	Update(suite.sandbox)

	for _, entityId := range removedEntities {
		suite.assertDeletedEntity(entityId, entityNotUnlinkedMsg, entityId)
		suite.assertDeletedComponent(aHandler, entityId, componentNotUnlinkedMsg, tagA, entityId)
		suite.assertDeletedComponent(bHandler, entityId, componentNotUnlinkedMsg, tagB, entityId)
		suite.assertDeletedComponent(cHandler, entityId, componentNotUnlinkedMsg, tagC, entityId)
		suite.assertDeletedComponent(dHandler, entityId, componentNotUnlinkedMsg, tagD, entityId)
	}
	for _, entityId := range removedA {
		suite.assertDeletedComponent(aHandler, entityId, componentNotUnlinkedMsg, tagA, entityId)
	}
	for _, entityId := range removedB {
		suite.assertDeletedComponent(bHandler, entityId, componentNotUnlinkedMsg, tagB, entityId)
	}
	for _, entityId := range removedC {
		suite.assertDeletedComponent(cHandler, entityId, componentNotUnlinkedMsg, tagC, entityId)
	}
	for _, entityId := range removedD {
		suite.assertDeletedComponent(dHandler, entityId, componentNotUnlinkedMsg, tagD, entityId)
	}
}

func (suite *SandboxTestSuite) TestSimplePooling() {
	for index := range numEntities {
		LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(index))
		suite.positionLinker.Get(entity.Id(index)).X = float64(index + numEntities)
	}
	Update(suite.sandbox)
	for index := range numEntities {
		suite.positionLinker.Unlink(entity.Id(index))
	}
	Update(suite.sandbox)
	for index := range numEntities {
		suite.positionLinker.Link(entity.Id(index))
	}
	Update(suite.sandbox)

	poolSize := 0
	if suite.pooling {
		poolSize = DefaultPoolCapacity
	}
	for index := range poolSize {
		assert.NotEqual(suite.T(), float64(0), suite.positionLinker.Get(entity.Id(index)).X, componentValueMsg, positionComponent, index)
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
	for index := range numEntities {
		LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(index))
	}
	assert.Len(suite.T(), positionFilter.EntityIds(), 0, filterIncorrectNumEntitiesMsg)
	Update(suite.sandbox)
	assert.Len(suite.T(), positionFilter.EntityIds(), numEntities, filterIncorrectNumEntitiesMsg)
}

func (suite *SandboxTestSuite) TestHeavyFilter() {
	armorHandler := TagLinker(suite.sandbox, armorComponent)
	filter1 := Filter(suite.sandbox, filter.Match2[position, velocity](), filter.ExcludeTags(armorComponent))
	filter2 := Filter(suite.sandbox, filter.Match2[position, velocity](), filter.MatchTags(armorComponent))
	filter3 := Filter(suite.sandbox, filter.Union[velocity](), filter.UnionTags(armorComponent))
	filter4 := Filter(suite.sandbox, filter.Union2[velocity, position](), filter.UnionTags(armorComponent))

	filterCount1 := 0
	filterCount2 := 0
	filterCount3 := 0
	filterCount4 := 0

	for range numEntities {
		LinkEntity(suite.sandbox)
	}
	addedPositionCount := numEntities / 5
	for index := range addedPositionCount {
		suite.positionLinker.Link(entity.Id(index))
	}
	addedVelocityCount := numEntities - numEntities/4
	for index := numEntities / 4; index < numEntities*3/4; index++ {
		suite.velocityLinker.Link(entity.Id(index))
	}
	for index := numEntities / 2; index < numEntities; index++ {
		armorHandler.Link(entity.Id(index))
	}

	assert.Len(suite.T(), filter1.EntityIds(), filterCount1, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter2.EntityIds(), filterCount2, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter3.EntityIds(), filterCount3, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter4.EntityIds(), filterCount4, filterIncorrectNumEntitiesMsg)

	Update(suite.sandbox)

	filterCount3 += addedVelocityCount
	filterCount4 += addedVelocityCount + addedPositionCount

	assert.Len(suite.T(), filter1.EntityIds(), filterCount1, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter2.EntityIds(), filterCount2, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter3.EntityIds(), filterCount3, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter4.EntityIds(), filterCount4, filterIncorrectNumEntitiesMsg)

	addedVelocityCount = numEntities / 10
	for index := range addedVelocityCount {
		suite.velocityLinker.Link(entity.Id(index))
	}
	Update(suite.sandbox)

	filterCount1 += addedVelocityCount
	filterCount3 += addedVelocityCount

	assert.Len(suite.T(), filter1.EntityIds(), filterCount1, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter2.EntityIds(), filterCount2, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter3.EntityIds(), filterCount3, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter4.EntityIds(), filterCount4, filterIncorrectNumEntitiesMsg)

	addedArmorCount := numEntities / 14
	for index := range addedArmorCount {
		armorHandler.Link(entity.Id(index))
	}
	Update(suite.sandbox)

	filterCount1 -= addedArmorCount
	filterCount2 += addedArmorCount

	assert.Len(suite.T(), filter1.EntityIds(), filterCount1, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter2.EntityIds(), filterCount2, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter3.EntityIds(), filterCount3, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter4.EntityIds(), filterCount4, filterIncorrectNumEntitiesMsg)
}
