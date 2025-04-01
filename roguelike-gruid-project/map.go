package main

import (
	"codeberg.org/anaseto/gruid"
	// "codeberg.org/anaseto/gruid/fov" // We'll add this when needed by FOV code
)

// TileType represents the type of a map tile.
type TileType int

const (
	TileWall TileType = iota
	TileFloor
)

// Map represents the game map's logical state and visibility.
type Map struct {
	LogicalTiles [][]TileType // Stores the actual type of each tile
	Width        int
	Height       int
	Visible      [][]bool // Tiles currently visible to the player
	Explored     [][]bool // Tiles that have ever been visible
}

// NewMap creates a new map initialized with walls and visibility data.
func NewMap(width, height int) *Map {
	m := &Map{
		LogicalTiles: make([][]TileType, height),
		Visible:      make([][]bool, height), // Initialize visibility slice
		Explored:     make([][]bool, height), // Initialize explored slice
		Width:        width,
		Height:       height,
	}
	for y := 0; y < height; y++ {
		m.LogicalTiles[y] = make([]TileType, width)
		m.Visible[y] = make([]bool, width)   // Initialize inner slice
		m.Explored[y] = make([]bool, width)  // Initialize inner slice
		for x := 0; x < width; x++ {
			m.LogicalTiles[y][x] = TileWall // Initialize all as walls
			m.Visible[y][x] = false        // Start not visible
			m.Explored[y][x] = false       // Start not explored
		}
	}
	return m
}

// InBounds checks if coordinates are within map bounds.
func (m *Map) InBounds(p gruid.Point) bool {
	return p.X >= 0 && p.X < m.Width && p.Y >= 0 && p.Y < m.Height
}

// IsWall checks if the tile at the given point is a wall.
func (m *Map) IsWall(p gruid.Point) bool {
	if !m.InBounds(p) {
		return true
	}
	return m.LogicalTiles[p.Y][p.X] == TileWall
}

// IsOpaque checks if a tile blocks FOV. Required by fov.Compute.
func (m *Map) IsOpaque(p gruid.Point) bool {
	return m.IsWall(p) // For now, only walls block sight
}
