# Active Context: Go Roguelike with Gruid (As of 2025-04-06)

## Current Focus

- **Implement Basic Combat System:**
  - Define combat-related components (`Health`, `CombatStats`)
  - Create an `AttackAction`
  - Implement combat resolution logic
  - Handle player attacking monsters and vice versa

## Recent Changes / Activity

- Implemented comprehensive color system:
  - Created base palette and semantic color mapping in `ui/color.go`
  - Added support for different entity colors and states
  - Updated renderer to use the color system in both terminal and SDL/JS modes
  - Added support for explored-but-not-visible tile coloring
- Completed Field of View (FOV) implementation task
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
- Turn-based system with action costs is functional. Actors take turns based on a priority queue (`TurnQueue`), and actions have associated time costs. Basic monster AI exists.
- **Field of View (FOV) is now implemented.** Player can only see tiles and entities within their field of view. Tiles outside FOV but previously seen appear dimmed.
- **Comprehensive color system is implemented.** Colors are now defined semantically and mapped to appropriate terminal/SDL colors.
- Project lacks combat, advanced monster AI, and UI elements.
- Core documentation files are being updated to reflect the current state.

## Next Steps (Immediate)

1. **Implement Basic Combat:**
   - Define combat-related components (e.g., `Health`, `CombatStats`).
   - Implement `AttackAction` for bumping into monsters.
   - Create damage calculation logic.
   - Handle entity death.
2. **Add Game UI:**
   - Create message log for combat and game events.
   - Add player stats display.
   - Add health indicators.
