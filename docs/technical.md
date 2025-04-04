# Technical Details: Go Roguelike with Gruid

## 1. Technology Stack

- **Language:** Go (Golang)
- **Terminal UI Library:** `gruid` (`codeberg.org/anaseto/gruid`)
  - Provides grid-based rendering, input handling, and application structure.
- **Terminal Driver:** TBD (Likely `gruid-sdl` or `gruid-tcell`, needs confirmation/setup in `driver.go`). Currently placeholder `driver` variable used in `main.go`.

## 2. Development Environment

- **Operating System:** macOS Sequoia (as per SYSTEM INFORMATION)
- **IDE/Editor:** VS Code (implied by environment details)
- **Build/Run:** Standard Go toolchain (`go build`, `go run`). Project is within `roguelike-gruid-project` subdirectory.

## 3. Key Technical Decisions & Patterns

- **Architecture:** Entity-Component-System (ECS)
  - **Implementation:** Custom-built (in `ecs.go`).
  - **Entity ID:** `int` (`ecs.EntityID`).
  - **Component Storage:** Central map in `ecs.World`: `map[EntityID]map[reflect.Type]interface{}`.
  - **System Logic:** Implemented as functions (e.g., `RenderSystem`, `handleMonsterTurn` in `systems.go`) that query the `ecs.World`.
- **Turn Management:** (`turn.go`)
  - **Queue:** `TurnQueue` struct using `container/heap` to manage `TurnEntry` items based on scheduled time (`uint64`). Provides a priority queue for actor turns.
  - **Actions:** `GameAction` interface defines actions actors can take. `Execute` method performs the action and returns an action `cost` (`uint`). Specific actions like `MoveAction` implement this.
  - **AI:** Monsters identified by `AITag` component (`components.go`). Basic AI logic resides in `handleMonsterTurn` (`systems.go`).
- **Rendering:** (`model.go`, `map.go`, `systems.go`)
  - Managed by `gruid` via `model.Draw`. Map drawn first, then entities via `RenderSystem`.
- **Game Loop:** (`model.go`) Handled by `gruid` calling `model.Update` and `model.Draw`. `model.Update` now incorporates `processTurnQueue` to manage actor turns based on the `TurnQueue` before handling player input.

## 4. Project Structure

(See `docs/architecture.md` for component relationships)

- `roguelike-gruid-project/`: Main project directory.
  - `main.go`: Application entry point, initialization.
  - `model.go`: `gruid` model implementation (Update/Draw loop).
  - `game.go`: High-level game state struct.
  - `map.go`: Map definition and drawing logic.
  - `ecs.go`: Core ECS implementation.
  - `components.go`: Component struct definitions (includes `AITag`).
  - `systems.go`: System function implementations (includes `RenderSystem`, `handleMonsterTurn`).
  - `turn.go`: Turn queue, actions, action costs.
  - `color.go`: (Likely for color definitions - TBD).
  - `tiles.go`: (Likely for tile definitions/constants - TBD).
  - `driver.go`/`sdl.go`/`tcell.go`: Terminal driver setup (TBD).
  - `go.mod`, `go.sum`: Go module files.

## 5. Dependencies

- `codeberg.org/anaseto/gruid`: Core terminal UI library.
- (Potentially `gruid-sdl` or `gruid-tcell` for the driver).

## 6. Technical Constraints/Considerations

- Performance of the `reflect`-based component storage might need evaluation later for larger numbers of entities/components.
- Terminal capabilities (color support, rune support) depend on the chosen driver and the user's terminal emulator.
