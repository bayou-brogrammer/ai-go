# Active Context: Roguelike Gruid - Initial Setup

## Current Focus

*   Establishing the initial project structure and core components based on the provided analysis.
*   Setting up the Memory Bank for ongoing development.

## Recent Changes

*   Initial project analysis performed.
*   Memory Bank files generated.

## Next Steps

*   Begin implementing specific features or addressing user requests based on the established codebase.

## Active Decisions & Considerations

*   Adhering to the existing ECS pattern.
*   Leveraging `gruid` for UI and `rl` for roguelike utilities (FOV, map grid).
*   Following Go best practices and existing project conventions.

## Important Patterns & Preferences

*   Use the custom ECS implementation found in `internal/ecs`.
*   Manage turns using the `internal/turn_queue` package.
*   Generate maps using the logic in `internal/game/map.go`.
*   Handle configuration via `internal/config`.
*   Use `logrus` for logging.

## Learnings & Insights

*   The project has a well-defined structure based on ECS principles.
*   Key libraries (`gruid`, `rl`) are already integrated.
*   Core systems (FOV, turn management, basic map generation) are present.
