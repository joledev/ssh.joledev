package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors inspired by the Shinseiki No Love Song cover
	Magenta  = lipgloss.Color("#FF00FF")
	Cyan     = lipgloss.Color("#00D4FF")
	Pink     = lipgloss.Color("#FF69B4")
	White    = lipgloss.Color("#FFFFFF")
	Gray     = lipgloss.Color("#666666")
	DarkGray = lipgloss.Color("#333333")
	Blue     = lipgloss.Color("#4A90D9")
	Dim      = lipgloss.Color("#555555")

	// Art style
	ArtStyle = lipgloss.NewStyle().Foreground(Magenta)

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Cyan)

	NameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Magenta)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(White).
			Italic(true)

	BodyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA"))

	DimStyle = lipgloss.NewStyle().
			Foreground(Dim)

	LinkStyle = lipgloss.NewStyle().
			Foreground(Blue)

	AccentStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	// Navigation
	ActiveNavStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Cyan)

	InactiveNavStyle = lipgloss.NewStyle().
				Foreground(Gray)

	HelpStyle = lipgloss.NewStyle().
			Foreground(Dim)

	// Song card
	SongCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Magenta).
			Padding(1, 2)

	// Blog
	CursorStyle = lipgloss.NewStyle().
			Foreground(Magenta).
			Bold(true)

	PostTitleStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	PostDateStyle = lipgloss.NewStyle().
			Foreground(Dim)

	PostBodyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC"))

	// Badge
	BadgeStyle = func(bg string) lipgloss.Style {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color(bg)).
			Padding(0, 1)
	}
)
