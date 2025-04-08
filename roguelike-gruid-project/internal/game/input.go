package game

import "codeberg.org/anaseto/gruid"

var KEYS_NORMAL = map[gruid.Key]action{
	gruid.KeyArrowLeft:  ActionW,
	gruid.KeyArrowDown:  ActionS,
	gruid.KeyArrowUp:    ActionN,
	gruid.KeyArrowRight: ActionE,
	"h":                 ActionW,
	"j":                 ActionS,
	"k":                 ActionN,
	"l":                 ActionE,
	"a":                 ActionW,
	"s":                 ActionS,
	"w":                 ActionN,
	"d":                 ActionE,
	"4":                 ActionW,
	"2":                 ActionS,
	"8":                 ActionN,
	"6":                 ActionE,
	"Q":                 ActionQuit,
}

func keyToDir(k action) (p gruid.Point) {
	switch k {
	case ActionW:
		p = gruid.Point{X: -1, Y: 0}
	case ActionE:
		p = gruid.Point{X: 1, Y: 0}
	case ActionS:
		p = gruid.Point{X: 0, Y: 1}
	case ActionN:
		p = gruid.Point{X: 0, Y: -1}
	}
	return p
}
