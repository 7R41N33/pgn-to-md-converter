package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
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
	fen := "8/8/8/8/8/8/8/8 w - - 0 1"
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

func TestGenerateDiagram_ImageSizeWithBorder(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	outputPath := "/tmp/test_border_size.png"

	if err := generateDiagram(fen, outputPath); err != nil {
		t.Errorf("Failed to generate diagram: %v", err)
	}

	// Load the generated image and check dimensions
	file, err := os.Open(outputPath)
	if err != nil {
		t.Errorf("Failed to open output file: %v", err)
	}
	defer file.Close()

	img, err := decodeImage(file)
	if err != nil {
		t.Errorf("Failed to decode image: %v", err)
	}

	cellSize := 60
	borderSize := cellSize / 2
	expectedSize := 8*cellSize + 2*borderSize

	if img.Bounds().Dx() != expectedSize || img.Bounds().Dy() != expectedSize {
		t.Errorf("Expected image size %dx%d, got %dx%d", expectedSize, expectedSize, img.Bounds().Dx(), img.Bounds().Dy())
	}

	os.Remove(outputPath)
}

func TestGenerateDiagram_WhiteBorder(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	outputPath := "/tmp/test_white_border.png"

	if err := generateDiagram(fen, outputPath); err != nil {
		t.Errorf("Failed to generate diagram: %v", err)
	}

	file, err := os.Open(outputPath)
	if err != nil {
		t.Errorf("Failed to open output file: %v", err)
	}
	defer file.Close()

	img, err := decodeImage(file)
	if err != nil {
		t.Errorf("Failed to decode image: %v", err)
	}

	// Check that border pixels are white
	borderColor := img.At(5, 5) // Top-left corner should be white (border)

	r, g, b, _ := borderColor.RGBA()
	if r != 65535 || g != 65535 || b != 65535 {
		t.Errorf("Expected white border, got color at (5,5): %v", borderColor)
	}

	os.Remove(outputPath)
}

func TestDrawPiece_NoTransparency(t *testing.T) {
	// Create a test board
	cellSize := 60
	borderSize := cellSize / 2
	imgSize := 8*cellSize + 2*borderSize
	board := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))

	// Fill with dark color
	darkColor := color.RGBA{181, 136, 99, 255}
	draw.Draw(board, board.Bounds(), &image.Uniform{darkColor}, image.Point{}, draw.Src)

	// Create a simple piece image (just a solid color for testing)
	pieceImg := image.NewRGBA(image.Rect(0, 0, 60, 60))
	pieceColor := color.RGBA{0, 0, 0, 255}
	draw.Draw(pieceImg, pieceImg.Bounds(), &image.Uniform{pieceColor}, image.Point{}, draw.Src)

	// Draw piece
	destRect := image.Rect(borderSize, borderSize, borderSize+cellSize, borderSize+cellSize)
	drawPiece(board, pieceImg, destRect, 60, 60, cellSize)

	// Check that the cell still has the piece drawn (not transparent)
	// The piece should be on top of the dark cell
	centerColor := board.At(borderSize+cellSize/2, borderSize+cellSize/2)
	r, g, b, _ := centerColor.RGBA()

	// Should be black (piece color), not dark brown (cell color)
	if r == 181*257 && g == 136*257 && b == 99*257 {
		t.Error("Cell appears to be transparent - showing cell color instead of piece")
	}
}

func TestAddCoordinates_Position(t *testing.T) {
	cellSize := 60
	borderSize := cellSize / 2
	imgSize := 8*cellSize + 2*borderSize
	board := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))

	// Fill with white
	draw.Draw(board, board.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	addCoordinates(board, cellSize, borderSize)

	// Check that coordinates were drawn (pixels should be black)
	// The letter 'a' should be drawn in the bottom border area
	// We can't easily check exact position without knowing the font, but we can check
	// that some black pixels exist in the border area
	foundBlackPixel := false
	for y := 8*cellSize + borderSize; y < imgSize; y++ {
		for x := borderSize; x < borderSize+cellSize; x++ {
			r, _, _, _ := board.At(x, y).RGBA()
			if r == 0 {
				foundBlackPixel = true
				break
			}
		}
		if foundBlackPixel {
			break
		}
	}

	if !foundBlackPixel {
		t.Error("Expected coordinate letters to be drawn on bottom border")
	}
}

func TestGenerateDiagramBase64_ValidFEN(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	b64, err := generateDiagramBase64(fen)
	if err != nil {
		t.Errorf("Failed to generate base64 diagram: %v", err)
	}

	if b64 == "" {
		t.Error("Base64 output should not be empty")
	}

	// Verify it's valid base64
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		t.Errorf("Output is not valid base64: %v", err)
	}

	if len(decoded) == 0 {
		t.Error("Decoded output should not be empty")
	}

	// Verify it starts with PNG signature (89 50 4E 47)
	if len(decoded) < 4 || decoded[0] != 0x89 || decoded[1] != 0x50 || decoded[2] != 0x4E || decoded[3] != 0x47 {
		t.Error("Decoded output is not a valid PNG file")
	}
}

func TestGenerateDiagramBase64_InvalidFEN(t *testing.T) {
	fen := "invalid fen"

	_, err := generateDiagramBase64(fen)
	if err == nil {
		t.Error("Should return error for invalid FEN")
	}
}

func TestGenerateDiagramBase64_EmptyBoard(t *testing.T) {
	fen := "8/8/8/8/8/8/8/8 w - - 0 1"

	b64, err := generateDiagramBase64(fen)
	if err != nil {
		t.Errorf("Failed to generate base64 diagram: %v", err)
	}

	if b64 == "" {
		t.Error("Base64 output should not be empty")
	}

	// Decode and verify dimensions
	decoded, _ := base64.StdEncoding.DecodeString(b64)
	img, err := png.Decode(bytes.NewReader(decoded))
	if err != nil {
		t.Errorf("Failed to decode PNG: %v", err)
	}

	cellSize := 60
	borderSize := cellSize / 2
	expectedSize := 8*cellSize + 2*borderSize

	if img.Bounds().Dx() != expectedSize || img.Bounds().Dy() != expectedSize {
		t.Errorf("Expected image size %dx%d, got %dx%d", expectedSize, expectedSize, img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestGenerateBoard_ValidFEN(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	board, err := generateBoard(fen)
	if err != nil {
		t.Errorf("Failed to generate board: %v", err)
	}

	if board == nil {
		t.Error("Board should not be nil")
	}

	cellSize := 60
	borderSize := cellSize / 2
	expectedSize := 8*cellSize + 2*borderSize

	if board.Bounds().Dx() != expectedSize || board.Bounds().Dy() != expectedSize {
		t.Errorf("Expected board size %dx%d, got %dx%d", expectedSize, expectedSize, board.Bounds().Dx(), board.Bounds().Dy())
	}
}

func TestGenerateBoard_InvalidFEN(t *testing.T) {
	fen := "invalid fen"

	_, err := generateBoard(fen)
	if err == nil {
		t.Error("Should return error for invalid FEN")
	}
}

func TestGenerateBoard_WhiteBorder(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	board, err := generateBoard(fen)
	if err != nil {
		t.Errorf("Failed to generate board: %v", err)
	}

	// Check that border pixels are white
	borderColor := board.At(5, 5) // Top-left corner should be white (border)

	r, g, b, _ := borderColor.RGBA()
	if r != 65535 || g != 65535 || b != 65535 {
		t.Errorf("Expected white border, got color at (5,5): %v", borderColor)
	}
}

func decodeImage(file *os.File) (image.Image, error) {
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}
