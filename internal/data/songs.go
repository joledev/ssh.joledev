package data

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Song struct {
	Artist string
	Title  string
	Album  string
	Year   string
	URL    string
}

func dataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "data")
}

func LoadSongs(customPath string) ([]Song, error) {
	path := customPath
	if path == "" {
		path = filepath.Join(dataDir(), "songs.txt")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var songs []Song
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "|", 5)
		if len(parts) < 2 {
			continue
		}
		s := Song{
			Artist: parts[0],
			Title:  parts[1],
		}
		if len(parts) > 2 {
			s.Album = parts[2]
		}
		if len(parts) > 3 {
			s.Year = parts[3]
		}
		if len(parts) > 4 {
			s.URL = parts[4]
		}
		songs = append(songs, s)
	}
	return songs, scanner.Err()
}

// SongOfTheDay returns a deterministic "random" song for the current day
// based on the date in America/Tijuana timezone.
func SongOfTheDay(songs []Song) Song {
	if len(songs) == 0 {
		return Song{
			Artist: "Asian Kung-Fu Generation",
			Title:  "Shinseiki No Love Song",
			Album:  "Magic Disk",
			Year:   "2010",
			URL:    "https://youtu.be/xyY4IZ3JDFE",
		}
	}

	loc, err := time.LoadLocation("America/Tijuana")
	if err != nil {
		loc = time.FixedZone("PST", -8*60*60)
	}
	now := time.Now().In(loc)
	dateStr := now.Format("2006-01-02")

	// Hash the date to get a deterministic index
	h := sha256.Sum256([]byte(dateStr))
	seed := int64(binary.BigEndian.Uint64(h[:8]))
	r := rand.New(rand.NewSource(seed))
	idx := r.Intn(len(songs))

	return songs[idx]
}
