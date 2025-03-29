package tests

import (
	sandbox "github.com/andrei-cosmin/sandecs"
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math/rand/v2"
)

type sandboxSuite struct {
	suite.Suite
	sandbox *sandbox.Sandbox
}

func (suite *sandboxSuite) assertEntity(entityId entity.Id, message ...interface{}) {
	assert.True(suite.T(), sandbox.IsEntityLinked(suite.sandbox, entityId), message...)
}

func (suite *sandboxSuite) assertDeletedEntity(entityId entity.Id, message ...interface{}) {
	assert.True(suite.T(), !sandbox.IsEntityLinked(suite.sandbox, entityId), message...)
}

func (suite *sandboxSuite) assertComponent(linker component.BasicLinker, entityId entity.Id, message ...interface{}) {
	assert.True(suite.T(), linker.Has(entityId), message...)
}

func (suite *sandboxSuite) assertDeletedComponent(linker component.BasicLinker, entityId entity.Id, message ...interface{}) {
	assert.False(suite.T(), linker.Has(entityId), message...)
}

func getRandomIds(maxId, numIds int) []entity.Id {
	mappedEntityIds := make(map[entity.Id]bool)
	randomEntityIds := make([]entity.Id, 0)

	for len(mappedEntityIds) < numIds {
		entityId := entity.Id(rand.IntN(maxId))
		if _, exists := mappedEntityIds[entityId]; !exists {
			mappedEntityIds[entityId] = true
			randomEntityIds = append(randomEntityIds, entityId)
		}
	}

	return randomEntityIds
}
