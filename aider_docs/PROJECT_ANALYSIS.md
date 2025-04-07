# Roguelike Gruid Project Analysis

This document summarizes the structure, conventions, and key components of the Go-based roguelike project using the `gruid` library.

## Project Overview

*   **Language:** Go
*   **Core Framework:** `gruid` for terminal UI and grid management.
*   **Genre:** Roguelike game.
*   **Architecture:** Entity-Component-System (ECS).
*   **Key Features:** Turn-based movement, Field of View (FOV), map generation (rooms and tunnels), basic monster AI, player control.

## Core Packages and Responsibilities

*   **`internal/config`**:
    *   Handles command-line flag parsing (`config.go`).
    *   Defines game constants like map dimensions and FOV radius (`constants.go`).
*   **`internal/ecs`**:
    *   Implements a custom ECS (`ecs.go`).
        *   Entities are represented by `EntityID` (int).
        *   Components are stored in maps keyed by `EntityID` (e.g., `Positions`, `Renderables`, `TurnActors`).
        *   Provides methods for adding/removing entities and components, querying entities based on components.
        *   Uses `sync.RWMutex` for concurrent access safety.
    *   Defines component types (`components/components.go`, `components/fov.go`, `components/tags.go`, `components/turn.go`).
*   **`internal/game`**:
    *   Contains the main game logic and state.
    *   `game.go`: Defines the `Game` struct holding map, ECS, player ID, turn queue, etc. Handles level initialization.
    *   `model.go`: Defines the top-level `Model` struct for `gruid`, embedding the `Game` state.
    *   `model_update.go`: Handles input processing (`gruid.Msg`) and updates the game state.
    *   `actions.go`: Defines player actions and how they affect the game state.
    *   `map.go`: Implements map generation (using `rl.Grid`), stores map data, handles walkability and opacity checks, manages explored tiles (bitset).
    *   `player.go`: Handles player-specific actions like bumping into things.
    *   `spawner.go`: Functions for creating player and monster entities within the ECS.
    *   `systems.go`: Implements game systems like `RenderSystem`, `FOVSystem`, and turn processing logic (`processTurnQueue`, `monstersTurn`).
    *   `turn.go`: Defines game action interfaces and specific actions like `MoveAction`.
    *   `rooms.go`: Helper functions for creating rooms and tunnels during map generation.
    *   `log.go`: Game message logging system.
*   **`internal/turn_queue`**:
    *   Implements a turn queue using a min-heap (`heap.go`, `queue.go`).
    *   Manages entity turn order based on time/speed.
    *   Includes logic for adding entities, retrieving the next entity to act, and cleaning up dead entities.
*   **`internal/ui`**:
    *   Handles UI setup and drawing primitives.
    *   Defines color constants and styles (`color.go`).
    *   Initializes the `gruid` driver (`init.go`, `tcell.go`, `js.go`).
    *   Potentially uses tile-based rendering (`tiles.go`).
*   **`internal/utils`**:
    *   Provides utility functions.
    *   `assert.go`: Custom assertion functions (`Assert`, `Assertf`).
    *   `fov.go`: FOV calculation helpers (potentially using `rl` library).
    *   `signal.go`: Handles OS signals for graceful shutdown.

## Key Libraries & Dependencies

*   **`github.com/anaseto/gruid`**: Core terminal grid UI library.
*   **`github.com/disintegration/rl`**: Roguelike library, likely used for FOV (`rl.NewFOV`) and map grid (`rl.Grid`).
*   **`github.com/sirupsen/logrus`**: Structured logging.
*   Standard Go libraries: `container/list`, `container/heap`, `sync`, `math/rand`, `flag`, `fmt`, `time`, `os`, `context`, etc.

## Conventions & Patterns

*   **ECS:** Central architectural pattern.
*   **Component Storage:** Maps keyed by `EntityID`.
*   **Game Loop:** Driven by `gruid` messages and the `TurnQueue`. Input -> Action -> System Processing -> Render.
*   **Map Representation:** `rl.Grid` for tile data, bitsets (`[]uint64`) for explored/visible status.
*   **Concurrency:** `sync.RWMutex` used in the ECS to protect shared component maps.
*   **Logging:** `logrus` used for application logging, with debug level controlled by config.
*   **Error Handling:** Standard Go `error` interface, `fmt.Errorf`.
*   **Assertions:** Custom `Assert` functions for internal consistency checks.
*   **Configuration:** Command-line flags via `flag` package.

## Potential Areas for Change (Common Roguelike Features)

*   **Combat System:** Adding health, attack, defense components; implementing attack actions and damage calculation.
*   **Inventory System:** Components for inventory, items; actions for picking up, dropping, using items.
*   **More Complex AI:** Pathfinding, different monster behaviors.
*   **Level Progression:** Stairs, multiple dungeon levels.
*   **UI Enhancements:** Status panels, message log display, inventory screens.
*   **Persistence:** Saving and loading game state.
