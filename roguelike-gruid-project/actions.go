package main

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
)

type action int

const (
	ActionNone action = iota
	ActionW
	ActionS
	ActionN
	ActionE
	ActionQuit
)

type actionError int

const (
	actionErrorUnknown actionError = iota
)

func (e actionError) Error() string {
	switch e {
	case actionErrorUnknown:
		return "unknown action"
	}
	return ""
}

func (md *model) normalModeAction(action action) (again bool, eff gruid.Effect, err error) {
	g := md.game

	switch action {
	case ActionNone:
		again = true
		err = actionErrorUnknown
	case ActionW, ActionS, ActionN, ActionE:
		again, err = g.PlayerBump(keyToDir(action))
	default:
		fmt.Printf("Unknown action: %v\n", action)
		err = actionErrorUnknown

	}

	if err != nil {
		again = true
	}

	return again, eff, err

}
