# Progress: Roguelike Gruid

## What Works

*   **Basic Project Structure:** ECS, Game, UI, Config, Utils packages are set up.
*   **ECS Core:** Adding/removing entities and components, querying entities.
*   **Map Generation:** Basic procedural generation of rooms and tunnels using `rl.Grid`. Walls and floors are distinct.
*   **Player Entity:** Player ('@') is spawned at a valid starting position.
*   **Basic Monster Spawning:** Monsters ('o', 'T', etc.) are placed randomly in rooms.
*   **Turn Queue:** Entities with `TurnActor` component are added to a turn queue (`internal/turn_queue`). Basic turn processing logic exists.
*   **Player Movement:** Player can move North, South, East, West using arrow keys or vi keys. Movement respects map boundaries and wall collisions.
*   **Field of View (FOV):** Player has FOV calculated using `rl.NewFOV`. Visible tiles are tracked.
*   **Map Exploration:** Explored tiles are remembered and displayed differently from currently visible tiles and unseen tiles.
*   **Rendering:** The map, player, and monsters are rendered to the terminal using `gruid`. Visible/Explored status affects tile appearance.
*   **Basic Monster Turn:** Placeholder logic exists for monsters to take turns (likely needs significant expansion).
*   **Configuration:** Debug logging can be enabled via command-line flags (`-debug`).
*   **Graceful Shutdown:** Handles OS signals (Ctrl+C) for clean exit.

## What's Left to Build (Key Features)

*   **Combat System:** Health, attack, defense components; attack actions, damage calculation, entity death handling beyond just removing from queue.
*   **Meaningful Monster AI:** Pathfinding (e.g., A*), different behaviors (chasing, fleeing, idle). Currently, AI is likely very basic or placeholder.
*   **Interaction System:** Player bumping into monsters should trigger combat, not just block movement.
*   **Game State Management:** Win/loss conditions.
*   **UI Enhancements:** Message log display, status panel (HP, etc.), inventory screen.
*   **Inventory System:** Picking up, dropping, using items.
*   **Items & Equipment:** Defining items, effects, equipment slots.
*   **Level Progression:** Stairs, generating new levels, increasing difficulty.
*   **Persistence:** Saving and loading game state.

## Current Status

The project provides a solid foundation for a roguelike game. Core systems like ECS, map generation, FOV, turn management, and basic rendering are functional. The immediate next steps would likely involve implementing combat and more sophisticated AI.

## Known Issues

*   Monster AI is rudimentary.
*   Collision handling might only block movement, not trigger interactions (like attacks).
*   No health or damage system implemented yet.
*   Game log (`internal/game/log.go`) exists but might not be displayed or fully utilized.

## Evolution of Project Decisions

*   Initial decision to use `gruid` and `rl` libraries.
*   Choice of a custom ECS implementation.
*   Adoption of a turn-based system managed by a min-heap queue.
