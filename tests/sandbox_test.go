package tests

import (
	"github.com/andrei-cosmin/sandecs"
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/filter"
	"github.com/andrei-cosmin/sandecs/options"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math/rand/v2"
	"slices"
	"testing"
)

var pooledModes = map[options.Mode]bool{
	options.Pooled: true,
}

func TestSandboxSuite(t *testing.T) {
	suite.Run(t, &SandboxTestSuite{mode: options.Standard, poolSize: 0})
	suite.Run(t, &SandboxTestSuite{mode: options.Pooled, poolSize: 20000})
	suite.Run(t, &SandboxTestSuite{mode: options.Compact, poolSize: 20000})
}

type SandboxTestSuite struct {
	sandboxSuite
	mode           options.Mode
	poolSize       uint
	positionLinker component.Linker[position]
	velocityLinker component.Linker[velocity]
	renderedLinker component.Linker[rendered]
	armorLinker    component.Linker[armor]
	healthLinker   component.Linker[health]
	nameLinker     component.Linker[name]
}

func (suite *SandboxTestSuite) SetupTest() {
	suite.T().Log("Pooling size:", suite.poolSize)
	suite.sandbox = sandbox.New(suite.mode, options.DefaultNumEntities, options.DefaultNumComponents, suite.poolSize)

	suite.positionLinker = sandbox.ComponentLinker[position](suite.sandbox)
	suite.velocityLinker = sandbox.ComponentLinker[velocity](suite.sandbox)
	suite.renderedLinker = sandbox.ComponentLinker[rendered](suite.sandbox)
	suite.armorLinker = sandbox.ComponentLinker[armor](suite.sandbox)
	suite.healthLinker = sandbox.ComponentLinker[health](suite.sandbox)
	suite.nameLinker = sandbox.ComponentLinker[name](suite.sandbox)
}

func (suite *SandboxTestSuite) TestSandbox_LinkingEntity() {
	for range numEntities {
		sandbox.LinkEntity(suite.sandbox)
	}
	sandbox.Update(suite.sandbox)
	for index := range numEntities {
		suite.assertEntity(entity.Id(index), entityNotLinkedMsg, index)
	}
}

func (suite *SandboxTestSuite) TestSandbox_UnlinkingEntity() {
	for range numEntities {
		sandbox.LinkEntity(suite.sandbox)
	}
	sandbox.Update(suite.sandbox)
	for index := range numEntities / 2 {
		sandbox.UnlinkEntity(suite.sandbox, entity.Id(index))
	}
	sandbox.Update(suite.sandbox)
	for index := range numEntities / 2 {
		suite.assertDeletedEntity(entity.Id(index), entityNotUnlinkedMsg, index)
	}
	for index := numEntities / 2; index < numEntities; index++ {
		suite.assertEntity(entity.Id(index), entityNotLinkedMsg, index)
	}
}

func (suite *SandboxTestSuite) TestSandbox_DuplicateUnlinkingEntity() {
	for range numEntities {
		sandbox.LinkEntity(suite.sandbox)
	}
	sandbox.Update(suite.sandbox)
	for index := range numEntities / 2 {
		sandbox.UnlinkEntity(suite.sandbox, entity.Id(index))
	}
	for index := range numEntities / 2 {
		sandbox.UnlinkEntity(suite.sandbox, entity.Id(index))
	}
	for index := range numEntities / 2 {
		suite.assertEntity(entity.Id(index), entityNotLinkedMsg, index)
	}
	sandbox.Update(suite.sandbox)
	for index := range numEntities / 2 {
		sandbox.UnlinkEntity(suite.sandbox, entity.Id(index))
	}
	for index := range numEntities / 2 {
		suite.assertDeletedEntity(entity.Id(index), entityNotUnlinkedMsg, index)
	}
}

