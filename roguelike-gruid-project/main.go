package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	// Import for random placement if needed later
	"codeberg.org/anaseto/gruid"
	// "codeberg.org/anaseto/gruid/fov" // Removed - Incorrect path
)

const (
	gameWidth  = 80
	gameHeight = 24
	// fovRadius  = 8 // Commented out for now
)

func main() {
	// Seed random number generator (optional, for more complex spawning later)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize the game state and ECS World
	g := &game{
		Map: NewMap(gameWidth, gameHeight),
		ecs: NewWorld(),
	}

	// Create the player entity
	playerID := g.ecs.CreateEntity()

	// Add components to the player
	startPos := gruid.Point{X: gameWidth / 2, Y: gameHeight / 2}
	g.ecs.AddComponent(playerID, Position{Point: startPos})
	g.ecs.AddComponent(playerID, Renderable{Glyph: '@', Color: gruid.ColorDefault})
	g.ecs.AddComponent(playerID, BlocksMovement{}) // Player blocks movement

	// Create a new grid with standard 80x24 size.
	gd := gruid.NewGrid(gameWidth, gameHeight)

	// Create the main application's model, using grid gd.
	m := &model{grid: gd, game: g}

	app := gruid.NewApp(gruid.AppConfig{
		Driver: driver,
		Model:  m,
	})

	// start application
	if err := app.Start(context.Background()); err != nil {
		driver.Close()
		log.Fatal(err)
	}
}
