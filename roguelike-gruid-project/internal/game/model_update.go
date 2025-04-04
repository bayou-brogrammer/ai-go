package game

import (
	"fmt"
	"log"

	"codeberg.org/anaseto/gruid"
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
	// First, process the turn queue for any pending monster actions
	waitForPlayer := md.processTurnQueue()

	// If it's now the player's turn, handle their input
	if waitForPlayer {
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
func (md *Model) processTurnQueue() (waitForPlayerInput bool) {
	g := md.game
	for {
		// Check if the queue is empty
		if g.turnQueue.IsEmpty() {
			log.Println("Turn queue is empty.")
			return false // Nothing to process
		}

		// Peek at the next actor
		entry, ok := g.turnQueue.Peek()
		if !ok { // Should not happen if IsEmpty is false, but good practice
			log.Println("Error: Queue not empty but failed to peek.")
			return false
		}

		// If the next actor's turn is in the future, stop processing for now
		if entry.Time > g.turnQueue.CurrentTime {
			// log.Printf("Next turn (%d) is in the future (current: %d). Waiting.", entry.Time, g.turnQueue.CurrentTime) // Optional debug
			return false // Wait for time to advance (implicitly via player action)
		}

		// It's time for the next actor's turn, pop them
		actorEntry, _ := g.turnQueue.Next() // We already know ok is true from Peek

		// Update game time to the current actor's time
		g.turnQueue.CurrentTime = actorEntry.Time

		// Check if it's the player
		if actorEntry.EntityID == g.PlayerID {
			// log.Println("Player's turn.") // Optional debug
			return true // Signal that we need player input
		}

		// It's a monster's turn
		// log.Printf("Handling turn for monster %d at time %d", actorEntry.EntityID, actorEntry.Time) // Optional debug
		if g.ecs.HasAITag(actorEntry.EntityID) { // Ensure it's actually an AI
			cost, err := handleMonsterTurn(g, actorEntry.EntityID)
			if err != nil {
				log.Printf("Error during monster %d turn: %v", actorEntry.EntityID, err)
				// Decide how to handle monster errors - skip turn? Give basic cost?
				// For now, let's give it a standard cost and re-queue
				cost = 100
			}
			// Re-queue the monster for its next turn
			g.turnQueue.Add(actorEntry.EntityID, g.turnQueue.CurrentTime+uint64(cost))
		} else {
			// Entity in queue that isn't player or AI? Log warning.
			log.Printf("Warning: Entity %d in turn queue is not player or AI.", actorEntry.EntityID)
			// Don't re-queue for now.
		}
		// Loop continues to process next actor in the queue at the current time
	}
}
