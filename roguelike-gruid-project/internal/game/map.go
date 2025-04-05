package game

import (
	"math/rand"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/rl"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/components"
	"github.com/sirupsen/logrus"

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

// NewMap creates a new map initialized with walls and visibility data.
// It now returns the map and the player's starting position.
func NewMap(width, height int) *Map {
	m := &Map{
		Grid:     rl.NewGrid(width, height),
		Visible:  make(map[gruid.Point]bool), // Initialize visibility slice
		Explored: make(map[gruid.Point]bool), // Initialize explored slice
		Width:    width,
		Height:   height,
	}

	return m
}

// generateMap creates a new map layout with rooms and tunnels.
// generateMap creates a new map layout with rooms and tunnels, and spawns monsters.
// It now takes the game struct to access ECS and TurnQueue.
// generateMap creates a new map layout with rooms and tunnels, and spawns monsters.
// It now takes the game struct to access ECS and TurnQueue.
func (m *Map) generateMap(g *Game, width, height int) gruid.Point { // Added *Game param
	m.Grid.Fill(WallCell)

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
			createRoom(m.Grid, newRoom)
			newCenter := newRoom.Center()

			if len(rooms) == 0 {
				// This is the first room, place player here
				playerStart = newCenter
			} else {
				// Connect to the previous room's center
				prevCenter := rooms[len(rooms)-1].Center()

				// Randomly decide tunnel order (H then V or V then H)
				if rand.Intn(2) == 0 {
					createHTunnel(m.Grid, prevCenter.X, newCenter.X, prevCenter.Y)
					createVTunnel(m.Grid, prevCenter.Y, newCenter.Y, newCenter.X)
				} else {
					createVTunnel(m.Grid, prevCenter.Y, newCenter.Y, prevCenter.X)
					createHTunnel(m.Grid, prevCenter.X, newCenter.X, newCenter.Y)
				}
				// Spawn monsters in this room (if not the first room)
				m.placeMonsters(g, newRoom)
			}
			rooms = append(rooms, newRoom)
		}
	}

	return playerStart
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

// placeMonsters spawns monsters in a given room.
func (m *Map) placeMonsters(g *Game, room Rect) {
	// Determine number of monsters for this room (e.g., 0 to maxMonstersPerRoom)
	numMonsters := rand.Intn(maxMonstersPerRoom + 1) // +1 because Intn is exclusive upper bound
	logrus.Debugf("Placing %d monsters in room: %v", numMonsters, room)

	for i := 0; i < numMonsters; i++ {
		// Find a random walkable tile within the room bounds
		// Add +1 to x1, y1 and -1 to x2, y2 to avoid spawning on walls
		x := rand.Intn(room.X2-room.X1-1) + room.X1 + 1
		y := rand.Intn(room.Y2-room.Y1-1) + room.Y1 + 1
		pos := gruid.Point{X: x, Y: y}

		// Check if the tile is walkable and not already occupied
		if m.isWalkable(pos) && len(g.ecs.EntitiesAt(pos)) == 0 {
			// Create Orc entity (example)
			monsterID := g.ecs.AddEntity(struct{}{}) // Using empty struct for now
			g.ecs.AddName(monsterID, "Orc")
			g.ecs.AddPosition(monsterID, pos)
			g.ecs.AddRenderable(monsterID, components.Renderable{Glyph: 'o', Color: gruid.ColorDefault}) // Example renderable (using default color for now)
			g.ecs.AddAITag(monsterID, components.AITag{})                                                // Mark as AI-controlled

			// Add to turn queue at current turn time
			logrus.Debugf("Created monster ID=%d at position %v, adding to turn queue at time %d",
				monsterID, pos, g.turnQueue.CurrentTime+100)
			g.turnQueue.Add(monsterID, g.turnQueue.CurrentTime+100)
		} else {
			// If tile is occupied or not walkable, we just skip spawning this monster for simplicity
			logrus.Debugf("Failed to spawn monster at position %v - not walkable or occupied", pos)
		}
	}
}

// Constants like maxRooms, roomMinSize, roomMaxSize, maxMonstersPerRoom should be defined centrally (e.g., in game.go)
