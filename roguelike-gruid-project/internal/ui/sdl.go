//go:build sdl
// +build sdl

package ui

import (
	"codeberg.org/anaseto/gruid"
	sdl "codeberg.org/anaseto/gruid-sdl"
	"github.com/sirupsen/logrus"
)

var driver gruid.Driver

func init() {
	t, err := GetTileDrawer()

	if err != nil {
		logrus.Fatal(err)
	}

	dr := sdl.NewDriver(sdl.Config{
		TileManager: t,
	})

	//dr.SetScale(2.0, 2.0)
	dr.PreventQuit()
	driver = dr
}
