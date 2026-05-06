package fenlib

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Sprite order: K, Q, B, N, R, P for white (row 0), k, q, b, n, r, p for black (row 1)
	// Each piece is 213x213 pixels in a 1280x427 sprite (6 pieces x 2 rows)
	pieceOrder = "KQBNRPkqbnrp"
	spritePath = "src/images/Chess_Pieces_Sprite.svg.png"
)

const (
	FENPosition    = 0
	FENActiveColor = 1
	FENCastling    = 2
	FENEnPassant   = 3
	FENHalfmove    = 4
	FENFullmove    = 5
)

// ValidateFEN validates a FEN string
func ValidateFEN(fen string) error {
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

// GenerateBoard generates a chess board image from a FEN string
func GenerateBoard(fen string) (*image.RGBA, error) {
	if err := ValidateFEN(fen); err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("failed to load sprite: %v", err)
	}

	spriteBounds := sprite.Bounds()
	pieceWidth := spriteBounds.Dx() / 6
	pieceHeight := spriteBounds.Dy() / 2

	cellSize := 60
	boardSize := 8 * cellSize
	borderSize := cellSize / 2 // Half cell border around board

	// Image size: board + border on all sides
	imgSize := boardSize + 2*borderSize
	board := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))

	// Fill entire image with white (border)
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(board, board.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	// Draw board cells starting at (borderSize, borderSize)
	lightColor := color.RGBA{240, 217, 181, 255}
	darkColor := color.RGBA{181, 136, 99, 255}

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			c := lightColor
			if (row+col)%2 == 1 {
				c = darkColor
			}
			x0 := borderSize + col*cellSize
			y0 := borderSize + row*cellSize
			rect := image.Rect(x0, y0, x0+cellSize, y0+cellSize)
			draw.Draw(board, rect, &image.Uniform{c}, image.Point{}, draw.Src)
		}
	}

	// Add coordinates on the white border
	addCoordinates(board, cellSize, borderSize)

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
					x0 := borderSize + col*cellSize
					y0 := borderSize + i*cellSize
					destRect := image.Rect(x0, y0, x0+cellSize, y0+cellSize)
					drawPiece(board, pieceImg, destRect, pieceWidth, pieceHeight, cellSize)
				}
				col++
			}
		}
	}

	return board, nil
}

// GenerateDiagram generates a chess diagram and saves it to a file
func GenerateDiagram(fen, outputPath string) error {
	board, err := GenerateBoard(fen)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(outputPath), 0755)
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, board)
}

// GenerateDiagramBase64 generates a chess diagram and returns it as a base64-encoded PNG string
func GenerateDiagramBase64(fen string) (string, error) {
	board, err := GenerateBoard(fen)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, board); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// GenerateDiagramToWriter generates a chess diagram and writes PNG to the provided writer
func GenerateDiagramToWriter(fen string, w io.Writer) error {
	board, err := GenerateBoard(fen)
	if err != nil {
		return err
	}

	return png.Encode(w, board)
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
	margin := 5 // Reduce piece size by 5 pixels on each side
	innerSize := destSize - 2*margin

	if innerSize <= 0 {
		return
	}

	for x := 0; x < innerSize; x++ {
		for y := 0; y < innerSize; y++ {
			srcX := x * srcW / innerSize
			srcY := y * srcH / innerSize
			pieceColor := piece.At(srcX, srcY)

			// Only draw non-transparent pixels (alpha > 0)
			r, g, b, a := pieceColor.RGBA()
			if a > 0 {
				board.Set(destRect.Min.X+margin+x, destRect.Min.Y+margin+y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
			}
		}
	}
}

func addCoordinates(board draw.Image, cellSize, borderSize int) {
	boardSize := 8 * cellSize

	// Draw letters a-h on bottom border
	letters := "abcdefgh"
	for i := 0; i < 8; i++ {
		x := borderSize + i*cellSize + cellSize/2 - 3 // Center horizontally in cell
		y := boardSize + borderSize + borderSize/2 - 4 // Center vertically in bottom border
		drawLetter(board, x, y, letters[i])
	}

	// Draw numbers 1-8 on left border (1 at bottom, 8 at top)
	for i := 0; i < 8; i++ {
		x := borderSize/2 - 3 // Center horizontally in left border
		y := borderSize + (7-i)*cellSize + cellSize/2 - 4 // Center vertically, 1 at bottom
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
