package art

import (
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// GenerateQR creates a QR code rendered in unicode block characters.
// Uses upper/lower half blocks to fit 2 rows per line.
func GenerateQR(content string) (string, error) {
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return "", err
	}
	qr.DisableBorder = true
	bitmap := qr.Bitmap()

	rows := len(bitmap)
	cols := 0
	if rows > 0 {
		cols = len(bitmap[0])
	}

	var lines []string
	for y := 0; y < rows; y += 2 {
		var line strings.Builder
		for x := 0; x < cols; x++ {
			top := bitmap[y][x]
			bottom := false
			if y+1 < rows {
				bottom = bitmap[y+1][x]
			}

			// black = true in bitmap = dark module
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

	return strings.Join(lines, "\n"), nil
}
