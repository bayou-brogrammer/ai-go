package ecs

import (
	"fmt"
	"sync"

	"codeberg.org/anaseto/gruid" // Added gruid import
	. "github.com/lecoqjacob/ai-go/roguelike-gruid-project/internal/ecs/components"
	"github.com/sirupsen/logrus"
)

// EntityID represents a unique identifier for an entity.
type EntityID int

type Entity interface{}

// ECS manages entities and their components.
type ECS struct {
	nextEntityID EntityID
	mu           sync.RWMutex // To handle concurrent access

	Entities    map[EntityID]Entity      // set of entities
	Positions   map[EntityID]gruid.Point // entity ID: map position
	Renderables map[EntityID]Renderable  // entity ID: map renderable
	Names       map[EntityID]string      // entity ID: map name
	AITags      map[EntityID]AITag       // entity ID: presence of AI tag
	TurnActors  map[EntityID]TurnActor   // entity ID: presence of TurnActor
	FOVs        map[EntityID]*FOV        // entity ID: FOV
}

// NewECS creates and initializes a new ECS.
func NewECS() *ECS {
	return &ECS{
		nextEntityID: 1, // Start IDs from 1
		Entities:     make(map[EntityID]Entity),
		Positions:    make(map[EntityID]gruid.Point),
		Renderables:  make(map[EntityID]Renderable),
		Names:        make(map[EntityID]string),
		AITags:       make(map[EntityID]AITag),
		TurnActors:   make(map[EntityID]TurnActor),
		FOVs:         make(map[EntityID]*FOV),
	}
}

// AddEntity creates a new entity and returns its ID.
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

// DestroyEntity removes an entity and all its components from the ECS.
func (ecs *ECS) DestroyEntity(id EntityID) {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()

	delete(ecs.Entities, id)
	delete(ecs.Positions, id)
	delete(ecs.Renderables, id)
	delete(ecs.Names, id)
	delete(ecs.AITags, id) // Remove AI tag as well
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
	ids := make([]EntityID, 0, len(ecs.Entities)) // Iterate over Entities map for completeness
	for id := range ecs.Entities {
		ids = append(ids, id)
	}
	return ids
}

// GetEntitiesWithPositionAndRenderable queries entities having both Position and Renderable components.
func (ecs *ECS) GetEntitiesWithPositionAndRenderable() []EntityID {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()

	var ids []EntityID
	for id := range ecs.Entities {
		// Check presence in both maps directly
		_, posOk := ecs.Positions[id]
		_, renOk := ecs.Renderables[id]
		if posOk && renOk {
			ids = append(ids, id)
		}
	}
	return ids
}

// GetEntitiesWithPositionAndFOV queries entities having both Position and FOV components.
func (ecs *ECS) GetEntitiesWithPositionAndFOV() []EntityID {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()

	var ids []EntityID
	// Iterate over the smaller map (likely FOVs) for potential efficiency
	for id := range ecs.FOVs {
		// Check if the entity also has a Position
		if _, posOk := ecs.Positions[id]; posOk {
			// Ensure the entity actually still exists (optional, but safer)
			if ecs.entityExistsUnlocked(id) {
				ids = append(ids, id)
			}
		}
	}
	return ids
}


// --- Setters ---

// AddPosition adds or updates the Position component for an entity.
func (ecs *ECS) AddPosition(id EntityID, position gruid.Point) *ECS {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()
	if !ecs.entityExistsUnlocked(id) {
		logrus.Debugf("Warning: Attempted to add Position to non-existent entity %d\n", id)
		return ecs
	}
	ecs.Positions[id] = position
	return ecs
}

// AddRenderable adds or updates the Renderable component for an entity.
func (ecs *ECS) AddRenderable(id EntityID, renderable Renderable) *ECS {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()
	if !ecs.entityExistsUnlocked(id) {
		logrus.Debugf("Warning: Attempted to add Renderable to non-existent entity %d\n", id)
		return ecs
	}
	ecs.Renderables[id] = renderable
	return ecs
}

// AddName adds or updates the Name component for an entity.
func (ecs *ECS) AddName(id EntityID, name string) *ECS {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()
	if !ecs.entityExistsUnlocked(id) {
		logrus.Debugf("Warning: Attempted to add Name to non-existent entity %d\n", id)
		return ecs
	}
	ecs.Names[id] = name
	return ecs
}

func (ecs *ECS) AddAITag(id EntityID) *ECS {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()
	if !ecs.entityExistsUnlocked(id) {
		logrus.Debugf("Warning: Attempted to add AITag to non-existent entity %d\n", id)
		return ecs
	}
	ecs.AITags[id] = AITag{}
	return ecs
}

func (ecs *ECS) AddTurnActor(id EntityID, actor TurnActor) *ECS {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()
	if !ecs.entityExistsUnlocked(id) {
		logrus.Debugf("Warning: Attempted to add TurnActor to non-existent entity %d\n", id)
		return ecs
	}

	ecs.TurnActors[id] = actor
	return ecs
}

func (ecs *ECS) AddFOV(id EntityID, fov *FOV) *ECS {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()
	if !ecs.entityExistsUnlocked(id) {
		logrus.Debugf("Warning: Attempted to add FOV to non-existent entity %d\n", id)
		return ecs
	}
	ecs.FOVs[id] = fov
	return ecs
}

// --- Getters ---

// GetPosition retrieves the Position component for an entity.
func (ecs *ECS) GetPosition(id EntityID) (gruid.Point, bool) {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()
	pos, ok := ecs.Positions[id]
	return pos, ok
}

// GetRenderable retrieves the Renderable component for an entity.
func (ecs *ECS) GetRenderable(id EntityID) (Renderable, bool) {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()
	renderable, ok := ecs.Renderables[id]
	return renderable, ok
}

// GetName retrieves the Name component for an entity.
func (ecs *ECS) GetName(id EntityID) (string, bool) {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()
	name, ok := ecs.Names[id]
	return name, ok
}

func (ecs *ECS) GetTurnActor(id EntityID) (TurnActor, bool) {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()
	actor, ok := ecs.TurnActors[id]
	return actor, ok
}

func (ecs *ECS) GetTurnActorUnchecked(id EntityID) TurnActor {
	return ecs.TurnActors[id]
}

func (ecs *ECS) GetFOV(id EntityID) (*FOV, bool) {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()
	fov, ok := ecs.FOVs[id]
	return fov, ok
}

// HasAITag checks if an entity has the AITag component.
func (ecs *ECS) HasAITag(id EntityID) bool {
	ecs.mu.RLock()
	defer ecs.mu.RUnlock()
	_, ok := ecs.AITags[id]
	return ok
}

// --- Helpers ---

// entityExistsUnlocked is an internal helper to check entity existence without locking (caller must hold lock).
func (ecs *ECS) entityExistsUnlocked(id EntityID) bool {
	_, ok := ecs.Entities[id]
	return ok
}

// MoveEntity updates the position of an existing entity.
func (ecs *ECS) MoveEntity(id EntityID, p gruid.Point) error {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()

	if !ecs.entityExistsUnlocked(id) {
		return fmt.Errorf("entity %d not found", id)
	}

	if _, ok := ecs.Positions[id]; !ok {
		// This case should ideally not happen if AddPosition checks existence
		return fmt.Errorf("entity %d exists but has no Position component", id)
	}

	ecs.Positions[id] = p
	return nil
}
