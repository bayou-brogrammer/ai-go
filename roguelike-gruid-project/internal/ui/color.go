package ui

import "codeberg.org/anaseto/gruid"

// Those constants represent the generic colors we use in this example.
const (
	ColorPlayer gruid.Color = 1 + iota // skip special zero value gruid.ColorDefault
	ColorLOS
	ColorDark
	ColorFlashingEnemy // New color for flashing enemies during their turn
)

// Those constants represent styling attributes.
const (
	AttrNone gruid.AttrMask = iota
	AttrReverse
)
