package main

import (
	"reflect"

	"codeberg.org/anaseto/gruid"
)

// checkCollision checks if a given position is a valid move
func (g *game) checkCollision(pos gruid.Point) bool {
	if !g.Map.InBounds(pos) {
		return true // Out of bounds
	}
	// Check for blocking entities
	for _, id := range g.ecs.QueryEntitiesWithComponents(reflect.TypeOf(Position{}), reflect.TypeOf(BlocksMovement{})) {
		if id == g.PlayerID {
			continue
		}
		otherPos, ok := g.ecs.GetComponent(id, reflect.TypeOf(Position{}))
		if ok && otherPos.(Position).Point == pos {
			return true // Collision with blocking entity
		}
	}
	return false
}

func (g *game) PlayerBump(delta gruid.Point) (bool, error) {
	pid := g.PlayerID

	// Update player position in ECS
	positionType := reflect.TypeOf(Position{})
	if comps, ok := g.ecs.entities[pid]; ok {
		if pos, found := comps[positionType]; found {
			currentPos := pos.(Position).Point
			newPos := currentPos.Add(delta)

			if !g.Map.isWalkable(newPos) { // Check if the new position is walkable
				return false, nil // Not walkable, don't move
			}

			g.ecs.entities[pid][positionType] = Position{Point: newPos} // Save the updated position back
			return true, nil                                            // Movement successful
		}
	} else {
		//Very unexpected, but handle the case where the player entity doesn't exist
		return false, nil
	}

	// If we reach here, something went wrong
	return false, nil // Movement failed
}
