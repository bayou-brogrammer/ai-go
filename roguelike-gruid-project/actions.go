package main

import "codeberg.org/anaseto/gruid"

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
	switch action {
	case ActionNone:
		again = true
		err = actionErrorUnknown
	case getMovementActions():
		// again, err = g.PlayerBump(g.Player.P.Add(keyToDir(action)))

	default:
		err = actionErrorUnknown

	}

	if err != nil {
		again = true
	}
	return again, eff, err

}
