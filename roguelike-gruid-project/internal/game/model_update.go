package game

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
	"github.com/sirupsen/logrus"
)

func (md *Model) Update(msg gruid.Msg) gruid.Effect {
	if _, ok := msg.(gruid.MsgInit); ok {
		return md.init()
	}

	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		if msg.Key == "q" {
			return gruid.End()
		}
	}

	return md.update(msg)
}

func (md *Model) update(msg gruid.Msg) gruid.Effect {
	g := md.game
	waitingForPlayer := g.ecs.WaitingForInput[g.PlayerID]

	// If it's now the player's turn, handle their input
	if waitingForPlayer {
		var eff gruid.Effect
		switch md.mode {
		case modeQuit:
			return nil // Should not happen if waiting for player?
		case modeNormal:
			// Player's turn: handle input message
			eff = md.updateNormal(msg)
		}
		return eff
	} else {
		logrus.Debug("Not waiting for player, processing turn queue")
		md.processTurnQueue()

		// Only monsters moved (or queue was empty/future turns)
		// No player input to process this cycle.
		// We still need to redraw if monsters moved.
		// Returning Flush ensures the screen updates.
		return nil
	}
}

func (md *Model) updateNormal(msg gruid.Msg) gruid.Effect {
	var eff gruid.Effect
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		eff = md.updateKeyDown(msg)
	case gruid.MsgMouse:
		eff = md.updateMouse(msg)
	}
	return eff

}

func (md *Model) updateKeyDown(msg gruid.MsgKeyDown) gruid.Effect {
	again, eff, err := md.normalModeKeyDown(msg.Key, msg.Mod&gruid.ModShift != 0)
	if err != nil {
		md.game.Print(err.Error())
	}

	if again {
		return eff
	}

	return md.EndTurn()

}

func (md *Model) updateMouse(msg gruid.MsgMouse) gruid.Effect {
	return nil
}

func (md *Model) normalModeKeyDown(key gruid.Key, shift bool) (again bool, eff gruid.Effect, err error) {
	action := KEYS_NORMAL[key]
	again, eff, err = md.normalModeAction(action)
	if _, ok := err.(actionError); ok {
		err = fmt.Errorf("key '%s' does nothing. Type ? for help", key)
	}
	return again, eff, err

}

// processTurnQueue processes turns for actors until it's the player's turn
// or the queue is exhausted for the current time step.
// It returns true if the game should wait for player input, false otherwise.
func (md *Model) processTurnQueue() {
	g := md.game
	logrus.Debug("========= processTurnQueue started =========")

	g.turnQueue.PrintQueue()

	for i := range 100 { // Limit to 100 iterations to prevent infinite loops
		logrus.Debugf("Turn queue iteration %d", i)
		logrus.Debugf("Current turn queue time: %d", g.turnQueue.CurrentTime)

		// Check if the queue is empty
		if g.turnQueue.IsEmpty() {
			logrus.Debug("Turn queue is empty.")
			logrus.Debug("========= processTurnQueue ended (queue empty) =========")
			return
		}

		// Peek at the next actor
		entry, ok := g.turnQueue.Peek()
		if !ok { // Should not happen if IsEmpty is false, but good practice
			logrus.Debug("Error: Queue not empty but failed to peek.")
			// Clear active entity
			logrus.Debug("========= processTurnQueue ended (peek error) =========")
			return
		}

		logrus.Debugf("Next entry in queue: EntityID=%d, Time=%d", entry.EntityID, entry.Time)

		// It's time for the next actor's turn, pop them
		actorEntry, _ := g.turnQueue.Next() // We already know ok is true from Peek
		logrus.Debugf("Popped actor from queue: EntityID=%d, Time=%d", actorEntry.EntityID, actorEntry.Time)

		// Update game time to the current actor's time
		g.turnQueue.CurrentTime = actorEntry.Time
		logrus.Debugf("Updated game time to: %d", g.turnQueue.CurrentTime)

		// Check if it's the player
		if actorEntry.EntityID == g.PlayerID {
			logrus.Debug("It's the player's turn.")
			g.ecs.WaitingForInput[g.PlayerID] = true
			logrus.Debug("========= processTurnQueue ended (player's turn) =========")
			return
		}

		// It's a monster's turn
		logrus.Debugf("Handling turn for monster %d at time %d", actorEntry.EntityID, actorEntry.Time)
		if g.ecs.HasAITag(actorEntry.EntityID) { // Ensure it's actually an AI
			cost, err := handleMonsterTurn(g, actorEntry.EntityID)
			if err != nil {
				logrus.Errorf("Error during monster %d turn: %v", actorEntry.EntityID, err)
				// Decide how to handle monster errors - skip turn? Give basic cost?
				// For now, let's give it a standard cost and re-queue
				cost = 100
			}
			// Re-queue the monster for its next turn
			nextTime := g.turnQueue.CurrentTime + uint64(cost)
			logrus.Debugf("Re-queuing monster %d for time %d (cost=%d)", actorEntry.EntityID, nextTime, cost)
			g.turnQueue.Add(actorEntry.EntityID, nextTime)
		} else {
			// Entity in queue that isn't player or AI? Log warning.
			logrus.Debugf("Warning: Entity %d in turn queue is not player or AI.", actorEntry.EntityID)
			// Don't re-queue for now.
		}
		// Loop continues to process next actor in the queue at the current time
	}

	// Clear active entity at the end
	logrus.Debug("========= processTurnQueue ended (iteration limit reached) =========")
}
