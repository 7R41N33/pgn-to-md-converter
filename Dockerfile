FROM golang:1.25.0-bookworm AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin/build_markdown ./src/builders/markdown_builder/build_markdown.go && \
    CGO_ENABLED=0 go build -o /app/bin/fen_builder ./src/builders/fen_builder/fen_to_diagram.go

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=build /app/bin/ /app/bin/
COPY src/images/ src/images/
ENTRYPOINT ["/app/bin/build_markdown"]
