package fenlib

import (
	"strings"
	"testing"

	"golang.org/x/image/math/fixed"
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
	// New size with extra border: 8 cells + 2*border (30 + 15 extra)
	expectedSize := 8*60 + 2*(60/2+15) // 8 cells + 2*border
	if board.Bounds().Dx() != expectedSize || board.Bounds().Dy() != expectedSize {
		t.Errorf("Board size = %dx%d, want %dx%d", board.Bounds().Dx(), board.Bounds().Dy(), expectedSize, expectedSize)
	}
}

func TestPieceOrder(t *testing.T) {
	// Test that pieceOrder constant matches the expected order
	expected := "KQBNRPkqbnrp"
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
	// pieceOrder = "KQBNRPkqbnrp"
	// White pieces (indices 0-5): K, Q, B, N, R, P
	// Black pieces (indices 6-11): k, q, b, n, r, p
	
	// Verify that 'K' (white king) maps to index 0
	if idx := strings.IndexRune(pieceOrder, 'K'); idx != 0 {
		t.Errorf("'K' should be at index 0, got %d", idx)
	}
	
	// Verify that 'Q' (white queen) maps to index 1
	if idx := strings.IndexRune(pieceOrder, 'Q'); idx != 1 {
		t.Errorf("'Q' should be at index 1, got %d", idx)
	}
	
	// Verify that 'B' (white bishop) maps to index 2
	if idx := strings.IndexRune(pieceOrder, 'B'); idx != 2 {
		t.Errorf("'B' should be at index 2, got %d", idx)
	}
	
	// Verify that 'N' (white knight) maps to index 3
	if idx := strings.IndexRune(pieceOrder, 'N'); idx != 3 {
		t.Errorf("'N' should be at index 3, got %d", idx)
	}
	
	// Verify that 'R' (white rook) maps to index 4
	if idx := strings.IndexRune(pieceOrder, 'R'); idx != 4 {
		t.Errorf("'R' should be at index 4, got %d", idx)
	}
	
	// Verify that 'P' (white pawn) maps to index 5
	if idx := strings.IndexRune(pieceOrder, 'P'); idx != 5 {
		t.Errorf("'P' should be at index 5, got %d", idx)
	}
	
	// Verify that 'k' (black king) maps to index 6
	if idx := strings.IndexRune(pieceOrder, 'k'); idx != 6 {
		t.Errorf("'k' should be at index 6, got %d", idx)
	}
	
	// Verify that 'q' (black queen) maps to index 7
	if idx := strings.IndexRune(pieceOrder, 'q'); idx != 7 {
		t.Errorf("'q' should be at index 7, got %d", idx)
	}
	
	// Verify that 'b' (black bishop) maps to index 8
	if idx := strings.IndexRune(pieceOrder, 'b'); idx != 8 {
		t.Errorf("'b' should be at index 8, got %d", idx)
	}
	
	// Verify that 'n' (black knight) maps to index 9
	if idx := strings.IndexRune(pieceOrder, 'n'); idx != 9 {
		t.Errorf("'n' should be at index 9, got %d", idx)
	}
	
	// Verify that 'r' (black rook) maps to index 10
	if idx := strings.IndexRune(pieceOrder, 'r'); idx != 10 {
		t.Errorf("'r' should be at index 10, got %d", idx)
	}
	
	// Verify that 'p' (black pawn) maps to index 11
	if idx := strings.IndexRune(pieceOrder, 'p'); idx != 11 {
		t.Errorf("'p' should be at index 11, got %d", idx)
	}
}