func (suite *SandboxTestSuite) TestSandbox_EntityRecycling() {
	for index := range numEntities {
		entityId := sandbox.LinkEntity(suite.sandbox)
		assert.Equal(suite.T(), entity.Id(index), entityId, entityLinkedOrderMsg, index)
		suite.assertEntity(entity.Id(index), entityNotLinkedMsg, index)
		sandbox.UnlinkEntity(suite.sandbox, entity.Id(index))
		suite.assertEntity(entity.Id(index), entityUnlinkedBeforeUpdateMsg, index)
	}

	sandbox.Update(suite.sandbox)

	for index := range numEntities {
		suite.assertDeletedEntity(entity.Id(index), entityNotUnlinkedMsg, index)
	}

	for index := range numEntities {
		entityId := sandbox.LinkEntity(suite.sandbox)
		assert.Equal(suite.T(), entity.Id(index), entityId, entityNotRecycledMsg, index)
		suite.assertEntity(entity.Id(index), entityNotLinkedMsg, index)
	}
}

func (suite *SandboxTestSuite) TestSandbox_EntityRecyclingRandom() {
	randomEntityIds := getRandomIds(numEntities, numRemoves)

	for range numEntities {
		sandbox.LinkEntity(suite.sandbox)
	}
	for _, entityId := range randomEntityIds {
		sandbox.UnlinkEntity(suite.sandbox, entityId)
	}
	sandbox.Update(suite.sandbox)
	slices.Sort(randomEntityIds)
	for _, entityId := range randomEntityIds {
		suite.assertDeletedEntity(entityId, entityNotUnlinkedMsg, entityId)
		newId := sandbox.LinkEntity(suite.sandbox)
		assert.Equal(suite.T(), entityId, newId, entityRecycleOrderMsg, newId)
	}
}

func (suite *SandboxTestSuite) TestSandbox_SimpleLinkingComponents() {
	for index := range numEntities {
		sandbox.LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(index)).X = float64(index)
		suite.assertComponent(suite.positionLinker, entity.Id(index), componentValueMsg, positionComponent, index)
	}

	sandbox.Update(suite.sandbox)

	for index := range numEntities {
		suite.assertComponent(suite.positionLinker, entity.Id(index), componentNotLinkedMsg, positionComponent, index)
		assert.Equal(suite.T(), float64(index), suite.positionLinker.Get(entity.Id(index)).X, componentValueMsg, positionComponent, index)
	}
}

func (suite *SandboxTestSuite) TestSandbox_DuplicateLinkingComponent() {
	entityId := sandbox.LinkEntity(suite.sandbox)
	xValue := 100.5

	suite.positionLinker.Link(entityId)
	suite.assertComponent(suite.positionLinker, entityId, componentNotLinkedMsg, positionComponent, entityId)

	suite.positionLinker.Get(entityId).X = xValue
	assert.Nil(suite.T(), suite.positionLinker.Link(entityId))
	assert.Equal(suite.T(), xValue, suite.positionLinker.Get(entityId).X, componentValueMsg, positionComponent, entityId)
}

func (suite *SandboxTestSuite) TestSandbox_DuplicateUnlinkingComponent() {
	entityId := sandbox.LinkEntity(suite.sandbox)
	suite.positionLinker.Link(entityId)
	suite.assertComponent(suite.positionLinker, entityId, componentNotLinkedMsg, positionComponent, entityId)

	suite.positionLinker.Unlink(entityId)
	suite.positionLinker.Unlink(entityId)
	suite.assertComponent(suite.positionLinker, entityId, componentNotLinkedMsg, positionComponent, entityId)
	sandbox.Update(suite.sandbox)
	assert.Nil(suite.T(), suite.positionLinker.Get(entityId))

	suite.assertDeletedComponent(suite.positionLinker, entityId, componentNotUnlinkedMsg, positionComponent, entityId)
	assert.Nil(suite.T(), suite.positionLinker.Get(entityId))
	suite.positionLinker.Unlink(entityId)
	suite.positionLinker.Unlink(entityId)
	sandbox.Update(suite.sandbox)
	suite.assertDeletedComponent(suite.positionLinker, entityId, componentNotUnlinkedMsg, positionComponent, entityId)
	assert.Nil(suite.T(), suite.positionLinker.Get(entityId))
}

