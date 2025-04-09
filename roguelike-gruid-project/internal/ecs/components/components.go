package components

import (
	"reflect"

	"codeberg.org/anaseto/gruid"
)

// ComponentType is a string identifier for component types
type ComponentType string

// Component type constants
const (
	CAITag          ComponentType = "AITag"
	CBlocksMovement ComponentType = "BlocksMovement"
	CCorpseTag      ComponentType = "CorpseTag"
	CFOV            ComponentType = "FOV"
	CHealth         ComponentType = "Health"
	CName           ComponentType = "Name"
	CPlayerTag      ComponentType = "PlayerTag"
	CPosition       ComponentType = "Position"
	CRenderable     ComponentType = "Renderable"
	CTurnActor      ComponentType = "TurnActor"
)

var TypeToComponent = map[ComponentType]reflect.Type{
	CAITag:          reflect.TypeOf(AITag{}),
	CBlocksMovement: reflect.TypeOf(BlocksMovement{}),
	CCorpseTag:      reflect.TypeOf(CorpseTag{}),
	CFOV:            reflect.TypeOf((*FOV)(nil)),
	CHealth:         reflect.TypeOf(Health{}),
	CName:           reflect.TypeOf(""),
	CPlayerTag:      reflect.TypeOf(PlayerTag{}),
	CPosition:       reflect.TypeOf(gruid.Point{}),
	CRenderable:     reflect.TypeOf(Renderable{}),
	CTurnActor:      reflect.TypeOf(TurnActor{}),
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

func NewHealth(maxHP int) Health {
	return Health{
		CurrentHP: maxHP,
		MaxHP:     maxHP,
	}
}

func (h *Health) IsDead() bool {
	return h.CurrentHP <= 0
}