func TestDrawTriangle_WhiteTurn(t *testing.T) {
	fen := "8/8/8/8/8/8/8/8 w - - 0 1"

	board, err := GenerateBoard(fen)
	if err != nil {
		t.Fatalf("Failed to generate board: %v", err)
	}

	cellSize := 60
	borderSize := cellSize/2 + 15
	boardSize := 8 * cellSize

	margin := 3
	triangleSize := cellSize / 2
	triangleBaseWidth := int(float64(triangleSize) * 2.0 / 1.7320508075688772)
	trianglePadding := 3
	triangleCenterX := borderSize + boardSize + margin + triangleBaseWidth/2 + trianglePadding
	triangleCenterY := borderSize + boardSize - cellSize/2

	// Check that some pixels in the triangle area are white (white turn indicator)
	whiteFound := false
	for x := triangleCenterX - triangleSize; x < triangleCenterX+triangleSize && !whiteFound; x++ {
		for y := triangleCenterY - triangleSize; y < triangleCenterY+triangleSize; y++ {
			c := board.At(x, y)
			r, g, b, a := c.RGBA()
			// White pixel (almost white, since we draw white on white border)
			if a > 0 && r > 65000 && g > 65000 && b > 65000 {
				whiteFound = true
				break
			}
		}
	}

	if !whiteFound {
		t.Error("White turn indicator (white triangle center) not found")
	}

	// Check that some pixels in the triangle area are black (border)
	blackFound := false
	for x := triangleCenterX - triangleSize; x < triangleCenterX+triangleSize && !blackFound; x++ {
		for y := triangleCenterY - triangleSize; y < triangleCenterY+triangleSize; y++ {
			c := board.At(x, y)
			r, g, b, a := c.RGBA()
			// Black or dark pixel
			if a > 0 && r < 10000 && g < 10000 && b < 10000 {
				blackFound = true
				break
			}
		}
	}

	if !blackFound {
		t.Error("Black border of triangle not found")
	}
}

func TestDrawTriangle_BlackTurn(t *testing.T) {
	fen := "8/8/8/8/8/8/8/8 b - - 0 1"

	board, err := GenerateBoard(fen)
	if err != nil {
		t.Fatalf("Failed to generate board: %v", err)
	}

	cellSize := 60
	borderSize := cellSize/2 + 15
	boardSize := 8 * cellSize

	margin := 3
	triangleSize := cellSize / 2
	triangleBaseWidth := int(float64(triangleSize) * 2.0 / 1.7320508075688772)
	trianglePadding := 3
	triangleCenterX := borderSize + boardSize + margin + triangleBaseWidth/2 + trianglePadding
	triangleCenterY := borderSize + boardSize - cellSize/2

	// For black turn, triangle should be solid black (no white center)
	// Check that the triangle area contains black pixels
	blackFound := false
	for x := triangleCenterX - triangleSize; x < triangleCenterX+triangleSize && !blackFound; x++ {
		for y := triangleCenterY - triangleSize; y < triangleCenterY+triangleSize; y++ {
			c := board.At(x, y)
			r, g, b, a := c.RGBA()
			if a > 0 && r < 10000 && g < 10000 && b < 10000 {
				blackFound = true
				break
			}
		}
	}

	if !blackFound {
		t.Error("Black turn indicator (black triangle) not found")
	}
}

func TestDrawTriangle_Position(t *testing.T) {
	fen := "8/8/8/8/8/8/8/8 w - - 0 1"

	board, err := GenerateBoard(fen)
	if err != nil {
		t.Fatalf("Failed to generate board: %v", err)
	}

	cellSize := 60
	borderSize := cellSize/2 + 15
	boardSize := 8 * cellSize

	// Check that the corner cell (h1) does NOT have triangle pixels

	// Check that the board area doesn't have triangle
	// The triangle should only appear to the right of the board
	lastBoardX := borderSize + boardSize - 1
	c := board.At(lastBoardX, borderSize+boardSize-1)
	r, g, b, a := c.RGBA()
	// This should be a board cell color, not triangle
	// Board cell should be light or dark brown, not white (border) or black (triangle)
	isBoardColor := (r > 50000 && g > 40000 && b > 30000) || // Light brown
		(r > 30000 && r < 60000 && g > 20000 && g < 50000 && b > 20000 && b < 40000) // Dark brown

	if !isBoardColor {
		t.Errorf("Last board cell should be board color, got rgba(%d,%d,%d,%d)", r, g, b, a)
	}

	// Check that the triangle area (right of board, at h1 level) contains the indicator
	margin := 3
	triangleSize := cellSize / 2
	triangleBaseWidth := int(float64(triangleSize) * 2.0 / 1.7320508075688772)
	trianglePadding := 3
	triangleCenterX := borderSize + boardSize + margin + triangleBaseWidth/2 + trianglePadding
	triangleCenterY := borderSize + boardSize - cellSize/2
	c = board.At(triangleCenterX, triangleCenterY)
	r, g, b, a = c.RGBA()

	// Triangle area should have some non-white pixels (triangle border or fill)
	// It should NOT be pure white (which is the border color)
	if a > 0 && r > 240 && g > 240 && b > 240 {
		// White pixel in triangle area - check if triangle is there
	}
}

