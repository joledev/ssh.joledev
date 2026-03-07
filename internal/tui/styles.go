package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors inspired by the Shinseiki No Love Song cover
	Magenta    = lipgloss.Color("#FF00FF")
	Cyan       = lipgloss.Color("#00D4FF")
	Pink       = lipgloss.Color("#FF69B4")
	White      = lipgloss.Color("#FFFFFF")
	Gray       = lipgloss.Color("#888888")
	DarkGray   = lipgloss.Color("#444444")
	Yellow     = lipgloss.Color("#FFD700")
	Blue       = lipgloss.Color("#4A90D9")

	// Tab styles
	ActiveTab = lipgloss.NewStyle().
			Bold(true).
			Foreground(White).
			Background(Magenta).
			Padding(0, 2)

	InactiveTab = lipgloss.NewStyle().
			Foreground(Gray).
			Padding(0, 2)

	TabBar = lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(DarkGray).
		MarginBottom(1)

	// Content styles
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Magenta).
		MarginBottom(1)

	Subtitle = lipgloss.NewStyle().
			Foreground(Pink).
			Italic(true)

	AsciiStyle = lipgloss.NewStyle().
			Foreground(Magenta)

	TextStyle = lipgloss.NewStyle().
			Foreground(White)

	DimStyle = lipgloss.NewStyle().
			Foreground(Gray)

	AccentStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	LinkStyle = lipgloss.NewStyle().
			Foreground(Blue).
			Underline(true)

	// Tech stack badge
	Badge = lipgloss.NewStyle().
		Foreground(White).
		Background(DarkGray).
		Padding(0, 1)

	// Song card
	SongCard = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Magenta).
			Padding(1, 2).
			MarginTop(1)

	// Blog post
	PostTitle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	PostDate = lipgloss.NewStyle().
			Foreground(Gray)

	PostBody = lipgloss.NewStyle().
			Foreground(White).
			MarginTop(1)

	// Help bar
	HelpStyle = lipgloss.NewStyle().
			Foreground(DarkGray).
			MarginTop(1)

	// Cursor for lists
	CursorStyle = lipgloss.NewStyle().
			Foreground(Magenta).
			Bold(true)

	// Container
	Container = lipgloss.NewStyle().
			Padding(1, 2)
)
