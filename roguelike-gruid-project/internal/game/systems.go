package game

import (
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

// processTurnQueue processes turns for actors until it's the player's turn
// or the queue is exhausted for the current time step.
func (md *Model) processTurnQueue() {
	g := md.game
	logrus.Debug("========= processTurnQueue started =========")

	// Periodically clean up the queue
	metrics := g.turnQueue.CleanupDeadEntities(g.ecs)
	if metrics.EntitiesRemoved > 10 {
		logrus.Infof(
			"Turn queue cleanup: removed %d entities in %v",
			metrics.EntitiesRemoved,
			metrics.ProcessingTime,
		)
	}

	g.turnQueue.PrintQueue()

	// Process turns until we need player input or run out of actors
	for i := 0; i < 100; i++ { // Limit iterations to prevent infinite loops
		logrus.Debugf("Turn queue iteration %d", i)

		// Check if the queue is empty
		if g.turnQueue.IsEmpty() {
			logrus.Debug("Turn queue is empty.")
			logrus.Debug("========= processTurnQueue ended (queue empty) =========")
			return
		}

		// Get the next actor's turn
		turnEntry, ok := g.turnQueue.Next()
		if !ok {
			logrus.Debug("Error: Queue does not have any more actors.")
			logrus.Debug("========= processTurnQueue ended (Next error) =========")
			return
		}

		logrus.Debugf("Processing actor: EntityID=%d, Time=%d", turnEntry.EntityID, turnEntry.Time)

		// Get the actor's TurnActor component
		actor, ok := g.ecs.GetTurnActor(turnEntry.EntityID)
		if !ok {
			logrus.Debugf("Error: Entity %d is not a valid actor.", turnEntry.EntityID)
			continue
		}

		if !actor.IsAlive() {
			logrus.Debugf("Entity %d is not alive, skipping turn.", turnEntry.EntityID)
			continue
		}

		isPlayer := turnEntry.EntityID == g.PlayerID
		action := actor.NextAction()

		// If it's the player's turn and they have no action queued
		if isPlayer && action == nil {
			g.waitingForInput = true
			logrus.Debug("It's the player's turn, waiting for input.")
			g.turnQueue.Add(turnEntry.EntityID, turnEntry.Time)
			logrus.Debug("========= processTurnQueue ended (player's turn) =========")
			return
		}

		// If no action is available, reschedule the turn
		if action == nil {
			logrus.Debugf("Entity %d has no actions, rescheduling turn at time %d", turnEntry.EntityID, turnEntry.Time)
			g.turnQueue.AddToBack(turnEntry.EntityID, turnEntry.Time)
			continue
		}

		// Execute the action
		cost, err := action.(GameAction).Execute(g)
		if err != nil {
			logrus.Debugf("Failed to execute action for entity %d: %v", turnEntry.EntityID, err)

			// On failure, reschedule with appropriate delay
			if isPlayer {
				g.turnQueue.Add(turnEntry.EntityID, turnEntry.Time)
			} else {
				g.turnQueue.Add(turnEntry.EntityID, turnEntry.Time+100)
			}
			continue
		}

		logrus.Debugf("Action executed for entity %d, cost: %d", turnEntry.EntityID, cost)

		// Update the game time and schedule next turn
		g.turnQueue.CurrentTime = turnEntry.Time + uint64(cost)
		g.turnQueue.Add(turnEntry.EntityID, g.turnQueue.CurrentTime)
	}

	logrus.Debug("========= processTurnQueue ended (iteration limit reached) =========")
}

// monstersTurn handles AI turns for all monsters in the game.
func (g *Game) monstersTurn() {
	// Get all entities with AI tag and TurnActor component
	for id := range g.ecs.AITags {

		// Skip if entity doesn't have a TurnActor
		actor, ok := g.ecs.GetTurnActor(id)
		if !ok {
			continue
		}

		// Skip if entity is dead
		if !actor.IsAlive() {
			continue
		}

		// Skip entities that already have actions queued
		if actor.PeekNextAction() != nil {
			continue
		}

		// Get the entity's current position
		pos, ok := g.ecs.GetPosition(id)
		if !ok {
			continue
		}

		// Try different directions in a random order
		directions := []gruid.Point{
			{X: -1, Y: 0}, // West
			{X: 1, Y: 0},  // East
			{X: 0, Y: -1}, // North
			{X: 0, Y: 1},  // South
		}

		// Shuffle the directions to add randomness
		// This is a simple way to randomize the order of directions
		rand.Shuffle(len(directions), func(i, j int) {
			directions[i], directions[j] = directions[j], directions[i]
		})

		// Find a valid direction to move (one that leads to a walkable tile)
		var validMove *gruid.Point
		for _, dir := range directions {
			newPos := pos.Add(dir)

			// Check if the new position is valid (walkable and not occupied)
			if g.Map.isWalkable(newPos) && len(g.ecs.EntitiesAt(newPos)) == 0 {
				validMove = &dir
				break
			}
		}

		// If we found a valid direction, queue the walk action
		if validMove != nil {
			logrus.Debugf("AI entity %d moving in direction %v", id, validMove)
			action := MoveAction{
				Direction: *validMove,
				EntityID:  id,
			}
			actor.AddAction(action)
		} else {
			// If no valid direction was found, just wait
			logrus.Debugf("AI entity %d has no valid move, waiting", id)
			// TODO: Add a wait action if needed
		}
	}
}
