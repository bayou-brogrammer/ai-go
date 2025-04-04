// This file defines the main model of the game: the Update function that
// updates the model state in response to user input, and the Draw function,
// which draws the grid.

package game

import (
	"runtime"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/utils"
)

type mode int

const (
	modeNormal mode = iota
	modeQuit
)

// Model represents the game model that implements gruid.Model
type Model struct {
	grid gruid.Grid
	game *Game
	mode mode
}

// NewModel creates a new game model
func NewModel(grid gruid.Grid) *Model {
	return &Model{
		grid: grid,
		game: NewGame(),
	}
}

func (md *Model) init() gruid.Effect {
	md.game.InitLevel()

	if runtime.GOOS == "js" {
		return nil
	}

	return gruid.Sub(utils.HandleSignals)
}

// EndTurn finalizes player's turn and runs other events until next player
// turn.
func (md *Model) EndTurn() gruid.Effect {
	md.mode = modeNormal
	return nil
}
