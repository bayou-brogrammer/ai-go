package game

import (
	"math/rand"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs/components"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ui"
	"github.com/sirupsen/logrus"
)

func (g *Game) SpawnPlayer(playerStart gruid.Point) {
	logrus.Debugf("Spawning player at %v", playerStart)
	playerID := g.ecs.AddEntity()
	g.PlayerID = playerID // Store the player ID in the game struct

	g.ecs.AddComponents(playerID,
		playerStart,
		components.PlayerTag{},
		components.BlocksMovement{},
		components.Name{Name: "Player"},
		components.Renderable{Glyph: '@', Color: ui.ColorPlayer},
		components.NewHealth(10),
		components.NewTurnActor(100),
		components.NewFOVComponent(4, g.Map.Width, g.Map.Height),
	)
	g.turnQueue.Add(playerID, g.turnQueue.CurrentTime)

}

func (g *Game) SpawnMonster(pos gruid.Point) {
	monsterID := g.ecs.AddEntity()

	monsterNames := []string{"Orc", "Troll", "Goblin", "Kobold"}
	monsterName := monsterNames[rand.Intn(len(monsterNames))]

	var rune rune
	var speed uint64
	var color gruid.Color = ui.ColorMonster // Default monster color
	var maxHP int

	switch monsterName {
	case "Orc":
		rune = 'o'
		speed = 100
		maxHP = 1
	case "Troll":
		rune = 'T'
		speed = 200
		maxHP = 1
	case "Goblin":
		rune = 'g'
		speed = 100
		color = ui.ColorSleepingMonster // Goblins use a different color
		maxHP = 1
	case "Kobold":
		rune = 'k'
		speed = 150
		maxHP = 1
	}

	g.ecs.AddComponents(monsterID,
		pos,
		components.AITag{},
		components.BlocksMovement{},
		components.Name{Name: monsterName},
		components.Renderable{Glyph: rune, Color: color},
		components.NewHealth(maxHP),
		components.NewFOVComponent(6, g.Map.Width, g.Map.Height),
		components.NewTurnActor(speed),
	)

	logrus.Debugf("Created monster ID=%d at position %v, adding to turn queue at time %d",
		monsterID, pos, g.turnQueue.CurrentTime+100)
	g.turnQueue.Add(monsterID, g.turnQueue.CurrentTime+100)

}
