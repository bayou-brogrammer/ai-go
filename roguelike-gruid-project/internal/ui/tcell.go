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

	// Map foreground colors
	switch st.Fg {
	case ColorPlayer:
		ts = ts.Foreground(tc.ColorBlue)
	case ColorMonster:
		ts = ts.Foreground(tc.ColorRed)
	case ColorSleepingMonster:
		ts = ts.Foreground(tc.ColorPurple)
	case ColorConfusedMonster:
		ts = ts.Foreground(tc.ColorGreen)
	case ColorParalyzedMonster:
		ts = ts.Foreground(tc.ColorAqua)
	case ColorItem:
		ts = ts.Foreground(tc.ColorYellow)
	case ColorSpecialItem:
		ts = ts.Foreground(tc.ColorPurple)
	case ColorVisibleWall, ColorVisibleFloor:
		ts = ts.Foreground(tc.ColorWhite)
	case ColorExploredWall, ColorExploredFloor:
		ts = ts.Foreground(tc.ColorDarkGray)
	case ColorUITitle:
		ts = ts.Foreground(tc.ColorWhite)
	case ColorUIText:
		ts = ts.Foreground(tc.ColorLightGray)
	case ColorUIHighlight:
		ts = ts.Foreground(tc.ColorYellow)
	case ColorHealthOk:
		ts = ts.Foreground(tc.ColorGreen)
	case ColorHealthWounded:
		ts = ts.Foreground(tc.ColorYellow)
	case ColorHealthCritical:
		ts = ts.Foreground(tc.ColorRed)
	case ColorStatusGood:
		ts = ts.Foreground(tc.ColorBlue)
	case ColorStatusBad:
		ts = ts.Foreground(tc.ColorRed)
	case ColorStatusNeutral:
		ts = ts.Foreground(tc.ColorYellow)
	}

	// Set background color
	ts = ts.Background(tc.ColorBlack)
	switch st.Bg {
	case ColorUIBackground:
		ts = ts.Background(tc.ColorBlack)
	case ColorUIBorder:
		ts = ts.Background(tc.ColorDarkGray)
	}

	// Handle attributes
	if st.Attrs&AttrReverse != 0 {
		ts = ts.Reverse(true)
	}
	if st.Attrs&AttrBlink != 0 {
		ts = ts.Blink(true)
	}
	if st.Attrs&AttrUnderline != 0 {
		ts = ts.Underline(true)
	}
	if st.Attrs&AttrBold != 0 {
		ts = ts.Bold(true)
	}

	return ts
}
