# RFC 002: High-Performance FOV System Design

**Status:** Proposed

## 1. Overview

This document outlines the design for a high-performance Field of View (FOV) calculation and state management system for the Go Roguelike project. The primary goals are raw calculation speed and efficient state updates, anticipating maps up to 512x512 and frequent recalculations.

## 2. Algorithm

*   **Recommendation:** Utilize the **Symmetric Shadow Casting** algorithm, specifically the implementation provided by the `gruid/fov` package (`codeberg.org/anaseto/gruid/fov`).
*   **Justification:**
    *   **Performance:** Offers a good balance of performance (roughly O(R²)), visual quality, and implementation simplicity for grid-based FOV.
    *   **Integration:** Designed to work seamlessly with `gruid` types and the existing `Map.IsOpaque` method.
    *   **Suitability:** Handles opaque obstacles correctly and produces visually appealing, symmetrical results common in roguelikes.

## 3. Data Structures

*   **Map Grid (`Map.Grid`):**
    *   **Recommendation:** Continue using the existing `rl.Grid` storing `rl.Cell` values.
    *   **Justification:** Memory-efficient and provides necessary data for `IsOpaque` check.
*   **FOV State (`Visible`, `Explored`):**
    *   **Recommendation:** Replace the current `map[gruid.Point]bool` with **Bitsets**, implemented as 1D slices of `uint64` within the `Map` struct:
        *   `Visible []uint64`
        *   `Explored []uint64`
    *   **Calculation:** Size `(width * height + 63) / 64`. Tile `(x, y)` maps to bit index `idx = y*width + x`. Slice index `idx / 64`, bit `1 << (idx % 64)`.
    *   **Justification:** Maximizes memory efficiency (~32KB vs. 256KB+ for 512x512) and cache locality for faster access compared to maps or boolean slices. Bitwise operations are fast.

## 4. State Storage Strategy

*   **Recommendation:** Use a **Centralized** storage strategy, storing the `Visible` and `Explored` bitsets directly within the `Map` struct.
*   **Justification:**
    *   **Cache Efficiency:** Provides contiguous memory layout optimal for FOV calculation and rendering loops.
    *   **Decentralized Inefficiency:** Embedding flags in tile structs (if map structure allowed) would scatter data, harming cache performance.
    *   **Update Complexity:** Central bitset updates are efficient for bulk FOV results due to data locality.

## 5. Go Implementation Patterns

*   **Minimize Allocations / GC Pressure:**
    *   Implement the `fov.Grid` interface on the `Map` struct. Ensure the `SetVisible` method updates the bitsets directly and efficiently without intermediate allocations.
    *   Create efficient, potentially inlined bitset helper functions (`isVisible`, `isExplored`, `setVisible`, `setExplored`).
    *   Use value types (`gruid.Point`) where possible in calculation logic.
*   **Cache-Friendly Access:**
    *   Leverage the inherent cache locality of the `[]uint64` bitsets.
    *   Ensure rendering loops access visibility/explored status sequentially via helper functions.
*   **Concurrency (Goroutines):**
    *   **Recommendation:** Avoid concurrency initially. Focus on single-threaded optimization.
    *   **Justification:** Parallelizing shadow casting is complex, and synchronization overhead might negate benefits for typical FOV radii. Revisit only if profiling proves FOV is a major bottleneck after optimization.
*   **Code Structure & Testability:**
    *   Encapsulate FOV logic in `map.go` or a new `fov.go`.
    *   `Map` type implements `fov.Grid`.
    *   Provide clear `ComputeFOV` function in `Map`.
    *   Add unit tests for bitset logic and `IsOpaque`.

## 6. System Interaction Diagram

```mermaid
graph TD
    subgraph MapStruct [Map Struct (map.go)]
        Grid["Grid (rl.Grid)"]
        Width["Width (int)"]
        Height["Height (int)"]
        VisibleBitset["Visible ([]uint64)"]
        ExploredBitset["Explored ([]uint64)"]
        IsOpaque["IsOpaque(p) bool"]
        SetVisible["SetVisible(p)"]
        ComputeFOV["ComputeFOV(origin, radius)"]
        IsVisibleHelper["isVisible(p) bool"]
        IsExploredHelper["isExplored(p) bool"]
    end

    subgraph FOVCalculation [FOV Calculation (fov.go / map.go)]
        GruidFOV["gruid/fov.Compute()"] -- Calls --> IsOpaque
        GruidFOV -- Calls --> SetVisible
        ComputeFOV -- Calls --> GruidFOV
        SetVisible -- Updates --> VisibleBitset
        SetVisible -- Updates --> ExploredBitset
    end

    subgraph Rendering [Rendering (model_draw.go / systems.go)]
        RenderSystem["RenderSystem / DrawMap"] -- Iterates Map --> CheckVisibility
        CheckVisibility -- Calls --> IsVisibleHelper
        CheckVisibility -- Calls --> IsExploredHelper
        IsVisibleHelper -- Reads --> VisibleBitset
        IsExploredHelper -- Reads --> ExploredBitset
    end

    MapStruct --> FOVCalculation
    MapStruct --> Rendering
```

## 7. Assumptions & Trade-offs

*   **Assumption:** `gruid/fov` implementation is sufficiently performant and allows efficient integration via the `fov.Grid` interface.
*   **Assumption:** Bitwise operation complexity is negligible compared to cache gains.
*   **Trade-off:** Bitset logic adds minor complexity for significant performance/memory benefits.
*   **Trade-off:** Sticking with `gruid/fov` saves development effort vs. a custom algorithm.
*   **Trade-off:** Avoiding concurrency simplifies implementation; revisit if profiling demands it later.
