package game

import (
	"math/rand"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs/components"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ui"
	"github.com/sirupsen/logrus"
)

func (g *Game) SpawnPlayer(playerStart gruid.Point) {
	// Create the player entity
	playerID := g.ecs.AddEntity(components.PlayerTag{})
	g.PlayerID = playerID // Store the player ID in the game struct

	// Add components to the player
	// Use the start position returned by NewMap
	g.ecs.AddPosition(playerID, playerStart).
		AddRenderable(playerID, components.Renderable{Glyph: '@', Color: ui.ColorPlayer}).
		AddTurnActor(playerID, components.NewTurnActor(100)).
		AddFOV(playerID, components.NewFOVComponent(4, g.Map.Width, g.Map.Height)) // Use correct constructor

	// Add the player to the turn queue
	g.turnQueue.Add(playerID, g.turnQueue.CurrentTime)
}

func (g *Game) SpawnMonster(pos gruid.Point) {
	// Create Orc entity (example)
	monsterID := g.ecs.AddEntity(struct{}{})

	monsterNames := []string{"Orc", "Troll", "Goblin", "Kobold"}
	monsterName := monsterNames[rand.Intn(len(monsterNames))]

	var rune rune
	var speed uint64
	var color gruid.Color = ui.ColorMonster // Default monster color

	switch monsterName {
	case "Orc":
		rune = 'o'
		speed = 100
	case "Troll":
		rune = 'T'
		speed = 200
	case "Goblin":
		rune = 'g'
		speed = 100
		color = ui.ColorSleepingMonster // Goblins use a different color
	case "Kobold":
		rune = 'k'
		speed = 150
	}

	g.ecs.AddName(monsterID, monsterName).
		AddPosition(monsterID, pos).
		AddRenderable(monsterID, components.Renderable{Glyph: rune, Color: color}).
		AddTurnActor(monsterID, components.NewTurnActor(speed)).
		AddAITag(monsterID).
		AddFOV(monsterID, components.NewFOVComponent(6, g.Map.Width, g.Map.Height)) // Use correct constructor

	// Add to turn queue at current turn time
	logrus.Debugf("Created monster ID=%d at position %v, adding to turn queue at time %d",
		monsterID, pos, g.turnQueue.CurrentTime+100)

	// Enqueue the monster
	g.turnQueue.Add(monsterID, g.turnQueue.CurrentTime+100)
}
