# Entity Component System

## Description

The Entity Component System (ECS) architecture is the backbone of the game's object management. It provides a flexible, cache-friendly way to define game objects as collections of components rather than deep inheritance hierarchies.

## Requirements

- [x] Define a basic Entity type that is a simple numeric ID
- [x] Implement component storage using maps for flexible component types
- [x] Support adding components to entities
- [x] Support removing components from entities
- [x] Support querying entities with specific components
- [x] Implement a system for retrieving single components from entities
- [x] Add helper methods for common component operations
- [ ] Support component serialization for save/load functionality
- [ ] Support entity templates or archetypes for common entity types
- [ ] Add component change notification system

## Components

- **Position** - Stores a position in the game world using gruid.Point
- **Renderable** - Contains glyph and color information for rendering
- **Health** - Tracks current and maximum hit points
- **BlocksMovement** - Tag component indicating entity blocks movement
- **PlayerTag** - Tag component identifying the player entity
- **AITag** - Tag component for entities controlled by AI
- **TurnActor** - Contains data for entities that take turns
- **FOV** - Field of view component for visibility calculation
- **Name** - Contains entity's name
- **CorpseTag** - Tag component for dead entities

## Systems

The ECS architecture uses systems that operate on entities with specific component combinations:

- **RenderSystem** - Draws entities with Position and Renderable components
- **MovementSystem** - Handles collision detection and position updates
- **CombatSystem** - Handles attacks and damage
- **FOVSystem** - Updates visibility for entities with FOV components

## Acceptance Criteria

- The ECS must efficiently support adding/removing components at runtime
- Component queries must be performant, even with thousands of entities
- The ECS must be thread-safe for future concurrent system processing
- Tag components must have minimal memory overhead
- Component access should be type-safe when possible
