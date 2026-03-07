package art

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"time"
)

type oembedResponse struct {
	ThumbnailURL string `json:"thumbnail_url"`
}

// FetchCoverImage downloads an album cover image from Spotify via oEmbed.
func FetchCoverImage(trackURL string) (image.Image, error) {
	if trackURL == "" {
		return nil, fmt.Errorf("no track URL")
	}

	client := &http.Client{Timeout: 10 * time.Second}

	oembedURL := fmt.Sprintf("https://open.spotify.com/oembed?url=%s", trackURL)
	resp, err := client.Get(oembedURL)
	if err != nil {
		return nil, fmt.Errorf("oembed request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("oembed status: %d", resp.StatusCode)
	}

	var oembed oembedResponse
	if err := json.NewDecoder(resp.Body).Decode(&oembed); err != nil {
		return nil, fmt.Errorf("oembed decode: %w", err)
	}

	if oembed.ThumbnailURL == "" {
		return nil, fmt.Errorf("no thumbnail in oembed")
	}

	imgResp, err := client.Get(oembed.ThumbnailURL)
	if err != nil {
		return nil, fmt.Errorf("image download: %w", err)
	}
	defer imgResp.Body.Close()

	img, _, err := image.Decode(imgResp.Body)
	if err != nil {
		return nil, fmt.Errorf("image decode: %w", err)
	}

	return img, nil
}
