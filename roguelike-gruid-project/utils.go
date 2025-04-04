package main

import (
	"fmt"
	"reflect"
)

func GetReflectType(component interface{}) reflect.Type {
	return reflect.TypeOf(component)
}

// getEntityDebugName tries to get a printable name for an entity.
// Corresponds to the Rust helper function.
func GetEntityDebugName(world *World, entityID EntityID) string {
	if name, found := world.GetEntityName(entityID); found {
		return name
	}
	// Fallback to the ID if no name component found
	return fmt.Sprintf("Entity(%d)", entityID)
}
