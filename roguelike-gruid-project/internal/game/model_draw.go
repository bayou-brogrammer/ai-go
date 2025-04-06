package game

import (
	"codeberg.org/anaseto/gruid"                                      // Needed for FOV type
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ui" // For colors
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/utils"
)

// Draw implements gruid.Model.Draw. It clears the grid, then renders the map
// and all entities using the RenderSystem.
func (md *Model) Draw() gruid.Grid {
	g := md.game

	utils.Assert(g != nil, "Game is nil")
	utils.Assert(g.ecs != nil, "ECS is nil")
	utils.Assert(g.Map != nil, "Map is nil")

	// Update FOV for all entities before drawing
	FOVSystem(g)

	// Get player's FOV component *after* FOVSystem runs
	playerFOVComp, _ := g.ecs.GetFOV(g.PlayerID)

	// Clear the grid before drawing
	md.grid.Fill(gruid.Cell{Rune: ' '}) // Fill with blank spaces

	// Draw the map tiles based on explored and visible status
	it := g.Map.Grid.Iterator()
	for it.Next() {
		p := it.P()

		isExplored := g.Map.IsExplored(p)
		if !isExplored {
			continue // Don't draw unexplored tiles
		}

		isVisible := playerFOVComp.IsVisible(p, g.Map.Width)

		if !isVisible {
			// Explored but not visible: Draw in a faded color
			md.grid.Set(p, gruid.Cell{
				Rune:  g.Map.Rune(it.Cell()),
				Style: gruid.Style{Fg: ui.ColorExploredWall, Bg: ui.ColorExploredFloor}, // Use faded colors
			})
		} else {
			// Currently visible: Draw normally
			md.grid.Set(p, gruid.Cell{
				Rune:  g.Map.Rune(it.Cell()),
				Style: gruid.Style{Fg: ui.ColorVisibleWall, Bg: ui.ColorVisibleFloor}, // Use bright colors
			})
		}
	}

	// Render entities using the ECS RenderSystem, passing player FOV if available
	RenderSystem(g.ecs, md.grid, playerFOVComp, g.Map.Width)

	return md.grid
}
