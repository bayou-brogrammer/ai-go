# Codebase Memory Bank: Roguelike Gruid

## Entity-Component-System (ECS)

### Core ECS Structure
- Located in `internal/ecs/ecs.go`
- Main structure: `ECS` with maps for entities and components
- `EntityID` (int) used to identify entities
- Generic component storage via maps: `map[components.ComponentType]map[EntityID]any`

### Entity Management
- `AddEntity()` - Creates a new entity and returns its ID
- `RemoveEntity(id EntityID)` - Removes an entity and all its components
- `EntityExists(id EntityID)` - Checks if an entity exists

### Component Management
- Components defined in `internal/ecs/components/`
- `AddComponent(id EntityID, compType components.ComponentType, component any)`
- `RemoveComponent(id EntityID, compType components.ComponentType)`
- `HasComponent(id EntityID, compType components.ComponentType)`
- `getComponent(id EntityID, compType components.ComponentType)` - Generic retrieval
- `GetComponentTyped<T>()` - Type-safe retrieval

### Query Functions
- `GetEntitiesWithComponent(compType ComponentType)` - Returns entities with specific component
- `GetEntitiesWithComponents(compTypes ...ComponentType)` - Returns entities with all specified components
- `EntitiesAt(p gruid.Point)` - Returns entities at a specific position
- Special composite queries like `GetEntitiesWithPositionAndRenderable()`

## Components

### Position Component
- `Position struct { Point gruid.Point }`
- Used for entity location on the map

### Renderable Component
- `Renderable struct { Glyph rune, Color gruid.Color }`
- Visual representation for entities

### Health Component
- `Health struct { CurrentHP int, MaxHP int }`
- Tracks entity health status

### Name Component
- `Name struct { Name string }`
- Entity identification

### FOV Component
- `FOV struct { Range int, Visible []uint64, fov *rl.FOV }`
- Handles field of view calculations
- Includes bitset for efficient visibility tracking

### Tag Components
- `PlayerTag struct{}` - Marks player entity
- `AITag struct{}` - Marks AI-controlled entities
- `BlocksMovement struct{}` - Marks entities that block movement
- `DeadTag struct{}`, `CorpseTag struct{}` - Status markers

### Turn Component
- `TurnActor struct { Speed uint64, Alive bool, NextTurnTime uint64, actions *list.List }`
- Manages entity turn information and action queue

## Turn Management

### Turn Queue
- Located in `internal/turn_queue/`
- Uses min-heap for efficient ordering: `container/heap`
- Core structure: `TurnQueue` with internal `turnHeap []TurnEntry`
- `TurnEntry struct { Time uint64, EntityID ecs.EntityID }`

### Turn Queue Operations
- `Add(entityID EntityID, time uint64)` - Adds entity with scheduled time
- `Next()` - Removes and returns next entity (smallest time)
- `Peek()` - Returns next entity without removing
- `Remove(entityID EntityID)` - Removes entity from queue
- `CleanupDeadEntities(world *ecs.ECS)` - Removes invalid entities

## Game Actions

### Action Interface
- `GameAction interface { Execute(g *Game) (cost uint, err error) }`

### Action Types
- `MoveAction struct { Direction gruid.Point, EntityID ecs.EntityID }`
- `AttackAction struct { AttackerID ecs.EntityID, TargetID ecs.EntityID }`
- `WaitAction struct { EntityID ecs.EntityID }`

## Map System

### Map Structure
- `Map struct { Grid rl.Grid, Width int, Height int, Explored []uint64 }`
- Uses `rl.Grid` for storing map cells
- `Explored []uint64` bitset for tracking explored tiles

### Map Generation
- Room-based generation in `generateMap()`
- Room representation: `Rect struct { X1, Y1, X2, Y2 int }`
- `createRoom()`, `createHTunnel()`, `createVTunnel()` for map construction
- Monster placement via `placeMonsters()`

### Map Queries
- `InBounds(p gruid.Point)` - Checks if coordinates are within map
- `IsWall(p gruid.Point)` - Checks if tile is a wall
- `isWalkable(p gruid.Point)` - Checks if tile can be walked on
- `IsOpaque(p gruid.Point)` - Checks if tile blocks FOV (for FOV calculations)
- `IsExplored(p gruid.Point)` - Checks if tile has been explored

## Game Systems

### FOV System
- `FOVSystem(g *Game)` - Updates visibility for entities with FOV component
- Uses `rl.FOV` for calculations
- Manages global exploration tracking

### Render System
- `RenderSystem(ecs *ecs.ECS, grid gruid.Grid, playerFOV *components.FOV, mapWidth int)`
- Draws entities with Position and Renderable components
- Respects player's field of view

### Turn Processing
- `processTurnQueue()` - Processes turns until player's turn
- Monster AI in `monstersTurn()`

## UI and Input

### Input Handling
- Keyboard input mapped to actions in `internal/game/input.go`
- Mouse input supported via `updateMouse()`

### Drawing
- Main draw method in `model_draw.go`
- Uses `gruid.Grid` for terminal/display representation
- Custom styling via `internal/ui/color.go`

### Messaging
- Message log: `MessageLog struct { Messages []Message }`
- `AddMessage(text string, color gruid.Color)`
- Used for game events communication

## Miscellaneous Utilities

### FOV Utilities
- `DrawFilledCircle()` - Calculates cells within radius
- `VisionRange()` - Gets range for FOV calculations

### Assertion Utilities
- `Assert(condition bool, message string)`
- `Assertf(condition bool, format string, args ...any)`

## Game Initialization

### Game Structure
- `Game struct { Depth int, Map *Map, ecs *ecs.ECS, PlayerID ecs.EntityID, turnQueue *turn.TurnQueue, Log *log.MessageLog, rand *rand.Rand }`

### Initialization
- `NewGame()` - Creates new game instance
- `InitLevel()` - Sets up new game level
- `SpawnPlayer()`, `SpawnMonster()` - Create game entities
