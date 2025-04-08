package game

import (
	"math/rand"
	"time"

	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/config"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/log"
	turn "github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/turn_queue"
)
// Game represents the main game state.
type Game struct {
	Depth           int
	waitingForInput bool

	Map       *Map
	ecs       *ecs.ECS        // The Entity-Component-System manager
	PlayerID  ecs.EntityID    // Store the player's entity ID
	turnQueue *turn.TurnQueue // Event queue for game events
	Log       *log.MessageLog // Add log field for game messages

	rand *rand.Rand
}

func NewGame() *Game {
	return &Game{
		ecs:       ecs.NewECS(),
		turnQueue: turn.NewTurnQueue(),
		Log:       log.NewMessageLog(), // Initialize message log
	}
}

func (g *Game) InitLevel() {
	if g.rand == nil {
		g.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	g.Depth = 1

	g.Map = NewMap(config.DungeonWidth, config.DungeonHeight)
	playerStart := g.Map.generateMap(g, config.DungeonWidth, config.DungeonHeight) // Pass game struct
	g.SpawnPlayer(playerStart)
}
