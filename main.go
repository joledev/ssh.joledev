package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/joledev/ssh.joledev/internal/tui"
)

func main() {
	host := envOr("SSH_HOST", "0.0.0.0")
	port := envOr("SSH_PORT", "2222")
	keyPath := envOr("SSH_KEY_PATH", ".ssh")
	songsPath := envOr("SONGS_PATH", "")
	postsDir := envOr("POSTS_DIR", "posts")

	// Make paths absolute relative to binary location
	if !filepath.IsAbs(postsDir) {
		if exe, err := os.Executable(); err == nil {
			postsDir = filepath.Join(filepath.Dir(exe), postsDir)
		}
	}
	if songsPath != "" && !filepath.IsAbs(songsPath) {
		if exe, err := os.Executable(); err == nil {
			songsPath = filepath.Join(filepath.Dir(exe), songsPath)
		}
	}

	// Also check relative to CWD as fallback
	if _, err := os.Stat(postsDir); os.IsNotExist(err) {
		if cwd, err := os.Getwd(); err == nil {
			postsDir = filepath.Join(cwd, "posts")
		}
	}
	if songsPath == "" {
		if cwd, err := os.Getwd(); err == nil {
			candidate := filepath.Join(cwd, "data", "songs.txt")
			if _, err := os.Stat(candidate); err == nil {
				songsPath = candidate
			}
		}
	}

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%s", host, port)),
		wish.WithHostKeyPath(keyPath),
		wish.WithMiddleware(
			bubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				m := tui.NewModel(songsPath, postsDir)
				return m, []tea.ProgramOption{tea.WithAltScreen()}
			}),
		),
	)
	if err != nil {
		log.Fatalf("Could not create SSH server: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting SSH server on %s:%s", host, port)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("SSH server error: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
