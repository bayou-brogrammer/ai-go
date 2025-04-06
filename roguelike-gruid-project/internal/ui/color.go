package ui

import "codeberg.org/anaseto/gruid"

// Those constants represent the generic colors we use in this example.
const (
	ColorPlayer gruid.Color = 1 + iota // skip special zero value gruid.ColorDefault
	ColorLOS
	ColorDark
	ColorMonster
	ColorVisibleFloor
	ColorVisibleWall
	ColorExploredFloor
	ColorExploredWall
)

// Those constants represent styling attributes.
const (
	AttrNone gruid.AttrMask = iota
	AttrReverse
)
