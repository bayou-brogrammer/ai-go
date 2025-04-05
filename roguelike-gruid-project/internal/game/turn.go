package game

import (
	"codeberg.org/anaseto/gruid"
)

// GameAction is an interface for actions that can be performed in the game.
type GameAction interface {
	Execute(g *Game) (cost uint, err error) // Modified return signature
}

type MoveAction struct {
	Direction gruid.Point
}

// Execute performs the move action, returning the time cost and any error.
func (a MoveAction) Execute(g *Game) (cost uint, err error) {
	// TODO: This currently assumes the action is for the player.
	// Need to generalize for any entity ID.
	// For now, we'll pass g.PlayerID, but this needs refactoring later.
	again, err := g.EntityBump(g.PlayerID, a.Direction) // Assuming EntityBump exists or will be created
	if err != nil {
		return 0, err // No cost if error occurred
	}

	if again {
		// Bumped into something, action didn't fully succeed in moving
		return 0, nil // No time cost for a bump? Or maybe a smaller cost? Let's use 0 for now.
	}

	// Successful move
	return 100, nil // Standard move cost
}