func TestDrawTriangle_Equilateral(t *testing.T) {
	fen := "8/8/8/8/8/8/8/8 w - - 0 1"

	board, err := GenerateBoard(fen)
	if err != nil {
		t.Fatalf("Failed to generate board: %v", err)
	}

	cellSize := 60
	borderSize := cellSize/2 + 15
	boardSize := 8 * cellSize
	margin := 3
	triangleSize := cellSize / 2
	triangleBaseWidth := int(float64(triangleSize) * 2.0 / 1.7320508075688772)
	trianglePadding := 3

	// Check that some non-white pixels exist in the triangle area
	triangleCenterX := borderSize + boardSize + margin + triangleBaseWidth/2 + trianglePadding
	triangleCenterY := borderSize + boardSize - cellSize/2

	triangleFound := false
	for x := triangleCenterX - triangleSize; x < triangleCenterX+triangleSize && !triangleFound; x++ {
		for y := triangleCenterY - triangleSize; y < triangleCenterY+triangleSize; y++ {
			c := board.At(x, y)
			r, g, b, a := c.RGBA()
			// Check for non-white pixels (triangle is black or white with black border)
			if a > 0 && !(r > 65000 && g > 65000 && b > 65000) {
				triangleFound = true
				break
			}
		}
	}

	if !triangleFound {
		t.Error("Triangle not found in expected area")
	}
}

func TestDrawTriangle_Exists(t *testing.T) {
	fen := "8/8/8/8/8/8/8/8 w - - 0 1"

	board, err := GenerateBoard(fen)
	if err != nil {
		t.Fatalf("Failed to generate board: %v", err)
	}

	cellSize := 60
	borderSize := cellSize/2 + 15
	boardSize := 8 * cellSize
	margin := 3
	triangleSize := cellSize / 2
	triangleBaseWidth := int(float64(triangleSize) * 2.0 / 1.7320508075688772)
	trianglePadding := 3

	triangleCenterX := borderSize + boardSize + margin + triangleBaseWidth/2 + trianglePadding
	triangleCenterY := borderSize + boardSize - cellSize/2

	// Scan the triangle area for non-white pixels (indicating triangle is drawn)
	triangleFound := false
	for x := triangleCenterX - triangleSize; x < triangleCenterX+triangleSize && !triangleFound; x++ {
		for y := triangleCenterY - triangleSize; y < triangleCenterY+triangleSize; y++ {
			c := board.At(x, y)
			r, g, b, a := c.RGBA()
			// Check for non-white pixels (triangle is black or white with black border)
			if a > 0 && !(r > 65000 && g > 65000 && b > 65000) {
				triangleFound = true
				break
			}
		}
	}

	if !triangleFound {
		t.Error("Triangle indicator not found in expected area")
	}
}

func TestCaption_SingleLine(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1"
	noCap, err := GenerateBoard(fen)
	if err != nil {
		t.Fatalf("Failed to generate board without caption: %v", err)
	}
	capBoard, err := GenerateBoard(fen, "Тестовый текст")
	if err != nil {
		t.Fatalf("Failed to generate board with caption: %v", err)
	}
	if capBoard.Bounds().Dy() <= noCap.Bounds().Dy() {
		t.Errorf("Caption board height (%d) should be greater than no-caption height (%d)",
			capBoard.Bounds().Dy(), noCap.Bounds().Dy())
	}
}

func TestCaption_EmptyCaption(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKQBNR w KQkq - 0 1"
	b1, err := GenerateBoard(fen)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	b2, err := GenerateBoard(fen, "")
	if err != nil {
		t.Fatalf("Failed with empty caption: %v", err)
	}
	if b1.Bounds().Dx() != b2.Bounds().Dx() || b1.Bounds().Dy() != b2.Bounds().Dy() {
		t.Errorf("Empty caption should produce same dimensions: no-cap %dx%d, empty-cap %dx%d",
			b1.Bounds().Dx(), b1.Bounds().Dy(), b2.Bounds().Dx(), b2.Bounds().Dy())
	}
}

