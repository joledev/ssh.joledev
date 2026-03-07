package art

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var spotifyTrackRegex = regexp.MustCompile(`track/([a-zA-Z0-9]+)`)

// FetchSpotifyCode downloads a Spotify Code barcode and renders it as braille art.
func FetchSpotifyCode(trackURL string, width int) (string, error) {
	match := spotifyTrackRegex.FindStringSubmatch(trackURL)
	if match == nil {
		return "", fmt.Errorf("no track ID found in URL")
	}
	trackID := match[1]

	codeURL := fmt.Sprintf(
		"https://scannables.scdn.co/uri/plain/png/FFFFFF/black/640/spotify:track:%s",
		trackID,
	)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(codeURL)
	if err != nil {
		return "", fmt.Errorf("spotify code download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("spotify code status: %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", fmt.Errorf("spotify code decode: %w", err)
	}

	return spotifyCodeToBraille(img, width), nil
}

// spotifyCodeToBraille renders a Spotify Code image as braille art optimized for scanability.
func spotifyCodeToBraille(img image.Image, width int) string {
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
	enhanceContrast(gray, pixW, pixH, 1.5)
	sharpen(gray, pixW, pixH, 2.0)

	for y := 0; y < pixH; y++ {
		for x := 0; x < pixW; x++ {
			if gray[y][x] < 128 {
				gray[y][x] = 0
			} else {
				gray[y][x] = 255
			}
		}
	}

	var lines []string
	for cy := 0; cy < charH; cy++ {
		var line []rune
		for cx := 0; cx < charW; cx++ {
			code := 0
			for _, dot := range dotMap {
				dx, dy, bit := dot[0], dot[1], dot[2]
				px := cx*2 + dx
				py := cy*4 + dy
				if px >= pixW || py >= pixH {
					continue
				}
				if gray[py][px] < 128 {
					code |= bit
				}
			}
			if code == 0 {
				line = append(line, ' ')
			} else {
				line = append(line, rune(brailleBase+code))
			}
		}
		lines = append(lines, string(line))
	}

	return strings.Join(lines, "\n")
}

// sharpen applies an unsharp mask to enhance edges.
func sharpen(gray [][]float64, w, h int, amount float64) {
	blurred := make([][]float64, h)
	for y := 0; y < h; y++ {
		blurred[y] = make([]float64, w)
		for x := 0; x < w; x++ {
			var sum float64
			var count int
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					ny, nx := y+dy, x+dx
					if ny >= 0 && ny < h && nx >= 0 && nx < w {
						sum += gray[ny][nx]
						count++
					}
				}
			}
			blurred[y][x] = sum / float64(count)
		}
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			v := gray[y][x] + amount*(gray[y][x]-blurred[y][x])
			gray[y][x] = math.Max(0, math.Min(255, v))
		}
	}
}
