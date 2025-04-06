# Active Context: Go Roguelike with Gruid (As of 2025-04-05)

## Current Focus

- **Implement Field of View (FOV):** (Next major task - see `tasks/tasks_plan.md`)
  - Integrate `gruid/fov` package.
  - Add `Visible` and `Explored` tracking to `Map` struct.
  - Implement `ComputeFOV` function.
  - Call `ComputeFOV` when player moves/acts.
  - Modify rendering to respect FOV.

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
- Rendering pipeline draws map tiles and entities.
- Basic map generation with rooms/corridors is implemented.
- Player starts in the center of the first room.
- Player can move around using arrow keys with wall collision detection.
- **Turn-based system with action costs is functional.** Actors take turns based on a priority queue (`TurnQueue`), and actions have associated time costs. Basic monster AI exists.
- Project lacks FOV, refined colors, advanced monster AI, and combat.
- Core documentation files are being updated to reflect the current state.

## Next Steps (Immediate)

1.  **Implement Field of View (FOV):** (See `tasks/tasks_plan.md` for details)
    - Integrate `gruid/fov`.
    - Update `Map` struct.
    - Implement `ComputeFOV`.
    - Update rendering logic.
2.  **Refine Colors:**
    - Define color constants.
    - Update rendering to use constants.
    - Consider explored-but-not-visible color.
