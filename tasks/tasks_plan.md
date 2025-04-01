# Task Plan: Go Roguelike with Gruid

## 1. Overall Goal

Create a basic, playable roguelike game in the terminal using Go and `gruid`, featuring an ECS architecture.

## 2. Completed Tasks (Based on ROGUELIKE_GUIDE_PLAN.md)

**Phase 1: Core ECS Structure**

- [x] Define Core Types (`ecs.go`: `EntityID`, `World`, component map)
- [x] Implement World Methods (`ecs.go`: `NewWorld`, `CreateEntity`, `AddComponent`, `GetComponent`, `QueryEntitiesWithComponents`, etc.)
- [x] Define Initial Components (`components.go`: `Position`, `Renderable`, `BlocksMovement`)

**Phase 2: Integration**

- [x] Integrate World into Game (`game.go`: Added `ecs *World` field)
- [x] Create Player Entity (`main.go`: Created player entity with `Position`, `Renderable`, `BlocksMovement` components)

**Phase 3: Systems & Game Loop (Partial)**

- [x] Implement Initial Render System (`systems.go`: `RenderSystem` created)
- [x] Modify Game Loop (`model.go`):
  - [x] `model.Draw` clears grid.
  - [x] `model.Draw` calls `DrawMap` (added).
  - [x] `model.Draw` calls `RenderSystem`.
  - [ ] `model.Update` is placeholder (needs input handling/system calls).

**Phase 4: Map Generation**

- [x] Define `Rect` struct and helpers (`map.go`)
- [x] Implement `generateMap` function (`map.go`)
- [x] Modify `NewMap` to use `generateMap` (`map.go`)
- [x] Modify `main` to handle `NewMap` return values and use `playerStart` (`main.go`)
- [x] Seed random number generator (`main.go`)

## 3. Current Status

- Core ECS is functional.
- Player entity is created and added to the ECS.
- Map tiles (`DrawMap`) and Entities (`RenderSystem`) are rendered in the correct order in `model.Draw`.
- Basic map generation with rooms and corridors is implemented.
- Player starts in the center of the first generated room.
- The game runs but has no interactivity yet (no player movement).

## 4. Next Steps / Backlog

_(Ordered roughly by priority/dependency)_

1. **Implement Basic Map Generation:**
   - Define room/corridor structures.
   - [x] Define room/corridor structures (`Rect` in `map.go`).
   - [x] Implement an algorithm (simple room placement + corridors) in `map.go` (`generateMap`).
   - [x] Modify `NewMap` to use the generation algorithm.
   - [x] Update player start position logic to place player in a valid floor tile (`generateMap` returns start, `main.go` uses it).
2. **Implement Player Movement:**
   - Handle keyboard input (arrow keys or hjkl) in `model.Update`.
   - Create `MovementSystem` (or add logic to `model.Update`).
   - Query for player entity's `Position`.
   - Check for collisions with walls (`Map.IsWall`) before updating `Position`.
3. **Implement Field of View (FOV):**
   - Integrate `gruid/fov` package.
   - Add `Visible` and `Explored` tracking to `Map` struct (already done).
   - Implement `ComputeFOV` function (likely in `map.go` or `fov.go`).
   - Call `ComputeFOV` when player moves.
   - Modify `DrawMap` and `RenderSystem` to only draw visible/explored tiles/entities based on `Map.Visible` and `Map.Explored`.
4. **Refine Colors:**
   - Define color constants (e.g., in `color.go`).
   - Update `DrawMap` and `RenderSystem` (and player creation) to use defined colors instead of `gruid.ColorDefault`.
   - Consider different colors for explored-but-not-visible tiles.
5. **Add Basic Monsters:**
   - Define monster components (e.g., `AI`, `Health`, `CombatStats`).
   - Create monster entities in `main.go` or map generation.
   - Implement basic AI system (e.g., `AISystem` for random movement).
   - Update `RenderSystem` to draw monsters.
   - Update `MovementSystem` to handle monster collisions.
6. **Implement Basic Combat:**
   - Define combat-related components/events.
   - Implement `CombatSystem`.
   - Handle player attacking monsters and vice-versa.
7. **Game UI:**
   - Message log.
   - Player stats display.

## 5. Known Issues / TODOs

- No player input handling or movement.
- No FOV implemented yet.
- Colors are placeholders (`gruid.ColorDefault`).
- `driver` variable in `main.go` is likely a placeholder for actual `gruid-sdl` or `gruid-tcell` driver initialization.
- Missing unit tests.
