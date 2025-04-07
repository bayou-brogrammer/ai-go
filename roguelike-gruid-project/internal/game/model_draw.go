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
		isWall := g.Map.IsWall(p)

		// Use the new helper function to get the appropriate style
		style := ui.GetMapStyle(isWall, isVisible, isExplored)

		md.grid.Set(p, gruid.Cell{
			Rune:  g.Map.Rune(it.Cell()),
			Style: style,
		})
	}

	// Render entities using the ECS RenderSystem, passing player FOV if available
	RenderSystem(g.ecs, md.grid, playerFOVComp, g.Map.Width)

	return md.grid
}
