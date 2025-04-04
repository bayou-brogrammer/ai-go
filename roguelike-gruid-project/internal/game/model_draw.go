package game

import "codeberg.org/anaseto/gruid"

// Draw implements gruid.Model.Draw. It clears the grid, then renders the map
// and all entities using the RenderSystem.
func (md *Model) Draw() gruid.Grid {
	g := md.game

	// Clear the grid before drawing
	md.grid.Fill(gruid.Cell{Rune: ' '}) // Fill with blank spaces

	// Draw the map tiles first
	if g != nil && g.Map != nil {
		// We draw the map tiles.
		it := g.Map.Grid.Iterator()
		for it.Next() {
			// if !g.Map.Explored[it.P()] {
			// 	continue
			// }

			c := gruid.Cell{Rune: g.Map.Rune(it.Cell())}
			// if g.InFOV(it.P()) {
			// 	c.Style.Bg = ColorFOV
			// }
			md.grid.Set(it.P(), c)
		}
	}

	// Render entities using the ECS RenderSystem
	if md.game != nil && md.game.ecs != nil {
		RenderSystem(md.game.ecs, md.grid)
	}

	return md.grid
}
