package game

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
	"github.com/sirupsen/logrus"
)

func (md *Model) Update(msg gruid.Msg) gruid.Effect {
	if _, ok := msg.(gruid.MsgInit); ok {
		return md.init()
	}

	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		if msg.Key == "q" {
			return gruid.End()
		}
	}

	return md.update(msg)
}

func (md *Model) update(msg gruid.Msg) gruid.Effect {

	g := md.game
	waitingForPlayer := g.waitingForInput
	if waitingForPlayer {
		var eff gruid.Effect
		switch md.mode {
		case modeQuit:
			return nil // Should not happen if waiting for player?
		case modeNormal:
			eff = md.updateNormal(msg)
		}
		return eff
	} else {
		logrus.Debug("Not waiting for player, processing turn queue")
		md.processTurnQueue()

		// Only monsters moved (or queue was empty/future turns)
		// No player input to process this cycle.
		// We still need to redraw if monsters moved.
		// Returning Flush ensures the screen updates.
		return nil
	}
}

func (md *Model) updateNormal(msg gruid.Msg) gruid.Effect {
	var eff gruid.Effect
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		eff = md.updateKeyDown(msg)
	case gruid.MsgMouse:
		eff = md.updateMouse(msg)
	}
	return eff

}

func (md *Model) updateKeyDown(msg gruid.MsgKeyDown) gruid.Effect {
	again, eff, err := md.normalModeKeyDown(msg.Key, msg.Mod&gruid.ModShift != 0)
	if err != nil {
		md.game.Print(err.Error())
	}

	if again {
		return eff
	}

	return md.EndTurn()

}

func (md *Model) updateMouse(msg gruid.MsgMouse) gruid.Effect {
	return nil
}

func (md *Model) normalModeKeyDown(key gruid.Key, shift bool) (again bool, eff gruid.Effect, err error) {
	action := KEYS_NORMAL[key]
	again, eff, err = md.normalModeAction(action)
	if _, ok := err.(actionError); ok {
		err = fmt.Errorf("key '%s' does nothing. Type ? for help", key)
	}
	return again, eff, err

}
