# Product Context: Roguelike Gruid

## Problem Solved

This project aims to create a classic terminal-based roguelike game experience, providing a foundation for learning and experimenting with game development concepts in Go, specifically using the `gruid` library and ECS patterns.

## How it Should Work

*   The player controls a character ('@') on a grid-based map.
*   The map is procedurally generated with rooms and tunnels.
*   The player can move in four cardinal directions using keyboard input.
*   The player has a limited Field of View (FOV); unexplored areas are hidden, previously seen areas are remembered but dimmed.
*   Monsters ('o', 'T', etc.) exist on the map and take turns acting.
*   Basic monster AI involves moving randomly or towards the player if visible (future goal).
*   The game operates in turns, processing player and monster actions sequentially.
*   The game runs within a terminal window.

## User Experience Goals

*   Clear and responsive terminal interface.
*   Intuitive keyboard controls for movement.
*   Classic roguelike feel regarding exploration and discovery.