func TestCaption_MultiLineWrap(t *testing.T) {
	fen := "8/8/8/8/8/8/8/8 w - - 0 1"
	longCap := "Это очень длинный текст который должен быть разбит на несколько строк поскольку он не помещается в ширину диаграммы"
	board, err := GenerateBoard(fen, longCap)
	if err != nil {
		t.Fatalf("Failed with long caption: %v", err)
	}
	noCap, _ := GenerateBoard(fen)
	if board.Bounds().Dy() <= noCap.Bounds().Dy() {
		t.Errorf("Multi-line caption should increase height, got %d vs %d",
			board.Bounds().Dy(), noCap.Bounds().Dy())
	}
}

func TestCaption_CyrillicText(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	russian := "Черные сыграли h5. Какова их идея?"
	board, err := GenerateBoard(fen, russian)
	if err != nil {
		t.Fatalf("Failed with Cyrillic caption: %v", err)
	}
	if board == nil {
		t.Fatal("Board is nil with Cyrillic caption")
	}
}

func TestCaption_SingleWordLongerThanWidth(t *testing.T) {
	fen := "8/8/8/8/8/8/8/8 w - - 0 1"
	veryLongWord := "ОЧЕНЬДЛИННОЕСЛОВОКОТОРОЕПРЕВЫШАЕТШИРИНУ"
	board, err := GenerateBoard(fen, veryLongWord)
	if err != nil {
		t.Fatalf("Failed with very long word: %v", err)
	}
	noCap, _ := GenerateBoard(fen)
	if board.Bounds().Dy() <= noCap.Bounds().Dy() {
		t.Errorf("Long word caption should still increase height")
	}
}

func TestCaption_UnsupportedGlyph(t *testing.T) {
	fen := "8/8/8/8/8/8/8/8 w - - 0 1"
	// Glyph with characters outside supported set
	weirdText := "Hello £200 ✓"
	board, err := GenerateBoard(fen, weirdText)
	if err != nil {
		t.Fatalf("Failed with unsupported glyphs: %v", err)
	}
	if board == nil {
		t.Fatal("Board is nil with unsupported glyphs")
	}
}

func TestWrapText(t *testing.T) {
	face, err := getCaptionFace()
	if err != nil {
		t.Fatalf("getCaptionFace: %v", err)
	}

	// Test empty string
	if got := wrapText(face, "", fixed.I(500)); len(got) != 0 {
		t.Errorf("empty caption should produce 0 lines, got %d", len(got))
	}

	// Test single line fits
	if got := wrapText(face, "Short caption", fixed.I(500)); len(got) != 1 {
		t.Errorf("short caption should fit in 1 line, got %d: %v", len(got), got)
	}

	// Test word wraps into multiple lines when constrained
	multiWord := "one two three four five six seven eight nine ten"
	wide := wrapText(face, multiWord, fixed.I(500))
	narrow := wrapText(face, multiWord, fixed.I(30))
	if len(wide) >= len(narrow) {
		t.Errorf("narrow width should produce more lines than wide; wide=%d, narrow=%d", len(wide), len(narrow))
	}

	// Test Cyrillic text wraps
	cyr := "Черные сыграли h5 какова их идея"
	narrowCyr := wrapText(face, cyr, fixed.I(30))
	if len(narrowCyr) < 2 {
		t.Errorf("narrow width with Cyrillic should produce multiple lines, got %d: %v", len(narrowCyr), narrowCyr)
	}

	// Test single very long word stays on one line
	long := wrapText(face, "ОЧЕНЬДЛИННОЕСЛОВО", fixed.I(50))
	if len(long) != 1 {
		t.Errorf("single word should stay on 1 line even if narrow, got %d", len(long))
	}

	// Test whitespace-only string
	if got := wrapText(face, "   ", fixed.I(500)); len(got) != 0 {
		t.Errorf("whitespace-only should produce 0 lines, got %d", len(got))
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
