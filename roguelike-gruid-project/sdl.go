//go:build sdl
// +build sdl

package main

import (
	"log"

	"codeberg.org/anaseto/gruid"
	sdl "codeberg.org/anaseto/gruid-sdl"
)

var driver gruid.Driver

func init() {
	t, err := GetTileDrawer()

	if err != nil {
		log.Fatal(err)
	}

	dr := sdl.NewDriver(sdl.Config{
		TileManager: t,
	})

	//dr.SetScale(2.0, 2.0)
	dr.PreventQuit()
	driver = dr
}
