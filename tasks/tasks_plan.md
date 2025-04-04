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

**Phase 3: Systems & Game Loop**

- [x] Implement Initial Render System (`systems.go`: `RenderSystem` created)
- [x] Modify Game Loop (`model.go`):
  - [x] `model.Draw` clears grid.
  - [x] `model.Draw` calls `DrawMap` (added).
  - [x] `model.Draw` calls `RenderSystem`.
  - [x] `model.Update` handles input and movement.

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
- Player movement with arrow keys and collision detection is implemented (`model.Update`, `actions.go`).
- Basic `TurnQueue` using `container/heap` is implemented (`turn.go`).
- `GameAction` interface and `MoveAction` struct defined (`turn.go`).

## 4. Next Steps / Backlog

_(Ordered roughly by priority/dependency)_

1.  **Implement Turn-Based System with Action Costs:**
    - [ ] Define `AITag` component (`components.go`) and add ECS methods (`ecs.go`).
    - [ ] Modify `GameAction` interface in `turn.go` to return `(cost uint, err error)`.
    - [ ] Update `MoveAction.Execute` in `turn.go` to return a cost (e.g., 100) and handle bump/fail cases.
    - [ ] Update monster spawning logic (e.g., in `game.go` or map generation) to add `AITag` and add monsters to `turnQueue`.
    - [ ] Create `handleMonsterTurn` function (`systems.go` or `ai.go`) for random movement and action cost calculation.
    - [ ] Create `processTurnQueue` method in `model_update.go` to handle the main turn loop logic (processing monsters until player).
    - [ ] Modify `model.Update` in `model_update.go` to call `processTurnQueue` and handle player input only when indicated.
    - [ ] Modify `EndTurn` (or related logic) in `model_update.go` to use action cost for re-queueing the player and call `processTurnQueue` afterwards.
2.  **Implement Field of View (FOV):**
    - [ ] Integrate `gruid/fov` package.
    - [ ] Add `Visible` and `Explored` tracking to `Map` struct (partially done).
    - [ ] Implement `ComputeFOV` function (likely in `map.go` or `fov.go`).
    - [ ] Call `ComputeFOV` when player moves/acts.
    - [ ] Modify `DrawMap` and `RenderSystem` to only draw visible/explored tiles/entities based on `Map.Visible` and `Map.Explored`.
2. **Refine Colors:**
   - Define color constants (e.g., in `color.go`).
   - Update `DrawMap` and `RenderSystem` (and player creation) to use defined colors instead of `gruid.ColorDefault`.
   - [ ] Consider different colors for explored-but-not-visible tiles.
   4.  **Add Basic Monsters (Initial setup done in Turn System task):**
       - [ ] Define other monster components (e.g., `Health`, `CombatStats`).
       - [ ] Refine monster spawning locations/types.
       - [ ] Update `RenderSystem` to draw monsters correctly (if not already covered).
5.  **Implement Basic Combat:**
    - [ ] Define combat-related components/events (e.g., `AttackAction`, `Health`, `CombatStats`).
    - [ ] Implement `CombatSystem` or integrate logic into action handlers.
    - [ ] Handle player attacking monsters and vice-versa within the turn system.
6.  **Game UI:**
    - [ ] Message log.
    - [ ] Player stats display.

## 5. Known Issues / TODOs

- FOV not implemented.
- Colors are placeholders.
- Monster AI is very basic (random movement).
- No combat implemented.
- `driver` variable in `main.go` likely needs proper initialization.
- Missing unit tests.
