package main

const (
// fovRadius = 8 // How far the player can see
)

// Game represents the main game state.
type game struct {
	Map       *Map
	ecs       *World     // The Entity-Component-System manager
	PlayerID  EntityID   // Store the player's entity ID
	turnQueue *TurnQueue // Event queue for game events
}

func NewGame(gameMap *Map) *game {
	return &game{
		Map:       gameMap,
		ecs:       NewWorld(),
		turnQueue: NewTurnQueue(),
	}
}
