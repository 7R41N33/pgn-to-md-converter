package fenlib

import (
	"os"
	"testing"
)

func TestValidateFEN_Valid(t *testing.T) {
	validFENs := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		"8/8/8/8/8/8/8/8 w - - 0 1",
	}
	for _, fen := range validFENs {
		if err := ValidateFEN(fen); err != nil {
			t.Errorf("Valid FEN rejected: %s - %v", fen, err)
		}
	}
}

func TestValidateFEN_InvalidPosition(t *testing.T) {
	invalidFENs := []string{
		"rnbqkbnr/pppppppp/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0", // Missing fullmove
		"rnbqkbnr/pppppppp/9/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", // Invalid row
	}
	for _, fen := range invalidFENs {
		if err := ValidateFEN(fen); err == nil {
			t.Errorf("Invalid FEN accepted: %s", fen)
		}
	}
}

func TestValidateFEN_InvalidActiveColor(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR x KQkq - 0 1"
	if err := ValidateFEN(fen); err == nil {
		t.Error("Should reject invalid active color")
	}
}

func TestValidateFEN_InvalidCastling(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w XYZ - 0 1"
	if err := ValidateFEN(fen); err == nil {
		t.Error("Should reject invalid castling")
	}
}

func TestValidateFEN_InvalidEnPassant(t *testing.T) {
	invalidEP := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq i9 0 1",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq e9 0 1",
	}
	for _, fen := range invalidEP {
		if err := ValidateFEN(fen); err == nil {
			t.Errorf("Should reject invalid en passant: %s", fen)
		}
	}
}

func TestGenerateDiagram_CreatesFile(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	outputPath := "/tmp/test_diagram.png"

	if err := GenerateDiagram(fen, outputPath); err != nil {
		t.Errorf("Failed to generate diagram: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
	os.Remove(outputPath)
}

func TestGenerateDiagram_InvalidFEN(t *testing.T) {
	fen := "invalid fen"

	if err := GenerateDiagram(fen, "/tmp/test.png"); err == nil {
		t.Error("Should return error for invalid FEN")
	}
}

func TestValidateFEN_EmptyString(t *testing.T) {
	if err := ValidateFEN(""); err == nil {
		t.Error("Empty string should be invalid")
	}
}

func TestValidateFEN_TooFewFields(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq"
	if err := ValidateFEN(fen); err == nil {
		t.Error("Should reject FEN with too few fields")
	}
}

func TestGenerateDiagram_OutputPath(t *testing.T) {
	fen := "8/8/8/8/8/8/8/8 w - - 0 1"
	outputPath := "/tmp/test_empty_board.png"

	if err := GenerateDiagram(fen, outputPath); err != nil {
		t.Errorf("Failed to generate diagram: %v", err)
	}

	content, _ := os.ReadFile(outputPath)
	if len(content) == 0 {
		t.Error("Output file should not be empty")
	}
	os.Remove(outputPath)
}

func TestGenerateDiagramBase64(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	b64, err := GenerateDiagramBase64(fen)
	if err != nil {
		t.Errorf("Failed to generate base64: %v", err)
	}
	if b64 == "" {
		t.Error("Base64 output should not be empty")
	}
}