func (suite *SandboxTestSuite) TestSandbox_ComponentForInvalidEntity() {
	suite.positionLinker.Link(10 * numEntities)
	suite.assertDeletedComponent(suite.positionLinker, 10*numEntities, componentNotUnlinkedMsg, positionComponent, 10*numEntities)
	suite.positionLinker.Unlink(10 * numEntities)
	sandbox.Update(suite.sandbox)
	suite.assertDeletedComponent(suite.positionLinker, 10*numEntities, componentNotLinkedMsg, positionComponent, 10*numEntities)
}

func (suite *SandboxTestSuite) TestSandbox_ComponentInstanceCleared() {
	for range numEntities/10 + 1 {
		sandbox.LinkEntity(suite.sandbox)
	}
	suite.positionLinker.Link(numEntities / 10)
	sandbox.Update(suite.sandbox)
	suite.positionLinker.Unlink(numEntities / 10)
	sandbox.Update(suite.sandbox)
	suite.assertDeletedComponent(suite.positionLinker, numEntities/10, componentNotUnlinkedMsg, positionComponent, numEntities/10)
	assert.Nil(suite.T(), suite.positionLinker.Get(numEntities/10))
}

func (suite *SandboxTestSuite) TestSandbox_HeavyLinkingComponents() {
	for index := range numEntities {
		sandbox.LinkEntity(suite.sandbox)
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

	sandbox.Update(suite.sandbox)

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

func (suite *SandboxTestSuite) TestSandbox_SimpleUnlinkingComponents() {
	removedEntityId := entity.Id(rand.IntN(numEntities))
	removedComponentEntityId := entity.Id(rand.IntN(numEntities))

	for index := range numEntities {
		sandbox.LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(index))
	}

	sandbox.UnlinkEntity(suite.sandbox, removedEntityId)
	suite.positionLinker.Unlink(removedComponentEntityId)

	sandbox.Update(suite.sandbox)

	suite.assertDeletedEntity(removedEntityId, entityNotUnlinkedMsg, removedEntityId)
	suite.assertDeletedComponent(suite.positionLinker, removedEntityId, componentNotUnlinkedMsg, positionComponent, removedEntityId)
	suite.assertDeletedComponent(suite.positionLinker, removedComponentEntityId, componentNotUnlinkedMsg, positionComponent, removedComponentEntityId)

	readdedEntityId := sandbox.LinkEntity(suite.sandbox)
	assert.Equal(suite.T(), removedEntityId, readdedEntityId, entityRecycleOrderMsg, readdedEntityId)
	suite.assertDeletedComponent(suite.positionLinker, readdedEntityId, componentNotUnlinkedMsg, positionComponent, readdedEntityId)

	for index := range numEntities {
		if index != int(removedEntityId) && index != int(removedComponentEntityId) {
			suite.assertComponent(suite.positionLinker, entity.Id(index), componentNotLinkedMsg, positionComponent, index)
		}
	}
}

func (suite *SandboxTestSuite) TestSandbox_HeavyUnlinkingComponents() {
	for index := range numEntities {
		sandbox.LinkEntity(suite.sandbox)
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
		sandbox.UnlinkEntity(suite.sandbox, entityId)
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

	sandbox.Update(suite.sandbox)

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

func (suite *SandboxTestSuite) TestSandbox_HeavyTags() {
	tagA := "A"
	tagB := "B"
	tagC := "C"
	tagD := "D"
	aHandler := sandbox.TagLinker(suite.sandbox, tagA)
	bHandler := sandbox.TagLinker(suite.sandbox, tagB)
	cHandler := sandbox.TagLinker(suite.sandbox, tagC)
	dHandler := sandbox.TagLinker(suite.sandbox, tagD)

	for index := range numEntities {
		sandbox.LinkEntity(suite.sandbox)
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
		sandbox.UnlinkEntity(suite.sandbox, entityId)
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

	sandbox.Update(suite.sandbox)

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

func (suite *SandboxTestSuite) TestSandbox_SimplePooling() {
	for index := range numEntities {
		sandbox.LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(index))
		suite.positionLinker.Get(entity.Id(index)).X = float64(index + numEntities)
	}
	sandbox.Update(suite.sandbox)
	for index := range numEntities {
		suite.positionLinker.Unlink(entity.Id(index))
	}
	sandbox.Update(suite.sandbox)
	for index := range numEntities {
		suite.positionLinker.Link(entity.Id(index))
	}
	sandbox.Update(suite.sandbox)

	if pooledModes[suite.mode] {
		for index := range suite.poolSize {
			assert.NotEqual(suite.T(), float64(index+numEntities), suite.positionLinker.Get(entity.Id(index)).X, componentValueMsg, positionComponent, index)
		}
		for i := suite.poolSize; i < numEntities; i++ {
			assert.Equal(suite.T(), float64(0), suite.positionLinker.Get(entity.Id(i)).X, componentValueMsg, positionComponent, i)
		}
	} else {
		for index := range numEntities {
			assert.NotEqual(suite.T(), float64(index+numEntities), suite.positionLinker.Get(entity.Id(index)).X, componentValueMsg, positionComponent, index)
		}
	}
}

func (suite *SandboxTestSuite) TestSandbox_SimpleFilter() {
	positionFilter := sandbox.Filter(suite.sandbox, filter.Match[position]())
	duplicateFilter := sandbox.Filter(suite.sandbox, filter.Match[position]())
	assert.Equal(suite.T(), positionFilter, duplicateFilter)
	sandbox.Filter(suite.sandbox, filter.Match[position]())

	assert.Len(suite.T(), positionFilter.EntityIds(), 0, filterIncorrectNumEntitiesMsg)
	for index := range numEntities {
		sandbox.LinkEntity(suite.sandbox)
		suite.positionLinker.Link(entity.Id(index))
	}
	assert.Len(suite.T(), positionFilter.EntityIds(), 0, filterIncorrectNumEntitiesMsg)
	sandbox.Update(suite.sandbox)
	assert.Len(suite.T(), positionFilter.EntityIds(), numEntities, filterIncorrectNumEntitiesMsg)

	for _, entityId := range positionFilter.EntityIds() {
		assert.True(suite.T(), positionFilter.EntityMask().Test(entityId))
	}
}

func (suite *SandboxTestSuite) TestSandbox_HeavyFilter() {
	armorTagHandler := sandbox.TagLinker(suite.sandbox, armorComponent)
	filter1 := sandbox.Filter(suite.sandbox, filter.Match2[position, velocity](), filter.ExcludeTags(armorComponent))
	filter2 := sandbox.Filter(suite.sandbox, filter.Match2[position, velocity](), filter.MatchTags(armorComponent))
	filter3 := sandbox.Filter(suite.sandbox, filter.Union[velocity](), filter.UnionTags(armorComponent))
	filter4 := sandbox.Filter(suite.sandbox, filter.Union2[velocity, position](), filter.UnionTags(armorComponent))

	filterCount1 := 0
	filterCount2 := 0
	filterCount3 := 0
	filterCount4 := 0

	for range numEntities {
		sandbox.LinkEntity(suite.sandbox)
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
		armorTagHandler.Link(entity.Id(index))
	}

	assert.Len(suite.T(), filter1.EntityIds(), filterCount1, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter2.EntityIds(), filterCount2, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter3.EntityIds(), filterCount3, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter4.EntityIds(), filterCount4, filterIncorrectNumEntitiesMsg)

	sandbox.Update(suite.sandbox)

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
	sandbox.Update(suite.sandbox)

	filterCount1 += addedVelocityCount
	filterCount3 += addedVelocityCount

	assert.Len(suite.T(), filter1.EntityIds(), filterCount1, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter2.EntityIds(), filterCount2, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter3.EntityIds(), filterCount3, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter4.EntityIds(), filterCount4, filterIncorrectNumEntitiesMsg)

	addedArmorCount := numEntities / 14
	for index := range addedArmorCount {
		armorTagHandler.Link(entity.Id(index))
	}
	sandbox.Update(suite.sandbox)

	filterCount1 -= addedArmorCount
	filterCount2 += addedArmorCount

	assert.Len(suite.T(), filter1.EntityIds(), filterCount1, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter2.EntityIds(), filterCount2, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter3.EntityIds(), filterCount3, filterIncorrectNumEntitiesMsg)
	assert.Len(suite.T(), filter4.EntityIds(), filterCount4, filterIncorrectNumEntitiesMsg)
}

func (suite *SandboxTestSuite) TestSandbox_Hooks() {
	count := 0
	armorHandler := sandbox.TagLinker(suite.sandbox, armorComponent)
	armorHandler.SetLinkHook(func() {
		count++
	})
	armorHandler.SetUnlinkHook(func() {
		count--
	})

	for range numEntities {
		entityId := sandbox.LinkEntity(suite.sandbox)
		armorHandler.Link(entityId)
	}
	for entityId := range numEntities {
		sandbox.UnlinkEntity(suite.sandbox, uint(entityId))
	}
	for entityId := range numEntities {
		armorHandler.Link(uint(entityId))
	}
	for entityId := range numEntities {
		armorHandler.Unlink(uint(entityId))
	}

	sandbox.Update(suite.sandbox)
	assert.Zero(suite.T(), count)

	for range numEntities {
		entityId := sandbox.LinkEntity(suite.sandbox)
		armorHandler.Link(entityId)
	}
	for entityId := range numEntities {
		armorHandler.Unlink(uint(entityId))
	}
	for entityId := range numEntities {
		armorHandler.Link(uint(entityId))
	}
	for entityId := range numEntities {
		armorHandler.Unlink(uint(entityId))
	}
	sandbox.Update(suite.sandbox)
	assert.Zero(suite.T(), count)

	for entityId := range numEntities {
		armorHandler.Link(uint(entityId))
	}
	for entityId := range numEntities {
		armorHandler.Link(uint(entityId))
	}
	for entityId := range numEntities {
		armorHandler.Unlink(uint(entityId))
	}
	for entityId := range numEntities {
		armorHandler.Unlink(uint(entityId))
	}
	sandbox.Update(suite.sandbox)
	assert.Zero(suite.T(), count)
}

func (suite *SandboxTestSuite) TestSandbox_HookTrigger() {
	linkCount := 0
	unlinkCount := 0
	armorHandler := sandbox.TagLinker(suite.sandbox, armorComponent)

	armorHandler.SetLinkHook(func() {
		linkCount++
	})
	suite.positionLinker.SetLinkHook(func(*position) {
		linkCount++
	})
	armorHandler.SetUnlinkHook(func() {
		unlinkCount--
	})
	suite.positionLinker.SetUnlinkHook(func(*position) {
		unlinkCount--
	})

	for range numEntities {
		entityId := sandbox.LinkEntity(suite.sandbox)
		suite.healthLinker.Link(entityId)
	}
	sandbox.Update(suite.sandbox)
	assert.Zero(suite.T(), linkCount)
	assert.Zero(suite.T(), unlinkCount)

	for entityId := range numEntities {
		sandbox.UnlinkEntity(suite.sandbox, uint(entityId))
	}
	sandbox.Update(suite.sandbox)
	assert.Zero(suite.T(), linkCount)
	assert.Zero(suite.T(), unlinkCount)
}
