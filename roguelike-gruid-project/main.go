//go:build !js
// +build !js

package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
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
	rand.New(rand.NewSource(time.Now().UnixNano())) // Correctly seed the global rand

	// Initialize the game state and ECS World
	gameMap, playerStart := NewMap(gameWidth, gameHeight) // Capture both return values
	g := NewGame(gameMap)

	// Create the player entity
	playerID := g.ecs.CreateEntity()

	// Add components to the player
	// Use the start position returned by NewMap
	g.ecs.AddComponent(playerID, Position{Point: playerStart})
	g.ecs.AddComponent(playerID, Renderable{Glyph: '@', Color: gruid.ColorDefault})
	g.ecs.AddComponent(playerID, BlocksMovement{}) // Player blocks movement

	g.PlayerID = playerID // Store the player ID in the game struct

	gd := gruid.NewGrid(gameWidth, gameHeight)
	m := &model{grid: gd, game: &game{}}

	app := gruid.NewApp(gruid.AppConfig{
		Model:  m,
		Driver: driver,
	})

	// start application
	if err := app.Start(context.Background()); err != nil {
		driver.Close()
		log.Fatal(err)
	}
}

func subSig(ctx context.Context, msgs chan<- gruid.Msg) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sig)
	select {
	case <-ctx.Done():
	case <-sig:
		msgs <- gruid.MsgQuit{}
	}
}
