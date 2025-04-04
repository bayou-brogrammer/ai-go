# Active Context: Go Roguelike with Gruid (As of 2025-04-04)

## Current Focus

- **Implement Turn-Based System with Action Costs:** (New Task - see `tasks/tasks_plan.md`)
  - Integrate `TurnQueue` into the main game loop (`model.Update`).
  - Process monster turns automatically based on the queue.
  - Pause processing when it's the player's turn to wait for input.
  - Implement action costs returned by `GameAction.Execute`.
  - Add `AITag` component and basic random monster AI handler.

## Recent Changes / Activity

- Created `docs/product_requirement_docs.md`.
- Created `docs/architecture.md`.
- Created `docs/technical.md`.
- Created `tasks/tasks_plan.md`.
- Planned the implementation of the Turn-Based System (2025-04-04).
- Added `DrawMap` function to `map.go` (using `gruid.ColorDefault` as placeholder).
- Integrated `DrawMap` call into `model.Draw` in `model.go`.
- Implemented basic map generation (`generateMap` in `map.go`).
- Updated `NewMap` to use `generateMap` and return player start position.
- Updated `main.go` to use the returned player start position.
- Created placeholder `.cursor/rules/error-documentation.mdc`.
- Created placeholder `.cursor/rules/lessons-learned.mdc`.
- Implemented player movement with arrow keys and collision detection (`model.Update`, `actions.go`).
- Added `TurnQueue` using `container/heap` (`turn.go`).
- Added `GameAction` interface and `MoveAction` (`turn.go`).

## Current State Summary

- Core ECS implementation is in place.
- Player entity is created with basic components.
- Rendering pipeline draws map tiles (currently static) and then entities.
- Basic map generation with rooms/corridors is implemented.
- Player starts in the center of the first room.
- Player can move around using arrow keys with wall collision detection.
- Project lacks FOV, colors, active monsters/AI, combat, and a proper turn-based game loop integrated with action costs.
- Core documentation files (`product_requirement_docs.md`, `architecture.md`, `technical.md`, `tasks_plan.md`, `active_context.md`) are being updated.

## Next Steps (Immediate)

1. **Implement Turn-Based System** (as per the plan):
   - Define `AITag` component and ECS methods.
   - Modify `GameAction` interface and `MoveAction` to return costs.
   - Ensure monsters are added to `TurnQueue` on spawn with `AITag`.
   - Create `handleMonsterTurn` function.
   - Implement `processTurnQueue` logic in `model_update.go`.
   - Update `EndTurn` logic for player action costs and re-queueing.
   - Call `processTurnQueue` after player turn.
2. Implement Field of View (FOV).
3. Refine Colors.
