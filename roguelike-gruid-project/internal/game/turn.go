package game

import (
	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
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
