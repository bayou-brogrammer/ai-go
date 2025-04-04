package main

import (
	"fmt"
	"sync"

	"codeberg.org/anaseto/gruid" // Added gruid import
)

// EntityID represents a unique identifier for an entity.
type EntityID int

type Entity interface{}

// ECS manages entities and their positions.
// This is a simplified version focusing only on position tracking.
type ECS struct {
	nextEntityID EntityID
	mu           sync.RWMutex // To handle concurrent access

	Entities    map[EntityID]Entity      // set of entities
	Positions   map[EntityID]gruid.Point // entity ID: map position
	Renderables map[EntityID]Renderable  // entity ID: map renderable
	Names       map[EntityID]string      // entity ID: map name
}

// NewECS creates and initializes a new ECS.
func NewECS() *ECS {
	return &ECS{
		nextEntityID: 1, // Start IDs from 1
		Entities:     make(map[EntityID]Entity),
		Positions:    make(map[EntityID]gruid.Point),
		Renderables:  make(map[EntityID]Renderable),
		Names:        make(map[EntityID]string),
	}
}

// AddEntity creates a new entity at the given position and returns its ID.
func (ecs *ECS) AddEntity(entity Entity) EntityID {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()

	id := ecs.nextEntityID
	ecs.Entities[id] = entity
	ecs.nextEntityID++
	return id
}

// EntityExists checks if an entity with the given ID exists.
func (ecs *ECS) EntityExists(id EntityID) bool {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()
	_, ok := ecs.Entities[id]
	return ok
}

// DestroyEntity removes an entity from the ECS.
func (ecs *ECS) DestroyEntity(id EntityID) {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()

	delete(ecs.Entities, id)
	delete(ecs.Positions, id)
}

// EntitiesAt returns a slice of EntityIDs located at the given point.
func (ecs *ECS) EntitiesAt(p gruid.Point) []EntityID {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()

	var ids []EntityID
	for id, pos := range ecs.Positions {
		if pos == p {
			ids = append(ids, id)
		}
	}
	return ids
}

// GetAllEntities returns all managed entity IDs.
func (ecs *ECS) GetAllEntities() []EntityID {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()
	ids := make([]EntityID, 0, len(ecs.Positions))
	for id := range ecs.Positions {
		ids = append(ids, id)
	}
	return ids
}

func (ecs *ECS) GetEntitiesWithPositionAndRenderable() []EntityID {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()

	var ids []EntityID
	for id := range ecs.Entities {
		if _, ok := ecs.GetPosition(id); ok {
			if _, ok := ecs.GetRenderable(id); ok {
				ids = append(ids, id)
			}
		}
	}
	return ids
}

// --- Setters ---

// AddPosition adds a position to the given entity.
func (ecs *ECS) AddPosition(id EntityID, position gruid.Point) *ECS {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()

	ecs.Positions[id] = position
	return ecs
}

// AddRenderable adds a renderable to the given entity.
func (ecs *ECS) AddRenderable(id EntityID, renderable Renderable) *ECS {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()

	ecs.Renderables[id] = renderable
	return ecs
}

func (ecs *ECS) AddName(id EntityID, name string) *ECS {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()

	ecs.Names[id] = name
	return ecs
}

// --- Getters ---
// GetPosition retrieves the position of the given entity.
// Returns the position and true if found, otherwise zero Point and false.
func (ecs *ECS) GetPosition(id EntityID) (gruid.Point, bool) {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()

	pos, ok := ecs.Positions[id]
	return pos, ok
}

// GetRenderable retrieves the renderable of the given entity.
// Returns the renderable and true if found, otherwise nil and false.
func (ecs *ECS) GetRenderable(id EntityID) (Renderable, bool) {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()

	renderable, ok := ecs.Renderables[id]
	return renderable, ok
}

func (ecs *ECS) GetName(id EntityID) (string, bool) {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()

	name, ok := ecs.Names[id]
	return name, ok
}

// -- Helpers --

// MoveEntity updates the position of an existing entity.
func (ecs *ECS) MoveEntity(id EntityID, p gruid.Point) error {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()

	if _, ok := ecs.Entities[id]; !ok {
		return fmt.Errorf("entity %d not found", id)
	}

	if _, ok := ecs.Positions[id]; !ok {
		return fmt.Errorf("position %d not found", id)
	}

	ecs.Positions[id] = p
	return nil
}
