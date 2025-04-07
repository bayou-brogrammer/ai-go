package game

import (
	"codeberg.org/anaseto/gruid"
	"fmt"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
	"github.com/sirupsen/logrus" // Added for logging
)

// GameAction is an interface for actions that can be performed in the game.
type GameAction interface {
	Execute(g *Game) (cost uint, err error)
}

type MoveAction struct {
	Direction gruid.Point
	EntityID  ecs.EntityID
}

// Execute performs the move action, returning the time cost and any error.
func (a MoveAction) Execute(g *Game) (cost uint, err error) {
	// Use the EntityID from the action
	again, err := g.EntityBump(a.EntityID, a.Direction)
	if err != nil {
		return 0, err // No cost if error occurred
	}

	if !again {
		// Bumped into something, action didn't fully succeed in moving
		return 0, nil // No time cost for a bump
	}

	// Successful move
	return 100, nil // Standard move cost
}

// AttackAction represents an entity attacking another entity.
type AttackAction struct {
	AttackerID ecs.EntityID
	TargetID   ecs.EntityID
}

// Execute performs the attack action.
func (a AttackAction) Execute(g *Game) (cost uint, err error) {
	attackerName, _ := g.ecs.GetName(a.AttackerID)
	targetName, _ := g.ecs.GetName(a.TargetID)
	targetHealth, ok := g.ecs.GetHealth(a.TargetID)
	if !ok {
		// Target might have died between action queuing and execution
		logrus.Debugf("%s (%d) tries to attack %s (%d), but target has no health component.", attackerName, a.AttackerID, targetName, a.TargetID)
		return 0, fmt.Errorf("target %d has no health", a.TargetID)
	}

	// --- Basic Damage Calculation (Placeholder) ---
	damage := 1 // Simple fixed damage for now
	targetHealth.CurrentHP -= damage

	logrus.Infof("%s (%d) attacks %s (%d) for %d damage. %s HP: %d/%d",
		attackerName, a.AttackerID,
		targetName, a.TargetID,
		damage,
		targetName, targetHealth.CurrentHP, targetHealth.MaxHP)

	// Update the target's health component in the ECS
	g.ecs.AddHealth(a.TargetID, targetHealth) // AddHealth also works for updates

	// TODO: Check for death (CurrentHP <= 0) and handle it (e.g., remove entity, drop items)

	return 100, nil // Standard attack cost
}
