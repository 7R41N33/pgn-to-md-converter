## ADDED Requirements

### Requirement: docker-compose.yml at project root
The project SHALL include a `docker-compose.yml` at the project root with services for local development.

#### Scenario: docker-compose.yml is parseable
- **WHEN** `docker compose config` is run
- **THEN** it SHALL exit successfully without errors

### Requirement: build service compiles both binaries
The `build` service SHALL run `go build` for both `build_markdown` and `fen_builder`, outputting to `/app/bin/`.

#### Scenario: Build service produces binaries
- **WHEN** `docker compose run build` is executed
- **THEN** the `bin/` directory on the host SHALL contain updated `build_markdown` and `fen_builder` binaries

### Requirement: dev service for interactive development
The `dev` service SHALL be a long-running container with no default command (`sleep infinity` or `tail -f /dev/null`), with bind mounts for `PGNs/`, `Markdowns/`, `tmp/` and `bin/`, so the developer can enter a shell and run any Go commands.

#### Scenario: Connect to dev service
- **WHEN** `docker compose run dev bash` is executed
- **THEN** an interactive shell opens inside the container with the Go toolchain

#### Scenario: Run go build through dev
- **WHEN** inside the dev container `go build -o bin/build_markdown src/builders/markdown_builder/build_markdown.go` is executed
- **THEN** the binary appears in `bin/` on the host

### Requirement: Makefile with convenience targets
A `Makefile` at the project root SHALL provide targets: `build` (native build) and `docker-build` (build via docker compose).

#### Scenario: make docker-build invokes docker compose
- **WHEN** `make docker-build` is run
- **THEN** it SHALL execute `docker compose run build`
