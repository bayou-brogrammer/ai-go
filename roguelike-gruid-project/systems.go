package main

import (
	"reflect"

	"codeberg.org/anaseto/gruid"
)

// RenderSystem draws all entities with Position and Renderable components onto the grid.
func RenderSystem(ecs *World, grid gruid.Grid) {
	// Get component types needed for querying
	posType := reflect.TypeOf((*Position)(nil)).Elem()
	renType := reflect.TypeOf((*Renderable)(nil)).Elem()

	// Query for entities that have both Position and Renderable
	entityIDs := ecs.QueryEntitiesWithComponents(posType, renType)

	// Iterate through entities and draw them
	for _, id := range entityIDs {
		// Retrieve components (we know they exist from the query)
		posComp, _ := ecs.GetComponent(id, posType)
		renComp, _ := ecs.GetComponent(id, renType)

		// Type assert components to access their fields
		pos := posComp.(Position)
		ren := renComp.(Renderable)

		// Set the cell in the grid
		// Note: This assumes the grid has been cleared beforehand.
		// We might add map bounds checking later.
		grid.Set(pos.Point, gruid.Cell{
			Rune:  ren.Glyph,
			Style: gruid.Style{Fg: ren.Color}, // Assuming default background for now
		})
	}
}

// Note: We will add other systems like MovementSystem, AISystem, etc., here later.
