package art

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"regexp"
	"time"
)

var spotifyTrackRegex = regexp.MustCompile(`track/([a-zA-Z0-9]+)`)

// FetchSpotifyCode downloads a Spotify Code barcode and converts it to braille art.
// The image is inverted so the bars render as white braille dots on dark terminal bg.
func FetchSpotifyCode(trackURL string, width int) (string, error) {
	match := spotifyTrackRegex.FindStringSubmatch(trackURL)
	if match == nil {
		return "", fmt.Errorf("no track ID found in URL")
	}
	trackID := match[1]

	// Request white bg with black bars, then invert for terminal display
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

	inverted := invertImage(img)
	return ImageToBraille(inverted, width, false), nil
}

// invertImage flips all pixel values (white becomes black, black becomes white).
func invertImage(img image.Image) image.Image {
	bounds := img.Bounds()
	inverted := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			inverted.Set(x, y, color.RGBA{
				R: uint8(255 - r>>8),
				G: uint8(255 - g>>8),
				B: uint8(255 - b>>8),
				A: uint8(a >> 8),
			})
		}
	}
	return inverted
}
