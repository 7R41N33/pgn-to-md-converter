## 1. Dockerfile

- [x] 1.1 Create `Dockerfile` with multi-stage build: `golang:1.25.0-bookworm` build stage, `gcr.io/distroless/static-debian12` runtime stage
- [x] 1.2 Build both binaries (`build_markdown`, `fen_builder`) in the build stage with `CGO_ENABLED=0`
- [x] 1.3 Copy static sprite sheet (`src/images/Chess_Pieces_Sprite.svg.png`) into runtime stage
- [x] 1.4 Set default entrypoint to `/app/bin/build_markdown` so `docker run` prints usage info

## 2. Docker Compose

- [x] 2.1 Create `docker-compose.yml` with `build` and `dev` services using the same image
- [x] 2.2 Configure `build` service: runs `go build -o bin/build_markdown src/builders/markdown_builder/build_markdown.go && go build -o bin/fen_builder src/builders/fen_builder/fen_to_diagram.go`
- [x] 2.3 Configure `dev` service: uses Dockerfile for build (build context from root), no entrypoint or default command (`sleep infinity`), bind mounts for all directories (`PGNs/`, `Markdowns/`, `tmp/`, `bin/`), working dir `/app`, so `docker compose run dev bash` opens an interactive shell with full Go toolchain

## 3. Supporting Files

- [x] 3.1 Create `.dockerignore` excluding `.git`, `node_modules`, `tmp/`, `Markdowns/`, `bin/*.go` (but not `bin/` itself for bind-mount)
- [x] 3.2 Create `Makefile` with targets: `build` and `docker-build`

## 4. Documentation

- [x] 4.1 Update `AGENTS.md` to reference Docker Compose workflow for building and development
- [x] 4.2 Document docker-compose usage in `README.md` (or reference AGENTS.md)

## 5. Verification

- [x] 5.1 Build the image: `docker compose build`
- [x] 5.2 Verify binaries exist: `docker compose run build && ls -la bin/`
- [x] 5.3 Run `docker compose config` to verify compose file validity
- [x] 5.4 Verify dev service: `docker compose run dev bash` and run `go test ./src/...` inside the container
