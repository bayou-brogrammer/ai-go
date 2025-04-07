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
- Player movement with arrow keys and collision detection is implemented.
- **Turn-based system with action costs is implemented:**
  - `TurnQueue` (using `container/heap`) manages actor turns based on scheduled time.
  - `GameAction` interface defines actions with associated costs.
  - `processTurnQueue` handles the main turn loop, processing monster turns automatically.
  - Player input is handled correctly within the turn structure.
  - Action costs are applied when actors are re-queued.
  - Basic monster AI (`handleMonsterTurn`) exists and integrates with the turn system.
- **Field of View (FOV) is implemented:**
  - Uses `gruid/fov` package for FOV calculations.
  - Map tracks visible and explored tiles using efficient bitsets.
  - Rendering only shows visible and explored areas with appropriate colors.
  - FOV is updated when player moves.
- **Comprehensive color system is implemented:**
  - Semantic color definitions in `ui/color.go`.
  - Different colors for entities, map features, and UI elements.
  - Support for different rendering modes (terminal, SDL/JS).
  - Special colors for different monster types and states.
- **Basic combat system is being implemented:**
  - `Health` component added to track entity health.
  - `AttackAction` implemented for basic combat.
  - Bump-to-attack mechanic working for player and monsters.

## 4. Next Steps / Backlog

_(Ordered roughly by priority/dependency)_

1.  **Implement Turn-Based System with Action Costs:**
    - [x] Define `AITag` component (`components.go`) and add ECS methods (`ecs.go`).
    - [x] Modify `GameAction` interface in `turn.go` to return `(cost uint, err error)`.
    - [x] Update `MoveAction.Execute` in `turn.go` to return a cost (e.g., 100) and handle bump/fail cases.
    - [x] Update monster spawning logic (e.g., in `game.go` or map generation) to add `AITag` and add monsters to `turnQueue`.
    - [x] Create `handleMonsterTurn` function (`systems.go` or `ai.go`) for random movement and action cost calculation.
    - [x] Create `processTurnQueue` method in `model_update.go` to handle the main turn loop logic (processing monsters until player).
    - [x] Modify `model.Update` in `model_update.go` to call `processTurnQueue` and handle player input only when indicated.
    - [x] Modify `EndTurn` (or related logic) in `model_update.go` to use action cost for re-queueing the player and call `processTurnQueue` afterwards.
2.  **Implement Field of View (FOV):**
    - [x] Integrate `gruid/fov` package.
    - [x] Add `Visible` and `Explored` tracking to `Map` struct (partially done).
    - [x] Implement `ComputeFOV` function (likely in `map.go` or `fov.go`).
    - [x] Call `ComputeFOV` when player moves/acts.
    - [x] Modify `DrawMap` and `RenderSystem` to only draw visible/explored tiles/entities based on `Map.Visible` and `Map.Explored`.
3. **Refine Colors:**
   - [x] Define color constants (in `ui/color.go`).
   - [x] Update `DrawMap` and `RenderSystem` to use defined colors.
   - [x] Add different colors for explored-but-not-visible tiles.
   - [x] Implement color support for different monster types.
   - [x] Create coherent color palette with semantic naming.
4.  **Implement Basic Combat:**
   - [x] Define combat-related components (e.g., `Health`, `CombatStats`).
   - [x] Implement `AttackAction` for bumping into monsters.
   - [x] Create damage calculation logic (basic implementation).
   - [ ] Handle entity death.
5.  **Add Game UI:**
   - [ ] Create message log for combat and game events.
   - [ ] Add player stats display.
   - [ ] Add health indicators.
6.  **Enhance Monster System:**
   - [x] Define other monster components (e.g., `Health`, `CombatStats`).
   - [x] Refine monster spawning locations/types with varying stats.
   - [ ] Improve AI behavior beyond random movement.

## 5. Known Issues / TODOs

- ~FOV not implemented.~ ✓
- ~Colors are placeholders.~ ✓
- Monster AI is very basic (random movement).
- ~No combat implemented.~ ✓ (Basic implementation complete, needs death handling)
- No UI elements for feedback.
- `driver` variable in `main.go` likely needs proper initialization.
- Missing unit tests.
