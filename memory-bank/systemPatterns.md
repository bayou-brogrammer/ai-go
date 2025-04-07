# System Patterns: Roguelike Gruid

## Architecture

*   **Entity-Component-System (ECS):** The core architectural pattern.
    *   Implemented in `internal/ecs`.
    *   Entities are `EntityID` (int).
    *   Components are stored in maps within the `ECS` struct (e.g., `Positions`, `Renderables`).
    *   Systems (like `RenderSystem`, `FOVSystem`) operate on entities possessing specific component combinations.
*   **Event/Turn-Based Loop:** Driven by `gruid` input messages and the `internal/turn_queue`.
    *   Input -> Action -> System Processing -> Render.
*   **Modular Packages:** Code is organized into packages based on responsibility (`ecs`, `game`, `ui`, `config`, `turn_queue`, `utils`).

## Key Technical Decisions

*   **Custom ECS Implementation:** The project uses a bespoke ECS rather than a third-party library.
*   **`gruid` Library:** Used for terminal grid management, input handling, and rendering.
*   **`rl` Library:** Used for specific roguelike algorithms, notably FOV calculation (`rl.NewFOV`) and potentially map grid representation (`rl.Grid`).
*   **Min-Heap Turn Queue:** `internal/turn_queue` uses `container/heap` for efficient turn ordering based on entity speed/time.
*   **Bitsets for Map State:** `[]uint64` bitsets are used for tracking explored and visible map tiles for efficiency.

## Component Relationships

*   Entities gain capabilities by having components added (e.g., adding `Position` and `Renderable` makes an entity appear on the map).
*   Systems query the ECS for entities with specific sets of components to perform updates (e.g., `FOVSystem` needs entities with `Position` and `FOV`).

## Critical Implementation Paths

*   **Game Loop:** `model_update.go` processes `gruid.Msg`, potentially queues actions, `processTurnQueue` in `systems.go` advances game time and executes actor turns.
*   **Rendering:** `model_draw.go` calls `RenderSystem` which queries ECS and draws entities onto the `gruid.Grid`. `FOVSystem` updates visibility first.
*   **Map Generation:** `internal/game/map.go` uses room/tunnel generation algorithms (`internal/game/rooms.go`) to populate the `rl.Grid`.
*   **Player Action:** Input maps to actions (`internal/game/input.go`), actions are processed (e.g., `EntityBump` in `player.go`), potentially modifying ECS state.
*   **Monster Action:** `monstersTurn` in `systems.go` iterates AI entities, determines actions (currently basic/placeholder), and queues them.
