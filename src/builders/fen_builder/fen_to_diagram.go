package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

const (
	pieceOrder = "KQRBNPkqrbnp"
	spritePath = "src/images/Chess_Pieces_Sprite.svg.png"
	defaultOut = "tmp/images/diagram.png"
)

const (
	FENPosition    = 0
	FENActiveColor = 1
	FENCastling    = 2
	FENEnPassant   = 3
	FENHalfmove    = 4
	FENFullmove    = 5
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <fen_string> [output_path]\n", os.Args[0])
		os.Exit(1)
	}

	fen := os.Args[1]
	outputPath := defaultOut
	if len(os.Args) >= 3 {
		outputPath = os.Args[2]
	}

	if err := validateFEN(fen); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid FEN - %v\n", err)
		os.Exit(1)
	}

	if err := generateDiagram(fen, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating diagram: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Diagram saved to: %s\n", outputPath)
}

func validateFEN(fen string) error {
	fields := strings.Fields(fen)
	if len(fields) < 6 {
		return fmt.Errorf("FEN must have 6 fields, got %d", len(fields))
	}

	position := fields[FENPosition]
	rows := strings.Split(position, "/")
	if len(rows) != 8 {
		return fmt.Errorf("position must have 8 rows, got %d", len(rows))
	}

	for _, row := range rows {
		if row == "" {
			return fmt.Errorf("empty row in position")
		}
		count := 0
		for _, ch := range row {
			if ch >= '1' && ch <= '8' {
				count += int(ch - '0')
			} else if strings.ContainsRune("KQRBNPkqrbnp", ch) {
				count++
			} else {
				return fmt.Errorf("invalid character in position: %c", ch)
			}
		}
		if count != 8 {
			return fmt.Errorf("row must have 8 squares, got %d", count)
		}
	}

	if fields[FENActiveColor] != "w" && fields[FENActiveColor] != "b" {
		return fmt.Errorf("active color must be 'w' or 'b'")
	}

	if !isValidCastling(fields[FENCastling]) {
		return fmt.Errorf("invalid castling field")
	}

	if fields[FENEnPassant] != "-" && !isValidEnPassant(fields[FENEnPassant]) {
		return fmt.Errorf("invalid en passant square")
	}

	return nil
}

func isValidCastling(s string) bool {
	if s == "-" {
		return true
	}
	for _, ch := range s {
		if !strings.ContainsRune("KQkq", ch) {
			return false
		}
	}
	return true
}

func isValidEnPassant(s string) bool {
	if s == "-" {
		return true
	}
	if len(s) != 2 {
		return false
	}
	return s[0] >= 'a' && s[0] <= 'h' && (s[1] == '3' || s[1] == '6')
}

func generateDiagram(fen, outputPath string) error {
	spritePaths := []string{
		spritePath,
		"../../images/Chess_Pieces_Sprite.svg.png",
		"../images/Chess_Pieces_Sprite.svg.png",
		"images/Chess_Pieces_Sprite.svg.png",
	}

	var sprite image.Image
	var err error
	for _, path := range spritePaths {
		sprite, err = loadSprite(path)
		if err == nil {
			break
		}
	}
	if sprite == nil {
		return fmt.Errorf("failed to load sprite: %v", err)
	}

	spriteBounds := sprite.Bounds()
	pieceWidth := spriteBounds.Dx() / 6
	pieceHeight := spriteBounds.Dy() / 2

	cellSize := 60
	boardSize := 8 * cellSize
	// Add extra space for coordinates: 20px left for numbers, 20px bottom for letters
	imgWidth := boardSize + 20
	board := image.NewRGBA(image.Rect(20, 0, imgWidth, boardSize))

	lightColor := color.RGBA{240, 217, 181, 255}
	darkColor := color.RGBA{181, 136, 99, 255}

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			c := lightColor
			if (row+col)%2 == 1 {
				c = darkColor
			}
			rect := image.Rect(20+col*cellSize, row*cellSize, 20+(col+1)*cellSize, (row+1)*cellSize)
			draw.Draw(board, rect, &image.Uniform{c}, image.Point{}, draw.Src)
		}
	}
	
	// Add coordinates: a-h on bottom, 1-8 on left
	addCoordinates(board, cellSize)

	fields := strings.Fields(fen)
	position := fields[FENPosition]
	rows := strings.Split(position, "/")

	for i, row := range rows {
		col := 0
		for _, ch := range row {
			if ch >= '1' && ch <= '8' {
				col += int(ch - '0')
			} else {
				pieceImg, err := extractPiece(sprite, pieceWidth, pieceHeight, ch)
				if err == nil {
					destRect := image.Rect(col*cellSize, i*cellSize, (col+1)*cellSize, (i+1)*cellSize)
					drawPiece(board, pieceImg, destRect, pieceWidth, pieceHeight, cellSize)
				}
				col++
			}
		}
	}

	os.MkdirAll(filepath.Dir(outputPath), 0755)
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, board)
}

