package game

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs/components"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/utils"
	"github.com/sirupsen/logrus"
)

// RenderSystem draws all entities with Position and Renderable components onto the grid.
// RenderSystem draws entities with Position and Renderable components onto the grid,
// respecting the player's field of view. It now requires the player's FOV component
// and the map width to perform visibility checks.
func RenderSystem(ecs *ecs.ECS, grid gruid.Grid, playerFOV *components.FOV, mapWidth int) bool {
	utils.Assert(ecs != nil, "ECS is nil")
	utils.Assert(playerFOV != nil, "Player FOV is nil")
	utils.Assert(mapWidth > 0, "Map width is not positive")

	needsRedraw := false
	entityIDs := ecs.GetEntitiesWithComponents(components.CPosition, components.CRenderable)
	for _, id := range entityIDs {
		pos, _ := ecs.GetPosition(id)
		if playerFOV.IsVisible(pos, mapWidth) {
			entityNeedsRedraw := drawEntity(ecs, pos, id, grid)
			if entityNeedsRedraw {
				fmt.Println("Entity needs redraw")
				needsRedraw = true
			}
		}
	}

	return needsRedraw
}

// Define the passable function once (reusing Map's IsOpaque)
func (g *Game) passable(p gruid.Point) bool {
	return !g.Map.IsOpaque(p)
}

// FOVSystem updates the visibility for all entities with an FOV component.
func FOVSystem(g *Game) {
	entities := g.ecs.GetEntitiesWithPositionAndFOV()

	for _, entity := range entities {
		id, pos, fov := entity.ID, entity.Position, entity.FOV
		fov.ClearVisible()
		fovCalculator := fov.GetFOVCalculator()
		visibleTiles := fovCalculator.SSCVisionMap(pos, fov.Range, g.passable, false) // Use asserted pos
		visibleTiles = utils.DrawFilledCircle(visibleTiles, fov.Range, pos)           // Use asserted pos

		// Update visibility and explored status for points in the circle
		for _, p := range visibleTiles {
			fov.SetVisible(p, g.Map.Width)

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
	for i := range 100 { // Limit iterations to prevent infinite loops
		logrus.Debugf("Turn queue iteration %d", i)

		if g.turnQueue.IsEmpty() {
			logrus.Debug("Turn queue is empty.")
			logrus.Debug("========= processTurnQueue ended (queue empty) =========")
			return
		}

		turnEntry, ok := g.turnQueue.Next()
		if !ok {
			logrus.Debug("Error: Queue does not have any more actors.")
			logrus.Debug("========= processTurnQueue ended (Next error) =========")
			return
		}

		logrus.Debugf("Processing actor: EntityID=%d, Time=%d", turnEntry.EntityID, turnEntry.Time)
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
		if isPlayer && action == nil {
			g.waitingForInput = true
			logrus.Debug("It's the player's turn, waiting for input.")
			g.turnQueue.Add(turnEntry.EntityID, turnEntry.Time)
			logrus.Debug("========= processTurnQueue ended (player's turn) =========")
			return
		}
		if action == nil {
			logrus.Debugf("Entity %d has no actions, rescheduling turn at time %d", turnEntry.EntityID, turnEntry.Time)
			g.turnQueue.Add(turnEntry.EntityID, turnEntry.Time)
			continue
		}
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
