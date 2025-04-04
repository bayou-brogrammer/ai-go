package main

import (
	"container/list"

	"codeberg.org/anaseto/gruid"
)

type Name struct {
	Name string
}

// Position component
type Position struct {
	Point gruid.Point
}

// Renderable component
type Renderable struct {
	Glyph rune
	Color gruid.Color
}

// BlocksMovement component
type BlocksMovement struct{}

type Dead struct{}

type Player struct{}

// TurnActor represents an entity that takes turns in the game.
type TurnActor struct {
	Speed        uint64
	Alive        bool
	NextTurnTime uint64
	actions      *list.List
}

// NewTurnActor creates a new TurnActor with the given speed.
func NewTurnActor(speed uint64) TurnActor {
	return TurnActor{
		Speed:        speed,
		Alive:        true,
		NextTurnTime: 0,
		actions:      list.New(),
	}
}

// QueueAction adds an action to the back of the action queue.
func (ta *TurnActor) QueueAction(action GameAction) *TurnActor {
	ta.actions.PushBack(action)
	return ta
}

// AddAction adds an action to the back of the action queue.
func (ta *TurnActor) AddAction(action GameAction) {
	ta.actions.PushBack(action)
}

// NextAction removes and returns the next action from the queue.
func (ta *TurnActor) NextAction() GameAction {
	if ta.actions.Len() == 0 {
		return nil
	}
	element := ta.actions.Front()
	ta.actions.Remove(element)
	return element.Value.(GameAction)
}

// PeakNextAction returns the next action from the queue without removing it.
func (ta *TurnActor) PeakNextAction() GameAction {
	if ta.actions.Len() == 0 {
		return nil
	}

	element := ta.actions.Front()
	return element.Value.(GameAction)
}

// IsAlive returns whether the actor is alive.
func (ta *TurnActor) IsAlive() bool {
	return ta.Alive
}

// WaitingForInput is a component that indicates that the entity is waiting for input.
type WaitingForInput struct{}
