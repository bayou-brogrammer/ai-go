package game

import (
	"math/rand"
	"time"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/turn"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ui"
	"github.com/sirupsen/logrus"

	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs/components"
)

// Game settings & map generation constants
const (
	Width              = 80 // Default width
	Height             = 24 // Default height
	maxRooms           = 10 // Max rooms on a level
	roomMinSize        = 6  // Min width/height of a room
	roomMaxSize        = 10 // Max width/height of a room
	maxMonstersPerRoom = 2  // Max monsters per room (excluding first)
	// fovRadius = 8 // How far the player can see - Keep commented for now
)

// Removed redundant const block

// Game represents the main game state.
type Game struct {
	Depth int

	Map       *Map
	ecs       *ecs.ECS        // The Entity-Component-System manager
	PlayerID  ecs.EntityID    // Store the player's entity ID
	turnQueue *turn.TurnQueue // Event queue for game events

	rand *rand.Rand
}

func NewGame() *Game {
	return &Game{
		ecs:       ecs.NewECS(),
		turnQueue: turn.NewTurnQueue(),
	}
}

func (g *Game) InitLevel() {
	if g.rand == nil {
		g.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	g.Depth = 1

	g.Map = NewMap(Width, Height)
	playerStart := g.Map.generateMap(g, Width, Height) // Pass game struct
	g.InitPlayer(playerStart)
}

func (g *Game) InitPlayer(playerStart gruid.Point) {
	// Create the player entity
	playerID := g.ecs.AddEntity(components.PlayerTag{})
	g.PlayerID = playerID // Store the player ID in the game struct

	logrus.Debugf("Player ID: %d\n", playerID)

	// Add components to the player
	// Use the start position returned by NewMap
	g.ecs.AddPosition(playerID, playerStart)
	g.ecs.AddRenderable(playerID, components.Renderable{Glyph: '@', Color: ui.ColorPlayer})

	// Add the player to the turn queue
	g.turnQueue.Add(playerID, g.turnQueue.CurrentTime)

}
