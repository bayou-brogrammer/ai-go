package game

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
)

// checkCollision checks if a given position is a valid move
func (g *Game) checkCollision(pos gruid.Point) bool {
	if !g.Map.InBounds(pos) {
		return true // Out of bounds
	}

	// Check for blocking entities (excluding the entity trying to move)
	for _, id := range g.ecs.EntitiesAt(pos) {
		// We need the ID of the entity trying to move to avoid self-collision check
		// This function needs the moving entity's ID passed in.
		// Let's assume it's passed as 'movingEntityID' for now.
		// if id == movingEntityID {
		//  continue
		// }
		// TODO: Refactor checkCollision to accept the moving entity's ID
		if _, ok := g.ecs.GetPosition(id); ok { // Simplified check for now
			return true // Collision with *any* other entity at the target position
		}
	}

	return false
}

// EntityBump attempts to move the entity with the given ID by the delta.
// It checks for map boundaries and collisions with other entities.
// It returns true if the entity successfully moved, false otherwise (due to wall or collision).
func (g *Game) EntityBump(entityID ecs.EntityID, delta gruid.Point) (moved bool, err error) {
	currentPos, ok := g.ecs.GetPosition(entityID)
	if !ok {
		return false, fmt.Errorf("entity %d position not found", entityID)
	}

	newPos := currentPos.Add(delta)

	// Check map bounds and walkability first
	if !g.Map.InBounds(newPos) || !g.Map.isWalkable(newPos) {
		// TODO: Differentiate between bumping wall and out of bounds?
		return false, nil // Bumped into wall or edge
	}

	// Check for collision with other entities at the target position
	// TODO: Refine collision check to potentially allow swapping or attacking
	for _, otherID := range g.ecs.EntitiesAt(newPos) {
		if otherID != entityID { // Don't collide with self
			// For now, any other entity blocks movement
			// TODO: Implement attack logic here if the other entity is hostile
			return false, nil // Bumped into another entity
		}
	}

	// If no collision, move the entity
	err = g.ecs.MoveEntity(entityID, newPos)
	if err != nil {
		return false, fmt.Errorf("failed to move entity %d: %w", entityID, err)
	}

	// Successfully moved
	return true, nil
}
