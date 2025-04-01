// This file defines the main model of the game: the Update function that
// updates the model state in response to user input, and the Draw function,
// which draws the grid.

package main

import (
	"runtime"

	"codeberg.org/anaseto/gruid"
)

// model represents our main application's state.
type model struct {
	grid gruid.Grid // drawing grid
	game *game      // game state
}

func (m *model) init() gruid.Effect {
	if runtime.GOOS == "js" {
		return nil
	}
	return gruid.Sub(subSig)

}

func (md *model) Update(msg gruid.Msg) gruid.Effect {
	if _, ok := msg.(gruid.MsgInit); ok {
		return md.init()
	}

	switch msg := msg.(type) {
	case gruid.MsgInit:
		_ = msg
		return nil
	case gruid.MsgQuit:
		return gruid.End()
	case gruid.MsgKeyDown:
		if msg.Key == "q" {
			return gruid.End()
		}
		return nil

	}

	return nil
}

// Draw implements gruid.Model.Draw. It clears the grid, then renders the map
// and all entities using the RenderSystem.
func (md *model) Draw() gruid.Grid {
	// Clear the grid before drawing
	md.grid.Fill(gruid.Cell{Rune: ' '}) // Fill with blank spaces

	// Draw the map tiles first
	if md.game != nil && md.game.Map != nil {
		DrawMap(md.game.Map, md.grid)
	}

	// Render entities using the ECS RenderSystem
	if md.game != nil && md.game.ecs != nil {
		RenderSystem(md.game.ecs, md.grid)
	}

	return md.grid
}
