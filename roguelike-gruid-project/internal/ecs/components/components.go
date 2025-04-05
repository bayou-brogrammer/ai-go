package components

import (
	"codeberg.org/anaseto/gruid"
)

// Name component represents an entity's name
type Name struct {
	Name string
}

// Position component represents an entity's position in the game world
type Position struct {
	Point gruid.Point
}

// Renderable component represents how an entity is rendered
type Renderable struct {
	Glyph rune
	Color gruid.Color
}
