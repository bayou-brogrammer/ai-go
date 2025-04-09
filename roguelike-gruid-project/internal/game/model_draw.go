package game

import (
	"codeberg.org/anaseto/gruid" // Needed for FOV type
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ui" // For colors
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/utils"
	"github.com/sirupsen/logrus"
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
	playerFOVComp, ok := g.ecs.GetFOV(g.PlayerID)
	if !ok {
		// Handle case where player FOV might be missing (though unlikely)
		logrus.Errorf("Player entity %d missing FOV component in Draw", g.PlayerID)
		return md.grid // Return the grid even if FOV is missing
	}

	// Clear the grid before drawing
	md.grid.Fill(gruid.Cell{Rune: ' '})

	it := g.Map.Grid.Iterator()
	for it.Next() {
		p := it.P()

		isExplored := g.Map.IsExplored(p)
		if !isExplored {
			continue
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

// When drawing an entity, check for HitFlash
func drawEntity(ecs *ecs.ECS, pos gruid.Point, entityID ecs.EntityID, grid gruid.Grid) bool {
	renderable, ok := ecs.GetRenderable(entityID)
	if !ok {
		return false
	}

	needsRedraw := false
	color := renderable.Color

	// Draw the entity with the appropriate color
	grid.Set(pos, gruid.Cell{Rune: renderable.Glyph, Style: gruid.Style{Fg: color}})
	return needsRedraw
}
