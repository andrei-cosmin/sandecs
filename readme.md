Creating an entity:
	Allocate next available entityId
	Set entityId bit in entityLinker.linkedEntities
	Mark entityLinker state as dirtied
	Return entityId


Removing an entity:
	Clear bit in entityLinker.linkedEntities
	Set bit in entityLinker.scheduledRemoves
	Mark entityLinker state as dirtied


Attaching componentInstance to entity:
	Get componentType (same as component name)
	Get componentResolver for componentType
	Set entityId bit in componentResolver.linkedEntities
	Set componentInstance in componentResolver.components at entityId position


Removing component from entity:
	Get componentType (same as component name)
	Get componentResolver for componentType
	Check if entityId bit is set in componentResolver.linkedEntities
	If entityId bit is set, set entityId bit in componentResolver.scheduledRemoves
	Mark componentLinker state as dirtied


Get componentInstance of componentType from entityId:
	Use componentType to find out componentId
	Get componentResolver for componentType
	Check if entityId bit is set in componentResolver.linkedEntities
	If entityId bit is set, get componentInstance from componentResolver.components[entityId]


Handling new components when encountering them:
Getting a componentResolver using componentType
	Use componentType to find out componentId
		If componentType doesn't exist create new componentResolver
		Use next avaialable compnentId
		Add componentResolver to array, using componentId
		Store componentType->componentId in map
		Return componentId
	Take componentResolver from array of resolvers, using componentId


On update:
	If componentLinker or entityLinker are not dirtied, skip
	Difference between entityLinker.linked entities and entityLinker.scheduledRemoves
	For every componentResolver:
		Union between componentResolver.scheduledRemoves and entityLinker.scheduledRemoves
		For each bit set in union clear componentResolver.components[bit] (bitset iterates checking words of 64 bits and findeing indexes of set bits)
		Difference between current componentResolver.linkedEntities and union
	Update QueryRegistry
		For every query registered update entities


------------------------------------------------------------------------------------------------------------------

Query creation:
	MatchAll: all components from the set must be present
	OneOf: only one component from the set must be present
	Exclude: none of the components from the set must be present
	After all sets have been defined generate a hash string


Query Registration
	Find queryId using queryHash 
	If queryId found, return queryCache (linked to queryId)
	// Create queryCache
	For every matchCompnent:
		Use componentType to find out matchCompnentId
		Store matchedComponentId

	For every oneOfComponent:
		Use componentType to find out oneOfComponentId
		Store oneOfComponentId

	For every excludeComponent:
		Use componentType to find out excludedComponentId
		Store excludedComponentId

	Add queryHash->queryId into queryMap
	Set queryCaches[queryId] = queryCache

Query Update:
	buffer = entityLinker.linkedEntities
	For every matchCompnentId:
		buffer = buffer intersection matchComponentResolver.linkedEntities

	For every oneOfComponent:
		buffer = buffer union oneOfComponentResolver.linkedEntities

	For every excludeComponent:
		buffer = buffer difference excludedComponentResolver.linkedEntities

	Check if buffer differes from queryCache.linkedEntities
	If it differs, rebuild queryCache.entityIds (list containing explicit ids of entities matching the queries)



