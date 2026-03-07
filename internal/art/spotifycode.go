package art

import (
	"fmt"
	"image"
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

	// Black bars → dark pixels → braille dots (visible white on terminal)
	// White bg → light pixels → spaces (transparent, terminal bg)
	return ImageToBraille(img, width, false), nil
}
