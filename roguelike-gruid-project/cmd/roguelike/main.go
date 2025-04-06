//go:build !js
// +build !js

package main

import (
	"context"
	"math/rand"
	"time"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/config"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/game"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ui"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize configuration and parse flags
	config.Init()

	// Setup logger format
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	// Display version information
	logrus.Infof("Starting roguelike game - Debug mode: %v", config.Config.DebugLogging)

	// Seed random number generator
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))

	// Create game grid
	gd := gruid.NewGrid(game.Width, game.Height)

	// Create game model
	m := game.NewModel(gd)

	// Get driver and initialize app
	driver := ui.GetDriver()
	app := gruid.NewApp(gruid.AppConfig{
		Model:  m,
		Driver: driver,
	})

	// Start application
	logrus.Info("Starting game loop")
	if err := app.Start(context.Background()); err != nil {
		driver.Close()
		logrus.Fatal(err)
	}
}
