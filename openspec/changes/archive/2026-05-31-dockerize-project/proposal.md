## Why

The project has no containerized build or runtime environment. Developers must manually install Go 1.25.0, clone the repo, and run `go build` commands. There is no reproducible, isolated environment for building, testing, or running the PGN→Markdown toolchain. AGENTS.md already prescribes a Dockerfile but none exists yet.

## What Changes

- Create a multi-stage `Dockerfile` at the project root that builds both binaries and serves as the reference runtime environment
- Create a `docker-compose.yml` at the project root for local development with volume mounts for PGN input and Markdown output
- Add a `Makefile` with targets (`build`, `docker-build`)
- Add `.dockerignore` to exclude unnecessary files from the Docker build context

## Capabilities

### New Capabilities
- `docker-image`: Multi-stage Dockerfile for building and running the chess-pgns toolchain. Includes build stage (Go 1.25.0) and runtime stage (distroless) with the sprite sheet asset.
- `docker-compose-dev`: Docker Compose services for local development — build and dev — with bind mounts for PGNs/, Markdowns/ and tmp/.

### Modified Capabilities
- *(none — no existing specs to modify)*

## Impact

- `Dockerfile` — new file, the reference container definition
- `docker-compose.yml` — new file at project root
- `Makefile` — new file at project root with `build` and `docker-build` targets
- `.dockerignore` — new file at project root
- `AGENTS.md` — update to reference Docker Compose commands for dev workflow
- All existing Go code, tests, and data directories remain unchanged
