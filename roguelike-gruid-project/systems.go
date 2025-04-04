package main

import (
	"log"

	"codeberg.org/anaseto/gruid"
)

// MovementEvent represents a request to move an entity
type MovementEvent struct {
	EntityID EntityID
	Delta    gruid.Point
}

// RenderSystem draws all entities with Position and Renderable components onto the grid.
func RenderSystem(ecs *World, grid gruid.Grid) {
	// Get component types needed for querying
	posType := GetReflectType(Position{})
	renType := GetReflectType(Renderable{})

	// Query for entities that have both Position and Renderable
	entityIDs := ecs.QueryEntitiesWithComponents(posType, renType)

	// Iterate through entities and draw them
	for _, id := range entityIDs {
		// Retrieve components (we know they exist from the query)
		posComp, _ := ecs.GetComponent(id, posType)
		renComp, _ := ecs.GetComponent(id, renType)

		// Type assert components to access their fields
		pos := posComp.(Position)
		ren := renComp.(Renderable)

		// Set the cell in the grid
		// Note: This assumes the grid has been cleared beforehand.
		// We might add map bounds checking later.
		grid.Set(pos.Point, gruid.Cell{
			Rune:  ren.Glyph,
			Style: gruid.Style{Fg: ren.Color}, // Assuming default background for now
		})
	}
}

// TurnSystem processes turns for entities
func TurnSystem(g *game) {
	processTurns(g)
	monstersTurn(g)
}

// process_turns function
func processTurns(g *game) {
	world := g.ecs
	turnQueue := g.turnQueue

	// Periodically clean up the queue
	metrics := turnQueue.CleanupDeadEntities(world)

	if metrics.EntitiesRemoved > 10 {
		log.Printf(
			"Turn queue cleanup: removed %d entities in %v",
			metrics.EntitiesRemoved,
			metrics.ProcessingTime,
		)
	}

	turnQueue.PrintQueue()

	for {
		entry, ok := turnQueue.Next()
		if !ok {
			break
		}

		entity := entry.EntityID
		time := entry.Time

		actorComponent, ok := world.GetComponent(entity, GetReflectType(TurnActor{}))
		if !ok {
			log.Printf("Actor not found: %v", entity)
			continue
		}

		actor, ok := actorComponent.(TurnActor)
		if !ok {
			log.Printf("Invalid turn_actor component for entity: %v", entity)
			continue
		}

		if !actor.IsAlive() {
			log.Printf("Actor is dead. Why is it still in the queue?")
			continue
		}

		hasAction := actor.PeakNextAction() != nil
		_, isPlayer := world.GetComponent(entity, GetReflectType(Name{}))

		if isPlayer && !hasAction {
			log.Printf("Player is awaiting input: %v", entity)
			world.AddComponent(entity, WaitingForInput{})
			turnQueue.Add(entity, time)
			return
		}

		action := actor.NextAction()
		if action == nil {
			log.Printf("No action for entity: %v. Rescheduling turn.", entity)
			turnQueue.Add(entity, time)
			return
		}

		dTime, err := action.Execute(world)
		if err != nil {
			log.Printf("Failed to execute action: %v", err)
			turnQueue.Add(entity, time)
			return
		}

		turnQueue.Add(entity, time+uint64(dTime))

	}
}

// monsters_turn function
func monstersTurn(g *game) {
	//log.Println("Monsters turn")

	//for entity, components := range world.Entities {
	//	_, isAI := components["ai_tag"]
	//	_, isPlayer := components["player_tag"]
	//	_, isDead := components["dead_tag"]

	//	if isAI && !isPlayer && !isDead {
	//		turnActorComponent, ok := components["turn_actor"]
	//		if !ok {
	//			log.Printf("Entity %d has ai_tag but no turn_actor", entity)
	//			continue
	//		}
	//		turnActor, ok := turnActorComponent.(TurnActor)
	//		if !ok {
	//			log.Printf("Entity %d has invalid turn_actor component", entity)
	//			continue
	//		}

	//		if turnActor.PeakNextAction() != nil {
	//			continue
	//		}

	//		positionComponent, ok := components["position"]
	//		if !ok {
	//			log.Printf("Entity %d has ai_tag but no position", entity)
	//			continue
	//		}
	//		position, ok := positionComponent.(Position)
	//		if !ok {
	//			log.Printf("Entity %d has invalid position component", entity)
	//			continue
	//		}

	//		// Try different directions in a random order
	//		directions := AllDirections()
	//		var validDirection *MoveDirection

	//		for i := 0; i < len(directions); i++ {
	//			direction := RandomDirection()
	//			newPosition := position.Add(direction)

	//			currentMapResource, ok := world.Resources["current_map"]
	//			if !ok {
	//				log.Println("current_map resource not found")
	//				continue
	//			}
	//			currentMap, ok := currentMapResource.(CurrentMap)
	//			if !ok {
	//				log.Println("Invalid current_map resource")
	//				continue
	//			}

	//			terrainEntity, ok := currentMap.GetTerrain(newPosition)
	//			if !ok {
	//				continue // No terrain at this position
	//			}

	//			terrainTypeComponent, ok := world.GetComponent(terrainEntity, "terrain_type")
	//			if !ok {
	//				continue // No terrain type at this position
	//			}
	//			terrainType, ok := terrainTypeComponent.(TerrainType)
	//			if !ok {
	//				continue // Invalid terrain type
	//			}

	//			if terrainType == Floor {
	//				validDirection = &direction
	//				break
	//			}
	//		}

	//		if validDirection != nil {
	//			log.Printf("AI entity %d moving in direction %v", entity, validDirection)
	//			walkAction := NewWalkBuilder().
	//				WithEntity(entity).
	//				WithDirection(*validDirection).
	//				Build()
	//			turnActor.AddAction(walkAction)
	//			world.Entities[entity]["turn_actor"] = turnActor // Update the turn actor
	//		} else {
	//			log.Printf("AI entity %d has no valid move, waiting", entity)
	//			// Here we could add a wait action if we had one
	//		}
	//	}
	//}

	////nextState.Set(ProcessTurns) // How to set next state in Go?
}
