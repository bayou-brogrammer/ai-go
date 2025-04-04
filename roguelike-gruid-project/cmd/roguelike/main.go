//go:build !js
// +build !js

package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/game"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ui"
)

func main() {
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
		log.Fatal(err)
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
