package art

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var spotifyTrackRegex = regexp.MustCompile(`track/([a-zA-Z0-9]+)`)

// FetchSpotifyCode downloads a Spotify Code barcode and renders it using half-block characters.
func FetchSpotifyCode(trackURL string, width, height int) (string, error) {
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

	return imageToBlocks(img, width, height), nil
}

// imageToBlocks renders an image using half-block unicode characters.
// Each character represents 2 vertical pixels using ▀▄█ and space.
// Much better than braille for geometric/barcode images.
func imageToBlocks(img image.Image, width, rows int) string {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	charW := width
	charH := rows * 2 // 2 pixels per row of half-blocks
	if charH == 0 {
		charH = int(float64(charW) * float64(srcH) / float64(srcW) * 2)
		if charH%2 != 0 {
			charH++
		}
	}

	// Resize to target pixel dimensions
	pixels := make([][]bool, charH)
	for y := 0; y < charH; y++ {
		pixels[y] = make([]bool, charW)
		for x := 0; x < charW; x++ {
			// Area-average sampling
			x0 := x * srcW / charW
			y0 := y * srcH / charH
			x1 := (x + 1) * srcW / charW
			y1 := (y + 1) * srcH / charH
			if x1 <= x0 {
				x1 = x0 + 1
			}
			if y1 <= y0 {
				y1 = y0 + 1
			}

			var lum float64
			count := 0
			for sy := y0; sy < y1 && sy < srcH; sy++ {
				for sx := x0; sx < x1 && sx < srcW; sx++ {
					r, g, b, _ := img.At(bounds.Min.X+sx, bounds.Min.Y+sy).RGBA()
					lum += 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
					count++
				}
			}
			if count > 0 {
				lum /= float64(count) * 257
			}
			pixels[y][x] = lum < 128 // dark = true
		}
	}

	var lines []string
	for y := 0; y < charH; y += 2 {
		var line strings.Builder
		for x := 0; x < charW; x++ {
			top := pixels[y][x]
			bottom := false
			if y+1 < charH {
				bottom = pixels[y+1][x]
			}

			switch {
			case top && bottom:
				line.WriteRune('█')
			case top && !bottom:
				line.WriteRune('▀')
			case !top && bottom:
				line.WriteRune('▄')
			default:
				line.WriteRune(' ')
			}
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}
