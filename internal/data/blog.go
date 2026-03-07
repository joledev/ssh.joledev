package data

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Post struct {
	Slug    string
	Title   string
	Date    time.Time
	Content string
	Lang    string
}

func LoadPosts(postsDir string, lang string) ([]Post, error) {
	langDir := filepath.Join(postsDir, lang)
	entries, err := os.ReadDir(langDir)
	if err != nil {
		// Fallback: try root posts dir
		entries, err = os.ReadDir(postsDir)
		if err != nil {
			return nil, err
		}
		langDir = postsDir
	}

	var posts []Post
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		content, err := os.ReadFile(filepath.Join(langDir, e.Name()))
		if err != nil {
			continue
		}

		slug := strings.TrimSuffix(e.Name(), ".md")
		post := Post{
			Slug:    slug,
			Lang:    lang,
			Content: string(content),
		}

		// Parse frontmatter-like first lines: "# Title" and "date: YYYY-MM-DD"
		lines := strings.SplitN(string(content), "\n", 10)
		body_start := 0
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "# ") {
				post.Title = strings.TrimPrefix(line, "# ")
				body_start = i + 1
			} else if strings.HasPrefix(line, "date:") {
				dateStr := strings.TrimSpace(strings.TrimPrefix(line, "date:"))
				if t, err := time.Parse("2006-01-02", dateStr); err == nil {
					post.Date = t
				}
				body_start = i + 1
			} else if line == "" {
				continue
			} else {
				break
			}
		}

		if post.Title == "" {
			post.Title = slug
		}
		if post.Date.IsZero() {
			info, _ := e.Info()
			if info != nil {
				post.Date = info.ModTime()
			}
		}

		// Set body to everything after frontmatter
		if body_start > 0 && body_start < len(lines) {
			post.Content = strings.TrimSpace(strings.Join(lines[body_start:], "\n"))
		}

		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, nil
}
