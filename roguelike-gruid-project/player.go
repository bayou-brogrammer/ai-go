package main

import (
	"errors"

	"codeberg.org/anaseto/gruid"
)

// checkCollision checks if a given position is a valid move
func (g *game) checkCollision(pos gruid.Point) bool {
	if !g.Map.InBounds(pos) {
		return true // Out of bounds
	}

	// Check for blocking entities
	for _, id := range g.ecs.EntitiesAt(pos) {
		if id == g.PlayerID {
			continue
		}
		otherPos, ok := g.ecs.GetPosition(id)
		if ok && otherPos == pos {
			return true // Collision with blocking entity
		}
	}

	return false
}

func (g *game) PlayerBump(delta gruid.Point) (bool, error) {
	pid := g.PlayerID
	playerPos, ok := g.ecs.GetPosition(pid)
	if !ok {
		return false, errors.New("player position not found")
	}

	newPos := playerPos.Add(delta)

	if !g.Map.isWalkable(newPos) || g.checkCollision(newPos) {
		return false, nil
	}

	g.ecs.MoveEntity(pid, newPos)

	return true, nil
}
