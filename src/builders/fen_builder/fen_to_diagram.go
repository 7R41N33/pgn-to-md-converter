package main

import (
	"fmt"
	"os"
	"strings"

	"fen-diagram/src/builders/fenlib"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <fen_string> [output_path] [--base64]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  --base64: output base64-encoded PNG to stdout instead of file\n")
		os.Exit(1)
	}

	fen := os.Args[1]
	outputPath := "tmp/images/diagram.png"
	base64Output := false

	// Parse arguments
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--base64" {
			base64Output = true
		} else if !strings.HasPrefix(os.Args[i], "--") {
			outputPath = os.Args[i]
		}
	}

	if err := fenlib.ValidateFEN(fen); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid FEN - %v\n", err)
		os.Exit(1)
	}

	if base64Output {
		b64, err := fenlib.GenerateDiagramBase64(fen)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating diagram: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(b64)
	} else {
		if err := fenlib.GenerateDiagram(fen, outputPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating diagram: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Diagram saved to: %s\n", outputPath)
	}
}
