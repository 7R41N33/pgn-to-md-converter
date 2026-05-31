.PHONY: build docker-build

build:
	go build -o bin/build_markdown src/builders/markdown_builder/build_markdown.go
	go build -o bin/fen_builder src/builders/fen_builder/fen_to_diagram.go

docker-build:
	docker compose run build
