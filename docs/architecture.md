# Architecture: Go Roguelike with Gruid

## 1. Overview

This project implements a simple roguelike game using the Go programming language and the `gruid` terminal UI library. The core game state management is handled by a custom-built Entity-Component-System (ECS).

## 2. Key Components

- **`main.go`**: Entry point of the application. Initializes the `gruid` app, the game state (`game`), the ECS `World`, the `model`, and creates the initial player entity.
- **`model.go`**: Implements the `gruid.Model` interface. Handles the main game loop (`Update` for input/logic, `Draw` for rendering). Owns the drawing `grid` and a reference to the `game` state.
- **`game.go`**: Contains the main `game` struct, holding references to the `Map` and the ECS `World`. (May hold other high-level game state later).
- **`map.go`**: Defines the `Map` struct (holding `LogicalTiles`, visibility, explored status) and related functions like `DrawMap`, `InBounds`, `IsWall`.
- **`ecs.go`**: Defines the core ECS `World` struct and its methods (`CreateEntity`, `AddComponent`, `GetComponent`, `QueryEntitiesWithComponents`, etc.). Manages entities and their components.
- **`components.go`**: Defines the component structs (e.g., `Position`, `Renderable`, `BlocksMovement`). These are simple data containers.
- **`systems.go`**: Contains system functions (e.g., `RenderSystem`) that operate on entities possessing specific sets of components. Systems query the `World` for relevant entities and perform actions.
- **`driver.go` / `sdl.go` / `tcell.go`**: (Assumed based on typical `gruid` setup) Handles the underlying terminal driver (SDL or Tcell).

## 3. Core Design: Entity-Component-System (ECS)

- **Entity:** Represented by an `EntityID` (int). A simple identifier.
- **Component:** Plain Go structs holding data (e.g., `Position{ Point gruid.Point }`).
- **System:** Functions (like `RenderSystem`) that contain game logic. They query the `World` for entities with specific components and update the game state or render output.
- **World:** The central manager (`ecs.World`) holding the component store (`map[EntityID]map[reflect.Type]interface{}`). Provides methods for entity/component manipulation and querying.

## 4. Data Flow (Rendering Example)

1.  `main` initializes `gruid`, `model`, `game`, `ecs.World`, `Map`.
2.  `main` creates the player `EntityID` and adds `Position`, `Renderable` components via `ecs.World`.
3.  `gruid` calls `model.Draw()`.
4.  `model.Draw()` clears the `gruid.Grid`.
5.  `model.Draw()` calls `DrawMap(game.Map, grid)` to render map tiles.
6.  `model.Draw()` calls `RenderSystem(game.ecs, grid)`.
7.  `RenderSystem` queries `ecs.World` for entities with `Position` and `Renderable`.
8.  `RenderSystem` iterates through results, gets component data, and sets cells on the `gruid.Grid`.
9.  `model.Draw()` returns the updated `gruid.Grid` to `gruid` for display.

## 5. Visualization

```mermaid
---
config:
  layout: fixed
---
flowchart LR
 subgraph subGraph0["ECS Core"]
        C_Store["(Component Store<br>map[EID]map[Type]Comp)"]
        World["World (ecs.go)"]
        M_Create["CreateEntity"]
        M_AddComp["AddComponent"]
        M_GetComp["GetComponent"]
        M_Destroy["DestroyEntity"]
        M_Query["QueryEntities"]
  end
 subgraph subGraph1["Game Logic"]
        ECS["World Ref"]
        Game["Game (game.go)"]
        Map["Map (map.go)"]
        Model["Model (model.go)"]
        Grid["gruid.Grid"]
  end
 subgraph Systems["Systems (systems.go)"]
        RenderSystem["RenderSystem"]
        MovementSystem["(Future) MovementSystem"]
        AISystem["(Future) AISystem"]
  end
 subgraph Components["Components (components.go)"]
        Position["Position"]
        Renderable["Renderable"]
        BlocksMovement["BlocksMovement"]
  end
    World --- C_Store & M_Create & M_AddComp & M_GetComp & M_Destroy & M_Query
    Game --- ECS & Map
    Model --- Game & Grid
    RenderSystem -- Queries --> World
    MovementSystem -- Queries --> World
    AISystem -- Queries --> World
    Position -- Added to --> World
    Renderable -- Added to --> World
    BlocksMovement -- Added to --> World
    Model -- Calls --> RenderSystem & MovementSystem & AISystem
    Model -- Calls --> DrawMap
    DrawMap -- Reads --> Map
    DrawMap -- Updates --> Grid
    RenderSystem -- Updates --> Grid
