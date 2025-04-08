package game

import (
	"fmt"

	"codeberg.org/anaseto/gruid"

	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs/components"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ui"
	"github.com/sirupsen/logrus"
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

// handleEntityDeath handles an entity's death, either removing it completely
// or turning it into a corpse (the preferred option)
func (g *Game) handleEntityDeath(entityID ecs.EntityID, entityName string) {
	// Log the death event
	g.Log.AddMessagef(ui.ColorDeath, "%s dies!", entityName)
	logrus.Infof("Entity %s (%d) has died.", entityName, entityID)

	// Check if this is the player
	if entityID == g.PlayerID {
		// Handle player death specially (could trigger game over screen later)
		g.Log.AddMessagef(ui.ColorCritical, "You died! Game over!")
		logrus.Info("Player has died. Game over!")
		// TODO: Implement game over state
		return
	}

	// Turn entity into a corpse
	// 1. Remove components that make it active or interactive
	g.ecs.RemoveComponents(entityID,
		components.CTurnActor,      // Remove from turn system
		components.CAITag,          // Remove AI behavior
		components.CBlocksMovement, // Allow walking over corpses
	)

	// 2. Change appearance to a corpse
	g.ecs.AddComponents(entityID,
		components.Renderable{Glyph: '%', Color: ui.ColorCorpse},
		components.CorpseTag{},
		components.DeadTag{},
	)

	// 3. Remove from turn queue
	// The queue's Remove method handles finding and removing by entity ID
	g.turnQueue.Remove(entityID)

	// 4. Optional: add decomposition logic or timer for corpse cleanup
	// TODO: Implement corpse decay system if desired
}
