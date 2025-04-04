package main

import (
	"codeberg.org/anaseto/gruid"
)

// MovementEvent represents a request to move an entity
type MovementEvent struct {
	EntityID EntityID
	Delta    gruid.Point
}

// RenderSystem draws all entities with Position and Renderable components onto the grid.
func RenderSystem(ecs *ECS, grid gruid.Grid) {
	entityIDs := ecs.GetEntitiesWithPositionAndRenderable()

	// Iterate through entities and draw them
	for _, id := range entityIDs {
		// Retrieve components (we know they exist from the query)
		pos, _ := ecs.GetPosition(id)
		ren, _ := ecs.GetRenderable(id)

		// Set the cell in the grid
		// Note: This assumes the grid has been cleared beforehand.
		// We might add map bounds checking later.
		grid.Set(pos, gruid.Cell{
			Rune:  ren.Glyph,
			Style: gruid.Style{Fg: ren.Color}, // Assuming default background for now
		})
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
