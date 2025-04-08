package game

import (
	"fmt"
	"time"

	"codeberg.org/anaseto/gruid" // Needed for FOV type
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs/components"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ui" // For colors
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/utils"
	"github.com/sirupsen/logrus"
)

// Draw implements gruid.Model.Draw. It clears the grid, then renders the map
// and all entities using the RenderSystem.
func (md *Model) Draw() gruid.Grid {
	fmt.Println("Draw called")
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
	needsRedraw := RenderSystem(g.ecs, md.grid, playerFOVComp, g.Map.Width)
	if needsRedraw {
		g.ClearHitFlashes()

		c1 := make(chan string, 1)
		go func() {
			fmt.Println("SLEEPING")
			// Clear hit flashes after rendering but before returning the grid
			time.Sleep(300 * time.Millisecond)
			fmt.Println("DRAWING")
			md.Draw()
			c1 <- "done"
		}()
	}

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

	// Check if entity has hit flash using the generic system
	if ecs.HasComponent(entityID, components.CHitFlash) {
		// Override color with hit flash color - define this in your color constants
		color = ui.ColorHitFlash // Typically white or bright red
		needsRedraw = true
	}

	// Draw the entity with the appropriate color
	grid.Set(pos, gruid.Cell{Rune: renderable.Glyph, Style: gruid.Style{Fg: color}})
	return needsRedraw
}

// ClearHitFlashes removes all HitFlash components after they've been rendered
func (g *Game) ClearHitFlashes() {
	// Get all entities with HitFlash component using the generic system
	entities := g.ecs.GetEntitiesWithComponent(components.CHitFlash)

	// Remove the HitFlash component from each entity
	for _, entityID := range entities {
		g.ecs.RemoveComponent(entityID, components.CHitFlash)
	}
}
