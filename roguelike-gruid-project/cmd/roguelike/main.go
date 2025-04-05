//go:build !js
// +build !js

package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
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
	logrus.Debugf("Initializing random number generator with seed: %d", seed)
	rand.New(rand.NewSource(seed))

	// Create game grid
	gd := gruid.NewGrid(game.Width, game.Height)
	logrus.Debugf("Created game grid with dimensions %dx%d", game.Width, game.Height)

	// Create game model
	m := game.NewModel(gd)
	logrus.Debug("Game model initialized")

	// Get driver and initialize app
	driver := ui.GetDriver()
	app := gruid.NewApp(gruid.AppConfig{
		Model:  m,
		Driver: driver,
	})
	logrus.Debug("Game app created with driver")

	// Start application
	logrus.Info("Starting game loop")
	if err := app.Start(context.Background()); err != nil {
		driver.Close()
		logrus.Fatal(err)
	}
}

func subSig(ctx context.Context, msgs chan<- gruid.Msg) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sig)
	select {
	case <-ctx.Done():
	case <-sig:
		msgs <- gruid.MsgQuit{}
	}
}
