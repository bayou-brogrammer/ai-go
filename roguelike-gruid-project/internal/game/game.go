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

	dungeon   *Map
	ecs       *ecs.ECS
	PlayerID  ecs.EntityID
	turnQueue *turn.TurnQueue
	log       *log.MessageLog

	rand *rand.Rand
}

func NewGame() *Game {
	return &Game{
		ecs:       ecs.NewECS(),
		turnQueue: turn.NewTurnQueue(),
		log:       log.NewMessageLog(),
	}
}

func (g *Game) InitLevel() {
	if g.rand == nil {
		g.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	g.Depth = 1

	g.dungeon = NewMap(config.DungeonWidth, config.DungeonHeight)
	playerStart := g.dungeon.generateMap(g, config.DungeonWidth, config.DungeonHeight)
	g.SpawnPlayer(playerStart)
}
