package tui

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"sync"
)

//go:embed assets/Llavero.jpg
var logoBytes []byte

var (
	renderedLogoOnce sync.Once
	cachedLogo       string
)

func colorToRGB(c color.Color) (uint8, uint8, uint8, uint8) {
	r, g, b, a := c.RGBA()
	return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)
}

func isBackgroundPixel(r, g, b uint8) bool {
	diffRG := math.Abs(float64(r) - float64(g))
	diffRB := math.Abs(float64(r) - float64(b))
	diffGB := math.Abs(float64(g) - float64(b))
	maxDiff := max(diffRG, max(diffRB, diffGB))
	// Neutral grays (checkerboard / drop shadows) have very low channel variance
	if maxDiff <= 5 || (r < 20 && g < 20 && b < 20) {
		return true
	}
	return false
}

func findKeyBounds(img image.Image) image.Rectangle {
	b := img.Bounds()
	minX, minY := b.Max.X, b.Max.Y
	maxX, maxY := b.Min.X, b.Min.Y

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bCol, _ := colorToRGB(img.At(x, y))
			if !isBackgroundPixel(r, g, bCol) {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	if minX >= maxX || minY >= maxY {
		return img.Bounds()
	}
	return image.Rect(minX, minY, maxX+1, maxY+1)
}

func generateANSIArt(img image.Image, targetWidth int, leftIndent int) string {
	bounds := findKeyBounds(img)
	cropMinX := bounds.Min.X
	cropMinY := bounds.Min.Y
	cropMaxX := bounds.Max.X
	cropMaxY := bounds.Max.Y

	cropW := cropMaxX - cropMinX
	cropH := cropMaxY - cropMinY

	if cropW <= 0 || cropH <= 0 {
		return ""
	}

	aspect := float64(cropW) / float64(cropH)
	subpixelHeight := int(math.Round(float64(targetWidth) / aspect))
	if subpixelHeight%2 != 0 {
		subpixelHeight++
	}
	charHeight := subpixelHeight / 2

	type Pixel struct {
		r, g, b uint8
		isBg    bool
	}

	pixels := make([][]Pixel, subpixelHeight)
	for py := 0; py < subpixelHeight; py++ {
		pixels[py] = make([]Pixel, targetWidth)
		for px := 0; px < targetWidth; px++ {
			srcX0 := cropMinX + (px * cropW / targetWidth)
			srcX1 := cropMinX + ((px + 1) * cropW / targetWidth)
			srcY0 := cropMinY + (py * cropH / subpixelHeight)
			srcY1 := cropMinY + ((py + 1) * cropH / subpixelHeight)

			if srcX1 <= srcX0 {
				srcX1 = srcX0 + 1
			}
			if srcY1 <= srcY0 {
				srcY1 = srcY0 + 1
			}

			var sumR, sumG, sumB, nonBgCount uint64
			for sy := srcY0; sy < srcY1; sy++ {
				for sx := srcX0; sx < srcX1; sx++ {
					r, g, b, _ := colorToRGB(img.At(sx, sy))
					if !isBackgroundPixel(r, g, b) {
						sumR += uint64(r)
						sumG += uint64(g)
						sumB += uint64(b)
						nonBgCount++
					}
				}
			}

			if nonBgCount == 0 {
				pixels[py][px] = Pixel{isBg: true}
			} else {
				pixels[py][px] = Pixel{
					r:    uint8(sumR / nonBgCount),
					g:    uint8(sumG / nonBgCount),
					b:    uint8(sumB / nonBgCount),
					isBg: false,
				}
			}
		}
	}

	indent := ""
	for i := 0; i < leftIndent; i++ {
		indent += " "
	}

	var buf bytes.Buffer
	for cy := 0; cy < charHeight; cy++ {
		buf.WriteString(indent)
		topY := cy * 2
		botY := cy*2 + 1

		for px := 0; px < targetWidth; px++ {
			topP := pixels[topY][px]
			botP := pixels[botY][px]

			if topP.isBg && botP.isBg {
				buf.WriteString(" ")
			} else if topP.isBg && !botP.isBg {
				fmt.Fprintf(&buf, "\x1b[38;2;%d;%d;%dm▄\x1b[0m", botP.r, botP.g, botP.b)
			} else if !topP.isBg && botP.isBg {
				fmt.Fprintf(&buf, "\x1b[38;2;%d;%d;%dm▀\x1b[0m", topP.r, topP.g, topP.b)
			} else {
				fmt.Fprintf(&buf, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀\x1b[0m", topP.r, topP.g, topP.b, botP.r, botP.g, botP.b)
			}
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

// GetLogo returns the high-fidelity truecolor ANSI rendered Llavero key logo.
func GetLogo() string {
	renderedLogoOnce.Do(func() {
		if len(logoBytes) == 0 {
			cachedLogo = ""
			return
		}

		img, _, err := image.Decode(bytes.NewReader(logoBytes))
		if err != nil {
			cachedLogo = ""
			return
		}

		// Generate a 48-char wide ANSI art key logo with 9 chars left indent for 68-char banner centering
		cachedLogo = generateANSIArt(img, 48, 9)
	})

	return cachedLogo
}
