# Game Mechanics: Roguelike Gruid

## Core Game Loop

1. Player inputs action via keyboard/mouse
2. Action is processed (movement, attack, etc.)
3. Player's turn ends
4. All AI entities take their turns
5. Game state is updated (FOV, positioning)
6. UI is redrawn
7. Return to step 1

## Movement System

- Grid-based movement on a 2D map
- `EntityBump(entityID, delta)` - Handles movement attempts
- Collision detection via `checkCollision(pos)`
- Movement blocked by:
  - Map boundaries
  - Wall tiles
  - Entities with `BlocksMovement` component
- Implemented via `MoveAction` executing `g.EntityBump()`
- Diagonal movement appears to be supported based on `keyToDir()`

## Combat System

- Basic attack mechanics in `AttackAction.Execute()`
- Target selection through position (bump into target)
- Damage calculation and health reduction
- Death handling via `handleEntityDeath()`
- Visual feedback through hit colors

## Field of Vision (FOV)

- Uses `rl.FOV` for calculations
- Entity FOV controlled by `FOV` component with configurable range
- Map exploration tracking via global bitset
- Visible areas updated by `FOVSystem()`
- Rendering respects FOV through `RenderSystem()`
- Areas outside FOV rendered differently or not at all

## Turn Management

- Turn-based system using priority queue
- Entities have speed parameters determining action frequency
- `TurnQueue` manages entity turn order
- Player turns handled by input system
- Monster turns handled by `monstersTurn()`
- Action execution returns time cost to determine next turn time

## Map Generation

- Procedural generation using room-and-corridor approach
- Rooms created with random dimensions and positions
- Corridors connect rooms via horizontal and vertical tunnels
- Room collision detection prevents overlapping rooms
- Monster placement with random distribution in rooms
- Player start position centered in first room

## Monster AI

- Basic AI through `monstersTurn()`
- Current implementation appears to be simple movement toward player
- Potential for expansion with more complex behaviors
- AI entities identified by `AITag` component

## Game State and Progression

- Current level tracking in `Game.Depth`
- Entity state management through ECS
- Death state through `DeadTag` component
- Message log for game events
- No explicit level transitions observed in the code

## User Interface

- Terminal-based UI via `gruid`
- Map representation with different glyphs for walls/floors
- Entities represented by characters with colors
- Different rendering styles for visible/explored/unseen areas
- Message log for game events display

## Items and Inventory

- No explicit inventory system observed in the current code
- Potential for implementation by adding inventory-related components

## Saving and Loading

- No explicit save/load system observed in the current code
- Potential for implementation using Go's serialization capabilities

## Difficulties and Challenges

- Combat appears to be simple, no complex enemy types observed
- No explicit difficulty scaling with depth
- No explicit player progression system (e.g., experience, levels)

## Random Number Generation

- Game uses its own random source (`Game.rand`)
- Used for map generation, monster placement, and potentially combat
