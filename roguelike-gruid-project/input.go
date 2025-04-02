package main

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

func getMovementActions() action {
	return ActionW | ActionS | ActionN | ActionE
}

func keyToDir(key gruid.Key) action {
	if action, ok := KEYS_NORMAL[key]; ok {
		return action
	}
	return ActionNone
}
