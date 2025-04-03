package main

import (
	"math/rand"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/rl"

	"slices"
)

// TileType represents the type of a map tile.

const (
	WallCell rl.Cell = iota
	FloorCell
)

// Map represents the game map's logical state and visibility.
type Map struct {
	Grid     rl.Grid // Stores the map cells (rune, style, attributes)
	Width    int
	Height   int
	Visible  map[gruid.Point]bool // Tiles currently visible to the player
	Explored map[gruid.Point]bool // Tiles that have ever been visible
}

// generateMap creates a new map layout with rooms and tunnels.
func generateMap(width, height int) (rl.Grid, gruid.Point) {
	grid := rl.NewGrid(width, height)
	grid.Fill(WallCell)

	var rooms []Rect
	var playerStart gruid.Point

	for range maxRooms {
		w := rand.Intn(roomMaxSize-roomMinSize+1) + roomMinSize
		h := rand.Intn(roomMaxSize-roomMinSize+1) + roomMinSize
		x := rand.Intn(width - w - 1)  // -1 to ensure room fits
		y := rand.Intn(height - h - 1) // -1 to ensure room fits

		newRoom := NewRect(x, y, w, h)

		// Check for intersections with existing rooms
		intersects := slices.ContainsFunc(rooms, newRoom.Intersects)

		if !intersects {
			createRoom(grid, newRoom)
			newCenter := newRoom.Center()

			if len(rooms) == 0 {
				// This is the first room, place player here
				playerStart = newCenter
			} else {
				// Connect to the previous room's center
				prevCenter := rooms[len(rooms)-1].Center()

				// Randomly decide tunnel order (H then V or V then H)
				if rand.Intn(2) == 0 {
					createHTunnel(grid, prevCenter.X, newCenter.X, prevCenter.Y)
					createVTunnel(grid, prevCenter.Y, newCenter.Y, newCenter.X)
				} else {
					createVTunnel(grid, prevCenter.Y, newCenter.Y, prevCenter.X)
					createHTunnel(grid, prevCenter.X, newCenter.X, newCenter.Y)
				}
			}
			rooms = append(rooms, newRoom)
		}
	}

	return grid, playerStart
}

// NewMap creates a new map initialized with walls and visibility data.
// It now returns the map and the player's starting position.
func NewMap(width, height int) (*Map, gruid.Point) {
	m := &Map{
		Grid:     rl.NewGrid(width, height),
		Visible:  make(map[gruid.Point]bool), // Initialize visibility slice
		Explored: make(map[gruid.Point]bool), // Initialize explored slice
		Width:    width,
		Height:   height,
	}

	// Generate the map layout
	var playerStart gruid.Point
	m.Grid, playerStart = generateMap(width, height)

	// Initialize visibility/explored based on the generated map (optional, could be done by FOV later)
	for y := range m.Height {
		for x := range m.Width {
			p := gruid.Point{X: x, Y: y}
			m.Visible[p] = false  // Start not visible
			m.Explored[p] = false // Start not explored
		}
	}

	return m, playerStart
}

// InBounds checks if coordinates are within map bounds.
func (m *Map) InBounds(p gruid.Point) bool {
	return p.X >= 0 && p.X < m.Width && p.Y >= 0 && p.Y < m.Height
}

func (m *Map) isWalkable(p gruid.Point) bool {
	if !m.InBounds(p) {
		return false
	}

	return m.Grid.At(p) == FloorCell
}

// IsWall checks if the tile at the given point is a wall.
func (m *Map) IsWall(p gruid.Point) bool {
	if !m.InBounds(p) {
		return true
	}
	return m.Grid.At(p) == WallCell
}

// IsOpaque checks if a tile blocks FOV. Required by fov.Compute.
func (m *Map) IsOpaque(p gruid.Point) bool {
	return m.IsWall(p) // For now, only walls block sight
}

func (m *Map) Rune(c rl.Cell) (r rune) {
	switch c {
	case WallCell:
		r = '#'
	case FloorCell:
		r = '.'
	}
	return r
}
