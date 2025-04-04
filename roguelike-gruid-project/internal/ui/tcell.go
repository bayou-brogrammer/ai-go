//go:build !sdl && !js
// +build !sdl,!js

package ui

import (
	"codeberg.org/anaseto/gruid"
	tcell "codeberg.org/anaseto/gruid-tcell"
	tc "github.com/gdamore/tcell/v2"
)

var driver gruid.Driver

func init() {
	st := styler{}
	dr := tcell.NewDriver(tcell.Config{StyleManager: st})
	dr.PreventQuit()
	driver = dr
}

// styler implements the tcell.StyleManager interface.
type styler struct{}

func (sty styler) GetStyle(st gruid.Style) tc.Style {
	ts := tc.StyleDefault
	switch st.Fg {
	case ColorPlayer:
		ts = ts.Foreground(tc.ColorNavy) // blue color for the player
	case ColorLOS:
		ts = ts.Foreground(tc.ColorYellow)
	case ColorFlashingEnemy:
		ts = ts.Foreground(tc.ColorYellow).Bold(true) // Bold yellow for flashing enemies
	}
	ts = ts.Background(tc.ColorBlack)
	switch st.Bg {
	case ColorDark:
		ts = ts.Background(tc.ColorDefault)
	}
	switch st.Attrs {
	case AttrReverse:
		ts = ts.Reverse(true)
	}
	return ts
}
