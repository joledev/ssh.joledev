package art

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

const brailleBase = 0x2800

var dotMap = [8][3]int{
	{0, 0, 0x01}, {1, 0, 0x08},
	{0, 1, 0x02}, {1, 1, 0x10},
	{0, 2, 0x04}, {1, 2, 0x20},
	{0, 3, 0x40}, {1, 3, 0x80},
}

// ImageToBraille converts an image to braille art.
// When colorMode is false, outputs plain braille characters (white on terminal bg).
// When colorMode is true, outputs braille with ANSI true-color fg/bg per cell.
func ImageToBraille(img image.Image, width int, colorMode bool) string {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	charW := width
	charH := int(float64(charW) * float64(srcH) / float64(srcW) * 0.5)
	if charH < 1 {
		charH = 1
	}

	pixW := charW * 2
	pixH := charH * 4

	gray := resizeGray(img, pixW, pixH)
	enhanceContrast(gray, pixW, pixH, 1.8)

	var colorPix [][]color.RGBA
	if colorMode {
		colorPix = resizeColor(img, pixW, pixH)
	}

	dither(gray, pixW, pixH)

	var lines []string
	for cy := 0; cy < charH; cy++ {
		var line strings.Builder
		for cx := 0; cx < charW; cx++ {
			code := 0
			var fgR, fgG, fgB float64
			var bgR, bgG, bgB float64
			fgCount := 0
			bgCount := 0

			for _, dot := range dotMap {
				dx, dy, bit := dot[0], dot[1], dot[2]
				px := cx*2 + dx
				py := cy*4 + dy
				if px >= pixW || py >= pixH {
					continue
				}

				if gray[py][px] < 128 {
					code |= bit
					if colorMode {
						c := colorPix[py][px]
						fgR += float64(c.R)
						fgG += float64(c.G)
						fgB += float64(c.B)
					}
					fgCount++
				} else {
					if colorMode {
						c := colorPix[py][px]
						bgR += float64(c.R)
						bgG += float64(c.G)
						bgB += float64(c.B)
					}
					bgCount++
				}
			}

			char := rune(brailleBase + code)

			if !colorMode {
				if code == 0 {
					line.WriteRune(' ')
				} else {
					line.WriteRune(char)
				}
				continue
			}

			if code == 0 {
				if bgCount > 0 {
					r := uint8(bgR / float64(bgCount))
					g := uint8(bgG / float64(bgCount))
					b := uint8(bgB / float64(bgCount))
					line.WriteString(fmt.Sprintf("\033[48;2;%d;%d;%dm \033[0m", r, g, b))
				} else {
					line.WriteRune(' ')
				}
			} else if code == 0xFF {
				if fgCount > 0 {
					r := uint8(fgR / float64(fgCount))
					g := uint8(fgG / float64(fgCount))
					b := uint8(fgB / float64(fgCount))
					line.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m", r, g, b, char))
				} else {
					line.WriteRune(char)
				}
			} else {
				fr, fg, fb := uint8(0), uint8(0), uint8(0)
				if fgCount > 0 {
					fr = uint8(fgR / float64(fgCount))
					fg = uint8(fgG / float64(fgCount))
					fb = uint8(fgB / float64(fgCount))
				}
				br, bgg, bb := uint8(0), uint8(0), uint8(0)
				if bgCount > 0 {
					br = uint8(bgR / float64(bgCount))
					bgg = uint8(bgG / float64(bgCount))
					bb = uint8(bgB / float64(bgCount))
				}
				line.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%d;48;2;%d;%d;%dm%c\033[0m",
					fr, fg, fb, br, bgg, bb, char))
			}
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}

func resizeGray(img image.Image, w, h int) [][]float64 {
	bounds := img.Bounds()
	srcW := float64(bounds.Dx())
	srcH := float64(bounds.Dy())

	result := make([][]float64, h)
	for y := 0; y < h; y++ {
		result[y] = make([]float64, w)
		for x := 0; x < w; x++ {
			x0 := int(float64(x) * srcW / float64(w))
			y0 := int(float64(y) * srcH / float64(h))
			x1 := int(float64(x+1) * srcW / float64(w))
			y1 := int(float64(y+1) * srcH / float64(h))
			if x1 <= x0 {
				x1 = x0 + 1
			}
			if y1 <= y0 {
				y1 = y0 + 1
			}

			var sum float64
			count := 0
			for sy := y0; sy < y1 && sy < bounds.Dy(); sy++ {
				for sx := x0; sx < x1 && sx < bounds.Dx(); sx++ {
					r, g, b, _ := img.At(bounds.Min.X+sx, bounds.Min.Y+sy).RGBA()
					lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
					sum += lum / 257
					count++
				}
			}
			if count > 0 {
				result[y][x] = sum / float64(count)
			}
		}
	}
	return result
}

func resizeColor(img image.Image, w, h int) [][]color.RGBA {
	bounds := img.Bounds()
	srcW := float64(bounds.Dx())
	srcH := float64(bounds.Dy())

	result := make([][]color.RGBA, h)
	for y := 0; y < h; y++ {
		result[y] = make([]color.RGBA, w)
		for x := 0; x < w; x++ {
			x0 := int(float64(x) * srcW / float64(w))
			y0 := int(float64(y) * srcH / float64(h))
			x1 := int(float64(x+1) * srcW / float64(w))
			y1 := int(float64(y+1) * srcH / float64(h))
			if x1 <= x0 {
				x1 = x0 + 1
			}
			if y1 <= y0 {
				y1 = y0 + 1
			}

			var rSum, gSum, bSum float64
			count := 0
			for sy := y0; sy < y1 && sy < bounds.Dy(); sy++ {
				for sx := x0; sx < x1 && sx < bounds.Dx(); sx++ {
					r, g, b, _ := img.At(bounds.Min.X+sx, bounds.Min.Y+sy).RGBA()
					rSum += float64(r)
					gSum += float64(g)
					bSum += float64(b)
					count++
				}
			}
			if count > 0 {
				result[y][x] = color.RGBA{
					R: uint8(rSum / float64(count) / 257),
					G: uint8(gSum / float64(count) / 257),
					B: uint8(bSum / float64(count) / 257),
					A: 255,
				}
			}
		}
	}
	return result
}

func enhanceContrast(gray [][]float64, w, h int, factor float64) {
	var sum float64
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			sum += gray[y][x]
		}
	}
	mean := sum / float64(w*h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			v := (gray[y][x]-mean)*factor + mean
			if v < 0 {
				v = 0
			}
			if v > 255 {
				v = 255
			}
			gray[y][x] = v
		}
	}
}

func dither(gray [][]float64, w, h int) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			old := gray[y][x]
			var newVal float64
			if old < 128 {
				newVal = 0
			} else {
				newVal = 255
			}
			gray[y][x] = newVal
			err := old - newVal

			if x+1 < w {
				gray[y][x+1] += err * 7 / 16
			}
			if y+1 < h {
				if x-1 >= 0 {
					gray[y+1][x-1] += err * 3 / 16
				}
				gray[y+1][x] += err * 5 / 16
				if x+1 < w {
					gray[y+1][x+1] += err * 1 / 16
				}
			}
		}
	}
}
