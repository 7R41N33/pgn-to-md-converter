package fenlib

import (
	"fmt"
	"image"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	captionFontSize  = 10.0 // pt
	captionDPI       = 72.0
	captionLineSpace = 2  // px extra between lines
	captionPad       = 0  // px extra at the bottom below last line
)

var (
	captionFaceOnce sync.Once
	captionFace     font.Face
	captionFaceErr  error
)

func getCaptionFace() (font.Face, error) {
	captionFaceOnce.Do(func() {
		f, err := opentype.Parse(gomono.TTF)
		if err != nil {
			captionFaceErr = fmt.Errorf("parse gomono: %w", err)
			return
		}
		face, err := opentype.NewFace(f, &opentype.FaceOptions{
			Size:    captionFontSize,
			DPI:     captionDPI,
			Hinting: font.HintingNone,
		})
		if err != nil {
			captionFaceErr = fmt.Errorf("new face: %w", err)
			return
		}
		captionFace = face
	})
	return captionFace, captionFaceErr
}

func lineHeight(face font.Face) int {
	return face.Metrics().Height.Ceil() + captionLineSpace
}

func wrapText(face font.Face, caption string, maxPixels fixed.Int26_6) []string {
	if caption == "" {
		return nil
	}
	words := splitWords(caption)
	if len(words) == 0 {
		return nil
	}

	d := font.Drawer{Face: face}
	var lines []string
	line := words[0]

	for _, word := range words[1:] {
		candidate := line + " " + word
		w := d.MeasureString(candidate)
		if w > maxPixels {
			lines = append(lines, line)
			line = word
		} else {
			line = candidate
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

func splitWords(s string) []string {
	var words []string
	start := -1
	for i, ch := range s {
		if ch == ' ' {
			if start != -1 {
				words = append(words, s[start:i])
				start = -1
			}
		} else if start == -1 {
			start = i
		}
	}
	if start != -1 {
		words = append(words, s[start:])
	}
	return words
}

func drawCaption(img *image.RGBA, caption string, imgWidth, boardEnd int) error {
	if caption == "" {
		return nil
	}

	face, err := getCaptionFace()
	if err != nil {
		return err
	}

	maxPx := fixed.I(imgWidth)
	lines := wrapText(face, caption, maxPx)
	if len(lines) == 0 {
		return nil
	}

	lh := lineHeight(face)
	ascent := face.Metrics().Ascent.Ceil()
	y := boardEnd + ascent

	for _, line := range lines {
		d := font.Drawer{
			Dst:  img,
			Src:  image.Black,
			Face: face,
		}
		lineW := d.MeasureString(line)
		x := (fixed.I(imgWidth) - lineW) / 2
		d.Dot = fixed.Point26_6{X: x, Y: fixed.I(y)}
		d.DrawString(line)
		y += lh
	}
	return nil
}
