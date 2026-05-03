package main

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
		if err := validateFEN(fen); err != nil {
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
		if err := validateFEN(fen); err == nil {
			t.Errorf("Invalid FEN accepted: %s", fen)
		}
	}
}

func TestValidateFEN_InvalidActiveColor(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR x KQkq - 0 1"
	if err := validateFEN(fen); err == nil {
		t.Error("Should reject invalid active color")
	}
}

func TestValidateFEN_InvalidCastling(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w XYZ - 0 1"
	if err := validateFEN(fen); err == nil {
		t.Error("Should reject invalid castling")
	}
}

func TestValidateFEN_InvalidEnPassant(t *testing.T) {
	invalidEP := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq i9 0 1",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq e9 0 1",
	}
	for _, fen := range invalidEP {
		if err := validateFEN(fen); err == nil {
			t.Errorf("Should reject invalid en passant: %s", fen)
		}
	}
}

func TestIsValidCastling_Valid(t *testing.T) {
	valid := []string{"-", "KQkq", "K", "Q", "k", "q", "KQ", "kq"}
	for _, s := range valid {
		if !isValidCastling(s) {
			t.Errorf("Valid castling rejected: %s", s)
		}
	}
}

func TestIsValidCastling_Invalid(t *testing.T) {
	invalid := []string{"KQkqX", "abc", "123"}
	for _, s := range invalid {
		if isValidCastling(s) {
			t.Errorf("Invalid castling accepted: %s", s)
		}
	}
}

func TestIsValidEnPassant_Valid(t *testing.T) {
	valid := []string{"-", "e3", "a6", "h3", "d6"}
	for _, s := range valid {
		if !isValidEnPassant(s) {
			t.Errorf("Valid en passant rejected: %s", s)
		}
	}
}

func TestIsValidEnPassant_Invalid(t *testing.T) {
	invalid := []string{"e9", "i3", "a7", ""}
	for _, s := range invalid {
		if isValidEnPassant(s) {
			t.Errorf("Invalid en passant accepted: %s", s)
		}
	}
}

func TestGenerateDiagram_CreatesFile(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	outputPath := "/tmp/test_diagram.png"

	if err := generateDiagram(fen, outputPath); err != nil {
		t.Errorf("Failed to generate diagram: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
	os.Remove(outputPath)
}

func TestGenerateDiagram_InvalidFEN(t *testing.T) {
	fen := "invalid fen"

	// Invalid FEN should fail validation
	if err := validateFEN(fen); err == nil {
		t.Error("Should return error for invalid FEN")
	}
}

func TestValidateFEN_EmptyString(t *testing.T) {
	if err := validateFEN(""); err == nil {
		t.Error("Empty string should be invalid")
	}
}

func TestValidateFEN_TooFewFields(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq"
	if err := validateFEN(fen); err == nil {
		t.Error("Should reject FEN with too few fields")
	}
}

func TestGenerateDiagram_OutputPath(t *testing.T) {
	fen := "8/8/8/8/8/8/8 w - - 0 1"
	outputPath := "/tmp/test_empty_board.png"

	if err := generateDiagram(fen, outputPath); err != nil {
		t.Errorf("Failed to generate diagram: %v", err)
	}

	content, _ := os.ReadFile(outputPath)
	if len(content) == 0 {
		t.Error("Output file should not be empty")
	}
	os.Remove(outputPath)
}

func TestDefaultOutputPath(t *testing.T) {
	// Test that default path is used when not specified
	// This is tested indirectly via command line args
	if defaultOut != "tmp/images/diagram.png" {
		t.Errorf("Default output path should be tmp/images/diagram.png, got %s", defaultOut)
	}
}

func TestGenerateDiagram_ProvidedFEN(t *testing.T) {
	// Test with the FEN from user: 8/pp3ppp/8/2p5/8/3P2P1/PP2P2P/8 w - - 0 1
	fen := "8/pp3ppp/8/2p5/8/3P2P1/PP2P2P/8 w - - 0 1"
	outputPath := "/tmp/test_user_fen.png"

	if err := generateDiagram(fen, outputPath); err != nil {
		t.Errorf("Failed to generate diagram: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
	os.Remove(outputPath)
}
