package main

import (
	"reflect"
	"sync"
)

// EntityID represents a unique identifier for an entity.
type EntityID int

// World manages all entities and components.
type World struct {
	entities     map[EntityID]map[reflect.Type]any
	nextEntityID EntityID
	mu           sync.RWMutex // To handle concurrent access if needed later
}

// NewWorld creates and initializes a new World.
func NewWorld() *World {
	return &World{
		entities:     make(map[EntityID]map[reflect.Type]any),
		nextEntityID: 1, // Start IDs from 1
	}
}

// CreateEntity generates a new unique EntityID and prepares its component map.
func (w *World) CreateEntity() EntityID {
	w.mu.Lock()
	defer w.mu.Unlock()

	id := w.nextEntityID
	w.nextEntityID++
	w.entities[id] = make(map[reflect.Type]any) // Initialize component map
	return id
}

// AddComponent adds a component to the specified entity.
// If the entity already has a component of the same type, it's overwritten.
func (w *World) AddComponent(id EntityID, component interface{}) {
	w.mu.Lock()
	defer w.mu.Unlock()

	compType := reflect.TypeOf(component)
	if _, ok := w.entities[id]; !ok {
		// Optionally handle error: entity does not exist
		// For now, we assume the entity exists or create its map entry
		w.entities[id] = make(map[reflect.Type]any)
	}
	w.entities[id][compType] = component
}

// GetComponent retrieves a component of a specific type for the given entity.
// Returns the component and true if found, otherwise nil and false.
func (w *World) GetComponent(id EntityID, compType reflect.Type) (any, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if comps, ok := w.entities[id]; ok {
		comp, found := comps[compType]
		return comp, found
	}
	return nil, false
}

// HasComponent checks if an entity possesses a component of the specified type.
func (w *World) HasComponent(id EntityID, compType reflect.Type) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if comps, ok := w.entities[id]; ok {
		_, found := comps[compType]
		return found
	}
	return false
}

// RemoveComponent removes a component of a specific type from an entity.
func (w *World) RemoveComponent(id EntityID, compType reflect.Type) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if comps, ok := w.entities[id]; ok {
		delete(comps, compType)
	}
}

// DestroyEntity removes an entity and all its associated components from the world.
func (w *World) DestroyEntity(id EntityID) {
	w.mu.Lock()
	defer w.mu.Unlock()

	delete(w.entities, id)
}

// --- System Query Helpers (Example - can be expanded) ---

// QueryEntitiesWithComponents returns a list of EntityIDs that have all the specified component types.
func (w *World) QueryEntitiesWithComponents(compTypes ...reflect.Type) []EntityID {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var result []EntityID
	for id, comps := range w.entities {
		match := true
		for _, reqType := range compTypes {
			if _, found := comps[reqType]; !found {
				match = false
				break
			}
		}
		if match {
			result = append(result, id)
		}
	}
	return result
}
