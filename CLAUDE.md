# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands
- Build: `just build`
- Run: `just run` or `just run-dev` (development mode)
- Run with SDL: `just run-sdl`
- Test all: `just test` or verbose: `just test-verbose`
- Test specific package: `just test-package internal/ecs`
- Format code: `just fmt`
- Lint: `just lint`
- Clean: `just clean`

## Code Conventions
- Use Go modules (`go.mod`) for dependencies
- Follow Go standard naming: CamelCase for exported, camelCase for internal
- Group imports: standard lib, external, internal
- Error handling: use explicit error checking with meaningful error messages
- Prefer composition over inheritance with the ECS pattern
- Use meaningful comments for functions but avoid redundant comments
- Prefer clear names over abbreviations
- Document all exported functions, types, and constants
- Use constants for magic values and enums
- Format with `go fmt` before committing

## Git Workflow
- Use conventional commits: `feat(auth):`, `fix(map):`, `refactor(ecs):`, etc.
- Test before committing changes