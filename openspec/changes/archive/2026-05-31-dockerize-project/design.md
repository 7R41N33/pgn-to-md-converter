## Context

The chess-pgns project is a Go 1.25.0 monorepo with two CLI binaries (`build_markdown`, `fen_builder`) sharing one library (`fenlib`). There are no external services — all operations are local filesystem I/O (read PGN from `PGNs/`, write Markdown to `Markdowns/`, render PNG to `tmp/`). The sprite sheet asset is bundled at `src/images/Chess_Pieces_Sprite.svg.png`. Tests invoke compiled binaries via relative paths (`../../../bin/build_markdown`).

AGENTS.md prescribes a Dockerfile as the single source of truth for environment setup. Currently none exists. `go.mod` and `go.sum` live at the project root, so the Dockerfile must account for the actual layout.

## Goals / Non-Goals

**Goals:**
- Provide a reproducible multi-stage Docker image that builds both binaries and can run them
- Provide Docker Compose services for common dev workflows: build and dev
- Bind-mount `PGNs/`, `Markdowns/`, `tmp/` and `bin/` so developers edit files on the host and see output immediately
- Add a `Makefile` as a DX layer over `docker compose` commands
- Keep the existing `go build -o bin/...` workflow working for users who prefer native Go

**Non-Goals:**
- CI/CD integration (separate concern, uses existing GitHub Actions workflow)
- Publishing images to a registry (local-only for now)
- Changing the project's Go module layout or moving `go.mod`
- Containerizing the Google Apps Script (`src/base64_to_image.gs`)
- Adding hot-reload or file-watching

## Decisions

1. **Multi-stage Dockerfile: build → runtime**
   - **Why**: Minimizes final image size. Build stage has full Go toolchain; runtime stage has only the binaries and sprite sheet.
   - **Alternatives considered**: Single-stage (larger image, but simpler) — rejected because runtime doesn't need Go.

2. **Runtime base image: `gcr.io/distroless/static-debian12`**
   - **Why**: Distroless images contain only the application and its runtime dependencies (libc, ca-certificates). No shell, no package manager — minimal attack surface and <10 MB overhead.
   - **Alternatives considered**: `alpine:3.21` (musl libc, slightly different behavior if CGO is ever needed) — viable, but distroless is more secure for a static Go binary. `debian:bookworm-slim` (larger, has package manager). Since the Go binary is fully static (`CGO_ENABLED=0`), distroless/static is ideal.

3. **Docker Compose services: build, dev**
   - **build**: Runs `go build` inside the container for both binaries, outputs to `bin/` (bind-mounted).
   - **dev**: Long-running container with no default command (`sleep infinity`), all directories (`PGNs/`, `Markdowns/`, `tmp/`, `bin/`) bind-mounted so the developer can enter a shell and run any Go commands.

4. **`.dockerignore` to exclude development artifacts**
   - **Why**: Prevents sending unnecessary context (`.git/`, `node_modules/`, `tmp/`, `Markdowns/`) to the Docker daemon, speeding up builds. Also prevents cache invalidation from irrelevant changes.

## Risks / Trade-offs

- **Bind mount performance**: On macOS (Docker Desktop) and Windows (WSL2), bind mounts can be slower than native filesystem. Mitigation: The project's I/O is lightweight (text files, small PNGs), so this is unlikely to be noticeable.
- **Test path assumptions**: Tests reference `../../../bin/build_markdown`. Inside the container, the working directory is `/app` (project root), so relative paths from test files resolve the same way. This should work unchanged.
- **Sprite path**: The Go code reads `src/images/Chess_Pieces_Sprite.svg.png` as a relative path from the working directory. In the container, the sprite must be present at the same relative location. The Dockerfile COPY preserves this structure.
