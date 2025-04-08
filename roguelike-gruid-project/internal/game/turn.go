package game

import (
	"fmt"

	"codeberg.org/anaseto/gruid"

	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs/components"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ui"
	"github.com/sirupsen/logrus" // Added for logging
)

// GameAction is an interface for actions that can be performed in the game.
type GameAction interface {
	Execute(g *Game) (cost uint, err error)
}

type WaitAction struct {
	EntityID ecs.EntityID
}

func (a WaitAction) Execute(g *Game) (cost uint, err error) {
	return 100, nil // Standard wait cost
}

type MoveAction struct {
	Direction gruid.Point
	EntityID  ecs.EntityID
}

// Execute performs the move action, returning the time cost and any error.
func (a MoveAction) Execute(g *Game) (cost uint, err error) {
	again, err := g.EntityBump(a.EntityID, a.Direction)
	if err != nil {
		return 0, err // No cost if error occurred
	}

	if !again {
		// Bumped into something, action didn't fully succeed in moving
		return 0, nil // No time cost for a bump
	}
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

	// --- Basic Damage Calculation ---
	damage := 1 // Simple fixed damage for now
	targetHealth.CurrentHP -= damage

	// Determine message color based on who is attacking
	var msgColor gruid.Color
	if a.AttackerID == g.PlayerID {
		msgColor = ui.ColorPlayerAttack // Define in ui/color.go
	} else if a.TargetID == g.PlayerID {
		msgColor = ui.ColorEnemyAttack // Define in ui/color.go
	} else {
		msgColor = ui.ColorNeutralAttack // Define in ui/color.go
	}
	g.Log.AddMessagef(msgColor, "%s attacks %s for %d damage.", attackerName, targetName, damage)

	logrus.Infof("%s (%d) attacks %s (%d) for %d damage. %s HP: %d/%d",
		attackerName, a.AttackerID,
		targetName, a.TargetID,
		damage,
		targetName, targetHealth.CurrentHP, targetHealth.MaxHP)
	g.ecs.AddComponent(a.TargetID, components.CHealth, targetHealth)

	// Check for death (CurrentHP <= 0) and handle it
	if targetHealth.CurrentHP <= 0 {
		g.handleEntityDeath(a.TargetID, targetName)
	}

	return 100, nil // Standard attack cost
}

// Add this new method to Game to handle entity death
func (g *Game) handleEntityDeath(entityID ecs.EntityID, entityName string) {
	g.Log.AddMessagef(ui.ColorDeath, "%s dies!", entityName)
	logrus.Infof("Entity %s (%d) has died.", entityName, entityID)

	// Option 1: Remove entity completely
	// g.ecs.RemoveEntity(entityID)

	// Option 2: Turn into a corpse (preferred as per your feedback)
	// Remove components that make it active
	g.ecs.RemoveComponents(entityID, components.CTurnActor, components.CAITag)

	// Change appearance to a corpse
	g.ecs.AddComponents(entityID, components.CRenderable, components.Renderable{
		Glyph: '%',
		Color: ui.ColorCorpse,
	}, components.CorpseTag{})

	// Remove from turn queue if it's there
	// This depends on how your turn queue is structured
	// g.turnQueue.Remove(entityID) // If your queue allows removal by entity ID
}
