// This file defines the main model of the game: the Update function that
// updates the model state in response to user input, and the Draw function,
// which draws the grid.

package main

import (
	"runtime"

	"codeberg.org/anaseto/gruid"
)

type mode int

const (
	modeNormal mode = iota
	modeQuit
)

// model represents our main application's state.
type model struct {
	grid gruid.Grid // drawing grid
	game *game      // game state
	mode mode       // current mode
}

func (md *model) init() gruid.Effect {
	md.game.InitLevel()

	if runtime.GOOS == "js" {
		return nil
	}

	return gruid.Sub(subSig)
}

// EndTurn finalizes player's turn and runs other events until next player
// turn.
func (md *model) EndTurn() gruid.Effect {
	md.mode = modeNormal
	return nil
}
