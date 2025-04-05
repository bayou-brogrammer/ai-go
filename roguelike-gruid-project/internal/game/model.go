// This file defines the main model of the game: the Update function that
// updates the model state in response to user input, and the Draw function,
// which draws the grid.

package game

import (
	"runtime"

	"codeberg.org/anaseto/gruid"
	"github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/utils"
	"github.com/sirupsen/logrus"
)

type mode int

const (
	modeNormal mode = iota
	modeQuit
)

// Model represents the game model that implements gruid.Model
type Model struct {
	grid gruid.Grid
	game *Game
	mode mode
}

// NewModel creates a new game model
func NewModel(grid gruid.Grid) *Model {
	return &Model{
		grid: grid,
		game: NewGame(),
	}
}

func (md *Model) init() gruid.Effect {
	logrus.Debug("========= Game Initialization Started =========")
	md.game.InitLevel()
	logrus.Debug("Level initialized")

	// Log the state of the turn queue before processing
	logrus.Debug("About to process turn queue for the first time")

	md.processTurnQueue()
	logrus.Debug("Initial turn queue processing completed")

	// Log the player's state after initialization
	logrus.Debug("========= Game Initialization Completed =========")

	if runtime.GOOS == "js" {
		return nil
	}

	return gruid.Sub(utils.HandleSignals)
}

// EndTurn finalizes player's turn and runs other events until next player
// turn.
func (md *Model) EndTurn() gruid.Effect {
	logrus.Debug("EndTurn called - player finished their turn")
	md.mode = modeNormal

	// Add the player back to the turn queue with a standard cost of 100
	g := md.game
	logrus.Debugf("Current turn queue time: %d", g.turnQueue.CurrentTime)
	logrus.Debugf("Re-adding player (ID: %d) to turn queue with time: %d", g.PlayerID, g.turnQueue.CurrentTime+100)
	g.turnQueue.Add(g.PlayerID, g.turnQueue.CurrentTime+100)

	// Mark that we're no longer waiting for player input
	g.ecs.WaitingForInput[g.PlayerID] = false
	logrus.Debug("Set waitingForPlayer = false")

	// Process the turn queue to handle monster turns
	logrus.Debug("Calling processTurnQueue")
	md.processTurnQueue()
	logrus.Debug("processTurnQueue completed")

	// Return nil to indicate the screen should be redrawn
	return nil
}
