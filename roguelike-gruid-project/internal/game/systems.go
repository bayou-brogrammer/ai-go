package game

import (
	"fmt"
	"math/rand"

	"codeberg.org/anaseto/gruid"
	"github.com/sirupsen/logrus"
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
		grid.Set(pos, gruid.Cell{
			Rune:  ren.Glyph,
			Style: gruid.Style{Fg: ren.Color}, // Use chosen color
		})
	}
}

// handleMonsterTurn determines and executes an action for a monster.
// It returns the cost of the action taken.
func handleMonsterTurn(g *Game, entityID EntityID) (cost uint, err error) {
	logrus.Debugf("Handling monster turn for entity %d", entityID)

	// Get monster position
	_, ok := g.ecs.GetPosition(entityID)
	if !ok {
		return 0, fmt.Errorf("monster %d has no position", entityID)
	}

	// Simple random walk AI
	possibleMoves := []gruid.Point{
		{X: -1, Y: 0}, // West
		{X: 1, Y: 0},  // East
		{X: 0, Y: -1}, // North
		{X: 0, Y: 1},  // South
	}

	// Shuffle directions
	rand.Shuffle(len(possibleMoves), func(i, j int) {
		possibleMoves[i], possibleMoves[j] = possibleMoves[j], possibleMoves[i]
	})

	// Try to find a valid move
	for _, delta := range possibleMoves {
		// Check bounds, walkability, and collision (using EntityBump logic implicitly)
		moved, bumpErr := g.EntityBump(entityID, delta)
		if bumpErr != nil {
			// Log error but maybe try another direction?
			logrus.Errorf("Error during monster %d bump check: %v", entityID, bumpErr)
			continue // Try next direction
		}

		if moved {
			// Successfully moved, return move cost (defined in MoveAction.Execute)
			// We technically already executed the move in EntityBump, so just get the cost.
			// Let's assume MoveAction cost is standard for now.
			// TODO: Refactor MoveAction/EntityBump if Execute shouldn't *do* the move.
			return 100, nil
		}
		// If !moved and err == nil, it means a bump occurred, try another direction.
	}

	// If no valid move found, monster waits (costs time)
	logrus.Debugf("Skipping monster %d.", entityID) // Optional debug log
	return 100, nil                                 // Return standard action cost for waiting
}
