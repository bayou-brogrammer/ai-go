package main

const (
	fovRadius = 8 // How far the player can see
)

// Game represents the main game state.
type game struct {
	Map *Map
	ecs *World // The Entity-Component-System manager
}
