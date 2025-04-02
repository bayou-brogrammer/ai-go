package main

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
)

func (md *model) Update(msg gruid.Msg) gruid.Effect {
	if _, ok := msg.(gruid.MsgInit); ok {
		return md.init()
	}

	switch msg := msg.(type) {
	case gruid.MsgInit:
		_ = msg
		return nil
	case gruid.MsgQuit:
		return gruid.End()
	case gruid.MsgKeyDown:
		if msg.Key == "q" {
			return gruid.End()
		}
	}

	return md.update(msg)
}

func (md *model) update(msg gruid.Msg) gruid.Effect {
	var eff gruid.Effect
	switch md.mode {
	case modeQuit:
		return nil
	case modeNormal:
		eff = md.updateNormal(msg)
	}

	return eff
}

func (md *model) updateNormal(msg gruid.Msg) gruid.Effect {
	var eff gruid.Effect
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		eff = md.updateKeyDown(msg)
	case gruid.MsgMouse:
		eff = md.updateMouse(msg)
	}
	return eff

}

func (md *model) updateKeyDown(msg gruid.MsgKeyDown) gruid.Effect {
	again, eff, err := md.normalModeKeyDown(msg.Key, msg.Mod&gruid.ModShift != 0)
	if err != nil {
		md.game.Print(err.Error())
	}

	if again {
		return eff
	}

	return md.EndTurn()

}

func (md *model) updateMouse(msg gruid.MsgMouse) gruid.Effect {
	return nil
}

func (md *model) normalModeKeyDown(key gruid.Key, shift bool) (again bool, eff gruid.Effect, err error) {
	action := KEYS_NORMAL[key]
	again, eff, err = md.normalModeAction(action)
	if _, ok := err.(actionError); ok {
		err = fmt.Errorf("key '%s' does nothing. Type ? for help", key)
	}
	return again, eff, err

}
