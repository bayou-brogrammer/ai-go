package game

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs/components"
	"github.com/sirupsen/logrus"
)

// checkCollision checks if a given position is a valid move
func (g *Game) checkCollision(pos gruid.Point) bool {
	if !g.dungeon.InBounds(pos) {
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
		if g.ecs.HasComponent(id, components.CPosition) {
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
	if !g.dungeon.InBounds(newPos) || !g.dungeon.isWalkable(newPos) {
		// TODO: Differentiate between bumping wall and out of bounds?
		return false, nil // Bumped into wall or edge
	}

	// Check for collision with other entities at the target position
	for _, otherID := range g.ecs.GetEntitiesAtWithComponents(newPos, components.CBlocksMovement) {
		if otherID == entityID {
			continue // Don't interact with self
		}

		// Check if the target entity has health (i.e., is attackable)
		if g.ecs.HasComponent(otherID, components.CHealth) {
			// Target is attackable. Queue an AttackAction for the bumping entity.
			logrus.Debugf("Entity %d bumping into attackable entity %d. Queuing AttackAction.", entityID, otherID)

			// Get the TurnActor component of the bumping entity to queue the action
			actor, actorOk := g.ecs.GetTurnActor(entityID)
			if !actorOk {
				// This should not happen if the entity bumping can take turns
				return false, fmt.Errorf("entity %d cannot perform actions (missing TurnActor)", entityID)
			}

			// Create and queue the attack action
			attackAction := AttackAction{
				AttackerID: entityID,
				TargetID:   otherID,
			}
			actor.AddAction(attackAction)

			// Return moved=false because the bump resulted in an attack, not movement.
			// The turn cost will be handled by the AttackAction itself when executed.
			return false, nil
		} else {
			// Bumped into a non-attackable entity (e.g., another player, item, scenery)
			logrus.Debugf("Entity %d bumped into non-attackable entity %d.", entityID, otherID)
			return false, nil // Block movement
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
