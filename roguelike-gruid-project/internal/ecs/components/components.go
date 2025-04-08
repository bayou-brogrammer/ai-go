package components

import (
	"reflect"

	"codeberg.org/anaseto/gruid"
)

// ComponentType is a string identifier for component types
type ComponentType string

// Component type constants
const (
	CPlayerTag  ComponentType = "PlayerTag"
	CPosition   ComponentType = "Position"
	CRenderable ComponentType = "Renderable"
	CName       ComponentType = "Name"
	CAITag      ComponentType = "AITag"
	CTurnActor  ComponentType = "TurnActor"
	CFOV        ComponentType = "FOV"
	CHealth     ComponentType = "Health"
	CCorpseTag  ComponentType = "CorpseTag"
)

var TypeToComponent = map[ComponentType]reflect.Type{
	CPlayerTag:  reflect.TypeOf(PlayerTag{}),
	CPosition:   reflect.TypeOf(gruid.Point{}),
	CRenderable: reflect.TypeOf(Renderable{}),
	CName:       reflect.TypeOf(""),
	CAITag:      reflect.TypeOf(AITag{}),
	CTurnActor:  reflect.TypeOf(TurnActor{}),
	CFOV:        reflect.TypeOf((*FOV)(nil)),
	CHealth:     reflect.TypeOf(Health{}),
	CCorpseTag:  reflect.TypeOf(CorpseTag{}),
}

// GetGoType returns the corresponding Go type for a ComponentType
func GetGoType(compType ComponentType) (reflect.Type, bool) {
	t, ok := TypeToComponent[compType]
	return t, ok
}

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

// Health component represents an entity's health points
type Health struct {
	CurrentHP int
	MaxHP     int
}
