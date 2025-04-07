# Tech Context: Roguelike Gruid

## Technologies Used

*   **Programming Language:** Go (Golang)
*   **Core UI Library:** `github.com/anaseto/gruid`
*   **Roguelike Utilities:** `github.com/disintegration/rl` (primarily for FOV, potentially Grid)
*   **Logging:** `github.com/sirupsen/logrus`
*   **Concurrency:** Go standard library `sync` package (`sync.RWMutex` in ECS).
*   **Data Structures:** Go standard library `container/heap`, `container/list`.
*   **Configuration:** Go standard library `flag` package.

## Development Setup

*   Standard Go development environment (Go compiler, modules).
*   Dependencies managed via Go modules (`go.mod`, `go.sum`).
*   Build and run using standard `go build` and `go run`.

## Technical Constraints

*   Terminal-based UI limitations (character grid, limited colors/styles).
*   Performance considerations for FOV calculation and rendering on large maps or with many entities (addressed partly by bitsets and efficient ECS queries).

## Dependencies

*   `github.com/anaseto/gruid`
*   `github.com/disintegration/rl`
*   `github.com/sirupsen/logrus`
*   `github.com/gdamore/tcell/v2` (likely a dependency of `gruid` for terminal interaction)

## Tool Usage Patterns

*   **Go Modules:** For dependency management (`go get`, `go mod tidy`).
*   **Go Build/Run:** For compiling and executing the application.
*   **Git:** For version control.
*   **Logrus:** For structured logging, level controlled by `-debug` flag.
