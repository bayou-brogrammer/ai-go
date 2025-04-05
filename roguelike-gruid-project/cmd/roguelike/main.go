//go:build !js
// +build !js

package main

import (
	"context"
	"math/rand"
	"time"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/game"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ui"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Seed random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	gd := gruid.NewGrid(game.Width, game.Height)
	m := game.NewModel(gd)

	driver := ui.GetDriver()
	app := gruid.NewApp(gruid.AppConfig{
		Model:  m,
		Driver: driver,
	})

	// start application
	if err := app.Start(context.Background()); err != nil {
		driver.Close()
		logrus.Fatal(err)
	}
}
