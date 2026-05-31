## ADDED Requirements

### Requirement: Dockerfile at project root
The project SHALL include a multi-stage Dockerfile at the project root that builds both binaries and serves as the runtime environment.

#### Scenario: Build both binaries
- **WHEN** `docker compose build` is run
- **THEN** the image SHALL contain `/app/bin/build_markdown` and `/app/bin/fen_builder` as executable binaries

#### Scenario: Runtime stage includes sprite asset
- **WHEN** the image is built
- **THEN** the runtime stage SHALL include `src/images/Chess_Pieces_Sprite.svg.png` at the same relative path

#### Scenario: Minimal runtime image
- **WHEN** the image is built
- **THEN** the runtime stage SHALL use a distroless or otherwise minimal base image

#### Scenario: Default command prints usage
- **WHEN** `docker run chess-pgns` is executed without arguments
- **THEN** the container SHALL print usage information and exit with a non-zero code

### Requirement: Build stage SHALL use Go 1.25.0
The build stage SHALL use the official `golang:1.25.0-bookworm` image to match the project's Go version.

#### Scenario: Go version matches go.mod
- **WHEN** the build stage runs `go version`
- **THEN** it SHALL report Go 1.25.0

### Requirement: .dockerignore excludes unnecessary files
The `.dockerignore` file at the project root SHALL exclude files that are not needed for the build (`.git`, `node_modules`, `tmp/`, `Markdowns/`, `bin/`).

#### Scenario: .dockerignore is present
- **WHEN** the project root is inspected
- **THEN** a `.dockerignore` file SHALL exist

### Requirement: Build with CGO_ENABLED=0
The Dockerfile SHALL build with `CGO_ENABLED=0` to produce fully static Go binaries.

#### Scenario: Binary is statically linked
- **WHEN** `file bin/build_markdown` is run on the built binary
- **THEN** it SHALL show "statically linked"
