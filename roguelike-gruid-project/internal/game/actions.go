package game

import (
	"codeberg.org/anaseto/gruid"
	"github.com/sirupsen/logrus"
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

func (md *Model) normalModeAction(action action) (again bool, eff gruid.Effect, err error) {
	g := md.game

	logrus.Debugf("Normal mode action: %v\n", action)

	switch action {
	case ActionNone:
		again = true
		err = actionErrorUnknown
	case ActionW, ActionS, ActionN, ActionE:
		action := MoveAction{
			Direction: keyToDir(action),
			EntityID:  g.PlayerID,
		}
		actor, _ := g.ecs.GetTurnActor(g.PlayerID)
		actor.AddAction(action)

		return false, eff, nil

	default:
		logrus.Debugf("Unknown action: %v\n", action)
		err = actionErrorUnknown

	}

	if err != nil {
		again = true
	}

	return again, eff, err

}
