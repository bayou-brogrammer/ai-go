package main

import (
	"math/rand"
	"time"

	"codeberg.org/anaseto/gruid"
)

const (
// fovRadius = 8 // How far the player can see
)

// Game represents the main game state.
type game struct {
	Depth int

	Map       *Map
	ecs       *ECS       // The Entity-Component-System manager
	PlayerID  EntityID   // Store the player's entity ID
	turnQueue *TurnQueue // Event queue for game events

	rand *rand.Rand
}

func NewGame() *game {
	return &game{
		ecs:       NewECS(),
		turnQueue: NewTurnQueue(),
	}
}

func (g *game) InitLevel() {
	if g.rand == nil {
		g.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	g.Depth = 1

	g.Map = NewMap(gameWidth, gameHeight)
	playerStart := g.Map.generateMap(gameWidth, gameHeight)
	g.InitPlayer(playerStart)
}

func (g *game) InitPlayer(playerStart gruid.Point) {
	// Create the player entity
	playerID := g.ecs.AddEntity(Player{})
	g.PlayerID = playerID // Store the player ID in the game struct

	// Add components to the player
	// Use the start position returned by NewMap
	g.ecs.AddPosition(playerID, playerStart)
	g.ecs.AddRenderable(playerID, Renderable{Glyph: '@', Color: gruid.ColorDefault})

	// Add the player to the turn queue
	g.turnQueue.Add(playerID, g.turnQueue.CurrentTime)

}
