package game

import (
	"fmt"
	"math/rand"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs/components" // Added for FOV type
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/utils"
	"github.com/sirupsen/logrus"
)

// RenderSystem draws all entities with Position and Renderable components onto the grid.
// RenderSystem draws entities with Position and Renderable components onto the grid,
// respecting the player's field of view. It now requires the player's FOV component
// and the map width to perform visibility checks.
func RenderSystem(ecs *ecs.ECS, grid gruid.Grid, playerFOV *components.FOV, mapWidth int) {
	utils.Assert(ecs != nil, "ECS is nil")
	utils.Assert(playerFOV != nil, "Player FOV is nil")
	utils.Assert(mapWidth > 0, "Map width is not positive")

	entityIDs := ecs.GetEntitiesWithPositionAndRenderable()

	// Iterate through entities and draw them
	for _, id := range entityIDs {
		// Retrieve components (we know they exist from the query)
		pos, _ := ecs.GetPosition(id)
		ren, renOk := ecs.GetRenderable(id)
		if !renOk {
			continue
		} // Should not happen based on query, but safe check

		// Only draw the entity if its position is visible to the player
		if playerFOV.IsVisible(pos, mapWidth) {
			// Set the cell in the grid
			grid.Set(pos, gruid.Cell{
				Rune:  ren.Glyph,
				Style: gruid.Style{Fg: ren.Color}, // Use chosen color
			})
		}
	}
}

// FOVSystem updates the visibility for all entities with an FOV component.
func FOVSystem(g *Game) {
	// Query entities with Position and FOV components
	entities := g.ecs.GetEntitiesWithPositionAndFOV()

	// Define the passable function once (reusing Map's IsOpaque)
	passable := func(p gruid.Point) bool {
		return !g.Map.IsOpaque(p)
	}

	for _, id := range entities {
		posCompData, posOk := g.ecs.GetPosition(id) // Rename variable to avoid conflict
		fovComp, fovOk := g.ecs.GetFOV(id)

		if !posOk || !fovOk {
			logrus.Warnf("Entity %d missing Position or FOV component in FOVSystem", id)
			continue
		}

		// Clear the entity's current visibility
		fovComp.ClearVisible()

		// Calculate new visibility using the component's calculator
		fovCalculator := fovComp.GetFOVCalculator()
		visibleTiles := fovCalculator.SSCVisionMap(posCompData, fovComp.Range, passable, false)
		visibleTiles = utils.DrawFilledCircle(visibleTiles, fovComp.Range, posCompData)

		// Update visibility and explored status for points in the circle
		for _, p := range visibleTiles {
			fovComp.SetVisible(p, g.Map.Width)

			// If this is the player, also update the global explored map
			if id == g.PlayerID {
				g.Map.SetExplored(p)
			}
		}
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

		moveOrWait := rand.Intn(2)
		if moveOrWait == 0 {
			// Move
			action, err := moveMonster(g, id)
			if err != nil {
				logrus.Debugf("Failed to move monster %d: %v", id, err)
				continue
			}
			actor.AddAction(action)
		} else {
			actor.AddAction(WaitAction{EntityID: id})
		}
	}
}

func moveMonster(g *Game, id ecs.EntityID) (GameAction, error) {
	// Get the entity's current position
	pos, ok := g.ecs.GetPosition(id)
	if !ok {
		return nil, fmt.Errorf("entity %d has no position", id)
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
		return action, nil
	} else {
		// If no valid direction was found, just wait
		logrus.Debugf("AI entity %d has no valid move, waiting", id)
		return WaitAction{EntityID: id}, nil
	}
}
