package game

import (
	"fmt"
	"math/rand"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
	"github.com/sirupsen/logrus"
)

// RenderSystem draws all entities with Position and Renderable components onto the grid.
func RenderSystem(ecs *ecs.ECS, grid gruid.Grid) {
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
func handleMonsterTurn(g *Game, entityID ecs.EntityID) (cost uint, err error) {
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

// processTurnQueue processes turns for actors until it's the player's turn
// or the queue is exhausted for the current time step.
// It returns true if the game should wait for player input, false otherwise.
func (md *Model) processTurnQueue() {
	g := md.game
	logrus.Debug("========= processTurnQueue started =========")

	g.turnQueue.CleanupDeadEntities(g.ecs)
	g.turnQueue.PrintQueue()

	for i := range 100 { // Limit to 100 iterations to prevent infinite loops
		logrus.Debugf("Turn queue iteration %d", i)
		logrus.Debugf("Current turn queue time: %d", g.turnQueue.CurrentTime)

		// Check if the queue is empty
		if g.turnQueue.IsEmpty() {
			logrus.Debug("Turn queue is empty.")
			logrus.Debug("========= processTurnQueue ended (queue empty) =========")
			return
		}

		// Peek at the next actor
		entry, ok := g.turnQueue.Peek()
		if !ok { // Should not happen if IsEmpty is false, but good practice
			logrus.Debug("Error: Queue not empty but failed to peek.")
			// Clear active entity
			logrus.Debug("========= processTurnQueue ended (peek error) =========")
			return
		}

		logrus.Debugf("Next entry in queue: EntityID=%d, Time=%d", entry.EntityID, entry.Time)

		// It's time for the next actor's turn, pop them
		actorEntry, _ := g.turnQueue.Next() // We already know ok is true from Peek
		logrus.Debugf("Popped actor from queue: EntityID=%d, Time=%d", actorEntry.EntityID, actorEntry.Time)

		// Update game time to the current actor's time
		g.turnQueue.CurrentTime = actorEntry.Time
		logrus.Debugf("Updated game time to: %d", g.turnQueue.CurrentTime)

		// Check if it's the player
		if actorEntry.EntityID == g.PlayerID {
			logrus.Debug("It's the player's turn.")
			g.ecs.WaitingForInput[g.PlayerID] = true
			logrus.Debug("========= processTurnQueue ended (player's turn) =========")
			return
		}

		// It's a monster's turn
		logrus.Debugf("Handling turn for monster %d at time %d", actorEntry.EntityID, actorEntry.Time)
		if g.ecs.HasAITag(actorEntry.EntityID) { // Ensure it's actually an AI
			cost, err := handleMonsterTurn(g, actorEntry.EntityID)
			if err != nil {
				logrus.Errorf("Error during monster %d turn: %v", actorEntry.EntityID, err)
				// Decide how to handle monster errors - skip turn? Give basic cost?
				// For now, let's give it a standard cost and re-queue
				cost = 100
			}
			// Re-queue the monster for its next turn
			nextTime := g.turnQueue.CurrentTime + uint64(cost)
			logrus.Debugf("Re-queuing monster %d for time %d (cost=%d)", actorEntry.EntityID, nextTime, cost)
			g.turnQueue.Add(actorEntry.EntityID, nextTime)
		} else {
			// Entity in queue that isn't player or AI? Log warning.
			logrus.Debugf("Warning: Entity %d in turn queue is not player or AI.", actorEntry.EntityID)
			// Don't re-queue for now.
		}
		// Loop continues to process next actor in the queue at the current time
	}

	// Clear active entity at the end
	logrus.Debug("========= processTurnQueue ended (iteration limit reached) =========")
}