func loadSprite(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func extractPiece(sprite image.Image, pieceWidth, pieceHeight int, piece rune) (image.Image, error) {
	idx := strings.IndexRune(pieceOrder, piece)
	if idx == -1 {
		return nil, fmt.Errorf("unknown piece: %c", piece)
	}

	row := idx / 6
	col := idx % 6

	img := image.NewRGBA(image.Rect(0, 0, pieceWidth, pieceHeight))

	for x := 0; x < pieceWidth; x++ {
		for y := 0; y < pieceHeight; y++ {
			img.Set(x, y, sprite.At(col*pieceWidth+x, row*pieceHeight+y))
		}
	}

	return img, nil
}

func drawPiece(board *image.RGBA, piece image.Image, destRect image.Rectangle, srcW, srcH, destSize int) {
	for x := 0; x < destSize; x++ {
		for y := 0; y < destSize; y++ {
			srcX := x * srcW / destSize
			srcY := y * srcH / destSize
			board.Set(destRect.Min.X+x, destRect.Min.Y+y, piece.At(srcX, srcY))
		}
	}
}

func addCoordinates(board draw.Image, cellSize int) {
	// Draw letters a-h on bottom (y = 8*cellSize + 10)
	letters := "abcdefgh"
	for i := 0; i < 8; i++ {
		x := 20 + i*cellSize + cellSize/2 - 4 // Center the letter, accounting for 20px left offset
		y := 8*cellSize + 15 // Below the board
		drawLetter(board, x, y, letters[i])
	}

	// Draw numbers 1-8 on left (x = 5)
	for i := 0; i < 8; i++ {
		x := 5 // Left of the board
		y := (7-i)*cellSize + cellSize/2 + 4 // Center vertically, counting from top
		drawNumber(board, x, y, i+1)
	}
}

func drawLetter(board draw.Image, x, y int, letter byte) {
	// Simple bitmap for letters (5x7 pixel font)
	fontMap := map[byte][]string{
		'a': {"  *  ", " * * ", "*****", "*   *", "*   *", "*   *", "     "},
		'b': {"**** ", "*   *", "*   *", "**** ", "*   *", "*   *", "**** "},
		'c': {" ****", "*    ", "*    ", "*    ", "*    ", "*    ", " ****"},
		'd': {"**** ", "*   *", "*   *", "*   *", "*   *", "*   *", "**** "},
		'e': {"*****", "*    ", "*    ", "*****", "*    ", "*    ", "*****"},
		'f': {"*****", "*    ", "*    ", "**** ", "*    ", "*    ", "*    "},
		'g': {" ****", "*    ", "*    ", "* ***", "*   *", "*   *", " ****"},
		'h': {"*   *", "*   *", "*   *", "*****", "*   *", "*   *", "*   *"},
	}
	
	if glyph, ok := fontMap[letter]; ok {
		for gy, row := range glyph {
			for gx, ch := range row {
				if ch == '*' {
					board.Set(x+gx, y+gy, color.Black)
				}
			}
		}
	}
}

func drawNumber(board draw.Image, x, y int, num int) {
	str := fmt.Sprintf("%d", num)
	if len(str) == 1 {
		drawDigit(board, x, y, str[0])
	} else if len(str) == 2 {
		drawDigit(board, x, y, str[0])
		drawDigit(board, x+5, y, str[1])
	}
}

func drawDigit(board draw.Image, x, y int, digit byte) {
	fontMap := map[byte][]string{
		'0': {" *** ", "*   *", "*  **", "* * *", "**  *", "*   *", " *** "},
		'1': {"  *  ", " * * ", "   * ", "   * ", "   * ", "   * ", "*****"},
		'2': {" *** ", "*   *", "    *", "  ** ", " *   ", "*    ", "*****"},
		'3': {" *** ", "*   *", "    *", "  ** ", "    *", "*   *", " *** "},
		'4': {"*   *", "*   *", "*   *", "*****", "    *", "    *", "    *"},
		'5': {"*****", "*    ", "*    ", "**** ", "    *", "*   *", " *** "},
		'6': {" *** ", "*    ", "*    ", "**** ", "*   *", "*   *", " *** "},
		'7': {"*****", "    *", "   * ", "  *  ", "  *  ", "  *  ", "  *  "},
		'8': {" *** ", "*   *", "*   *", " *** ", "*   *", "*   *", " *** "},
		'9': {" *** ", "*   *", "*   *", " ****", "    *", "*   *", " *** "},
	}
	
	if glyph, ok := fontMap[digit]; ok {
		for gy, row := range glyph {
			for gx, ch := range row {
				if ch == '*' {
					board.Set(x+gx, y+gy, color.Black)
				}
			}
		}
	}
}
