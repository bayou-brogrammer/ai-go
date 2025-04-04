# Architecture: Go Roguelike with Gruid

## 1. Overview

This project implements a simple roguelike game using the Go programming language and the `gruid` terminal UI library. The core game state management is handled by a custom-built Entity-Component-System (ECS).

## 2. Key Components

- **`main.go`**: Entry point of the application. Initializes the `gruid` app, the game state (`game`), the ECS `World`, the `model`, and creates the initial player entity.
- **`model.go`**: Implements the `gruid.Model` interface. Handles the main game loop (`Update` for input/logic, `Draw` for rendering). Owns the drawing `grid` and a reference to the `game` state.
- **`game.go`**: Contains the main `game` struct, holding references to the `Map`, the ECS `World`, the `PlayerID`, and the `TurnQueue`.
- **`map.go`**: Defines the `Map` struct and related functions.
- **`ecs.go`**: Defines the core ECS `World` struct and methods.
- **`components.go`**: Defines component structs (e.g., `Position`, `Renderable`, `BlocksMovement`, `AITag`).
- **`systems.go`**: Contains system functions (e.g., `RenderSystem`, `handleMonsterTurn`).
- **`turn.go`**: Defines the `TurnQueue` (priority queue based on `container/heap`), `TurnEntry`, `GameAction` interface, and specific actions like `MoveAction`.
- **`driver.go` / `sdl.go` / `tcell.go`**: Handles the underlying terminal driver.

## 3. Core Design: Entity-Component-System (ECS)

- **Entity:** Represented by an `EntityID` (int). A simple identifier.
- **Component:** Plain Go structs holding data (e.g., `Position{ Point gruid.Point }`).
- **System:** Functions (like `RenderSystem`) that contain game logic. They query the `World` for entities with specific components and update the game state or render output.
- **World:** The central manager (`ecs.World`) holding the component store (`map[EntityID]map[reflect.Type]interface{}`). Provides methods for entity/component manipulation and querying.

## 4. Data Flow (Turn Processing Example)

1. `main` initializes `gruid`, `model`, `game` (including `ecs.World`, `Map`, `TurnQueue`).
2. Entities (Player, Monsters) are created and added to `ecs.World` with components (`Position`, `Renderable`, `AITag` for monsters).
3. Entities are added to `game.TurnQueue` with initial time (e.g., 0).
4. `gruid` calls `model.Update(msg)`.
5. `model.Update` calls `processTurnQueue()`.
6. `processTurnQueue` loops:
    a.  Peeks at `TurnQueue`.
    b.  If Monster is next and ready (`entry.Time <= CurrentTime`):
        i.  Pops monster from `TurnQueue`.
        ii. Updates `TurnQueue.CurrentTime`.
        iii. Calls `handleMonsterTurn(game, monsterID)`.
        iv. `handleMonsterTurn` determines action (e.g., `MoveAction`), calls `action.Execute(game)`.
        v.  `Execute` performs logic (e.g., updates position via ECS) and returns `cost`.
        vi. `handleMonsterTurn` returns `cost`.
        vii. `processTurnQueue` adds monster back to `TurnQueue` with `CurrentTime + cost`.
        viii. Continues loop.
    c.  If Player is next and ready:
        i.  Sets `waitForPlayerInput = true`.
        ii. Breaks loop.
    d.  If queue empty or next actor not ready:
        i.  Sets `waitForPlayerInput = false`.
        ii. Breaks loop.
7. If `waitForPlayerInput` is true:
    a.  `model.Update` processes `gruid.Msg` (e.g., `MsgKeyDown`).
    b.  Input handler determines `playerAction` (e.g., `MoveAction`).
    c.  Calls `playerAction.Execute(game)`, gets `cost`.
    d.  Calls `EndTurn()` (or similar).
    e.  `EndTurn` updates `TurnQueue.CurrentTime` (if needed) and adds Player back to `TurnQueue` with `CurrentTime + cost`.
    f.  `EndTurn` calls `processTurnQueue()` again to handle immediate reactions.
8. `model.Update` returns `gruid.Effect` (e.g., `gruid.Flush{}`).
9. `gruid` calls `model.Draw()` (which uses `RenderSystem` based on updated ECS state).

## 5. Visualization

```mermaid
---
config:
  layout: fixed
---
flowchart LR
    subgraph Gruid["Gruid Framework"]
        App["App Loop"]
        InputMsg["Input Msgs"]
        OutputGrid["Output Grid"]
    end

    subgraph ModelLayer["Model Layer (model.go)"]
        Model["Model"]
        M_Update["Update()"]
        M_Draw["Draw()"]
        Grid["gruid.Grid"]
        M_ProcessTQ["processTurnQueue()"]
        M_EndTurn["EndTurn()"]
    end

    subgraph GameState["Game State (game.go)"]
        Game["Game"]
        ECS["ECS World (ecs.go)"]
        Map["Map (map.go)"]
        TQ["TurnQueue (turn.go)"]
        PlayerID["PlayerID"]
    end

    subgraph TurnLogic["Turn Logic (turn.go, systems.go)"]
        TurnQueue["TurnQueue (Heap)"]
        GameAction["GameAction Interface"]
        MoveAction["MoveAction"]
        HMT["handleMonsterTurn()"]
    end

    subgraph ECSSubsystem["ECS Subsystem"]
        World["World (ecs.go)"]
        Components["Components (components.go)<br>Position, Renderable,<br>AITag, ..."]
        Systems["Systems (systems.go)<br>RenderSystem"]
    end

    App -- Calls --> M_Update & M_Draw
    M_Update -- Receives --> InputMsg
    M_Draw -- Returns --> OutputGrid
    M_Draw -- Uses --> Grid

    Model -- Owns --> Game & Grid
    Game -- Owns --> ECS & Map & TQ & PlayerID

    M_Update -- Calls --> M_ProcessTQ
    M_ProcessTQ -- Interacts --> TQ
    M_ProcessTQ -- Calls --> HMT

    M_Update -- Processes --> InputMsg
    M_Update -- Calls --> GameAction -- via EndTurn --> M_EndTurn

    M_EndTurn -- Interacts --> TQ
    M_EndTurn -- Calls --> M_ProcessTQ

    HMT -- Determines --> GameAction
    HMT -- Calls --> GameAction
    GameAction -- Interacts --> Game
    GameAction -- Modifies --> ECS
    GameAction -- Returns --> Cost

    TurnQueue -- Stores --> TurnEntry(EntityID, Time)
    TQ -- Uses --> TurnQueue

    M_Draw -- Calls --> Systems
    Systems -- Queries --> World
    Systems -- Updates --> Grid

    World -- Manages --> Components
