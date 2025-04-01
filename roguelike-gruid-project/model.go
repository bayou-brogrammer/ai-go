// This file defines the main model of the game: the Update function that
// updates the model state in response to user input, and the Draw function,
// which draws the grid.

package main

import (
	"codeberg.org/anaseto/gruid"
)

// model represents our main application's state.
type model struct {
	grid gruid.Grid // drawing grid
	game *game      // game state
}

func (m *model) Update(msg gruid.Msg) gruid.Effect {
	switch msg := msg.(type) {
	case gruid.MsgInit:
		_ = msg
		return nil
	case gruid.MsgQuit:
		return gruid.End()
	}

	return nil
}

// Draw implements gruid.Model.Draw. It clears the grid, then renders the map
// and all entities using the RenderSystem.
func (m *model) Draw() gruid.Grid {
	// Clear the grid before drawing
	m.grid.Fill(gruid.Cell{Rune: ' '}) // Fill with blank spaces

	// TODO: Draw the map tiles first (we'll add this later)
	// Example: DrawMap(m.game.Map, m.grid)

	// Render entities using the ECS RenderSystem
	if m.game != nil && m.game.ecs != nil {
		RenderSystem(m.game.ecs, m.grid)
	}

	return m.grid
}
