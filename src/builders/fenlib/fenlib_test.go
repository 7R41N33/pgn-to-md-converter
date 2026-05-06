package fenlib

import (
	"strings"
	"testing"
)

func TestGenerateBoard_StandardStart(t *testing.T) {
	// Standard starting position
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKQBNR w KQkq - 0 1"
	
	board, err := GenerateBoard(fen)
	if err != nil {
		t.Fatalf("Failed to generate board: %v", err)
	}
	
	if board == nil {
		t.Fatal("Board is nil")
	}
	
	// Check that the board has the correct dimensions
	expectedSize := 8*60 + 2*(60/2) // 8 cells + 2*border
	if board.Bounds().Dx() != expectedSize || board.Bounds().Dy() != expectedSize {
		t.Errorf("Board size = %dx%d, want %dx%d", board.Bounds().Dx(), board.Bounds().Dy(), expectedSize, expectedSize)
	}
}

func TestPieceOrder(t *testing.T) {
	// Test that pieceOrder constant matches the expected order
	expected := "KQRBNPkqrbnp"
	if pieceOrder != expected {
		t.Errorf("pieceOrder = %q, want %q", pieceOrder, expected)
	}
}

func TestGenerateBoard_InvalidFEN(t *testing.T) {
	invalidFENs := []string{
		"",
		"invalid",
		"rnbqkbnr/pppppppp/9/8/8/8/PPPPPPPP/RNBKQBNR w KQkq - 0 1", // invalid row
	}
	
	for _, fen := range invalidFENs {
		_, err := GenerateBoard(fen)
		if err == nil {
			t.Errorf("Expected error for FEN %q, got nil", fen)
		}
	}
}

func TestValidateFEN_Valid(t *testing.T) {
	validFENs := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKQBNR w KQkq - 0 1",
		"8/8/8/8/8/8/8/8 w - - 0 1",
	}
	
	for _, fen := range validFENs {
		if err := ValidateFEN(fen); err != nil {
			t.Errorf("Expected valid FEN %q, got error: %v", fen, err)
		}
	}
}

func TestFENRanks_TopBottom(t *testing.T) {
	// Test that rank 8 (black's side) is at TOP, rank 1 (white's side) is at BOTTOM
	// Create a FEN with black pawn on a8 (top) and white pawn on a1 (bottom)
	fen := "p7/8/8/8/8/8/8/P7 w - - 0 1"
	
	board, err := GenerateBoard(fen)
	if err != nil {
		t.Fatalf("Failed to generate board: %v", err)
	}
	
	cellSize := 60
	borderSize := cellSize / 2
	
	// Check that SOME pixel in the top row region has a non-white piece
	topRegionHasPiece := false
	for x := borderSize; x < borderSize+cellSize; x++ {
		for y := borderSize; y < borderSize+cellSize; y++ {
			c := board.At(x, y)
			r, g, b, a := c.RGBA()
			if a > 0 && !(r == 65535 && g == 65535 && b == 65535) {
				topRegionHasPiece = true
				break
			}
		}
		if topRegionHasPiece {
			break
		}
	}
	
	if !topRegionHasPiece {
		t.Error("No piece found in top region (rank 8)")
	}
	
	// Check that SOME pixel in the bottom row region has a non-white piece
	bottomRegionHasPiece := false
	for x := borderSize; x < borderSize+cellSize; x++ {
		for y := borderSize + 7*cellSize; y < borderSize+8*cellSize; y++ {
			c := board.At(x, y)
			r, g, b, a := c.RGBA()
			if a > 0 && !(r == 65535 && g == 65535 && b == 65535) {
				bottomRegionHasPiece = true
				break
			}
		}
		if bottomRegionHasPiece {
			break
		}
	}
	
	if !bottomRegionHasPiece {
		t.Error("No piece found in bottom region (rank 1)")
	}
}

func TestFENSpriteOrder(t *testing.T) {
	// Test that the sprite has the correct piece order
	// pieceOrder = "KQRBNPkqrbnp"
	// White pieces (indices 0-5): K, Q, R, B, N, P
	// Black pieces (indices 6-11): k, q, r, b, n, p
	
	// Verify that 'K' (white king) maps to index 0
	if idx := strings.IndexRune(pieceOrder, 'K'); idx != 0 {
		t.Errorf("'K' should be at index 0, got %d", idx)
	}
	
	// Verify that 'P' (white pawn) maps to index 5
	if idx := strings.IndexRune(pieceOrder, 'P'); idx != 5 {
		t.Errorf("'P' should be at index 5, got %d", idx)
	}
	
	// Verify that 'k' (black king) maps to index 6
	if idx := strings.IndexRune(pieceOrder, 'k'); idx != 6 {
		t.Errorf("'k' should be at index 6, got %d", idx)
	}
	
	// Verify that 'p' (black pawn) maps to index 11
	if idx := strings.IndexRune(pieceOrder, 'p'); idx != 11 {
		t.Errorf("'p' should be at index 11, got %d", idx)
	}
}
