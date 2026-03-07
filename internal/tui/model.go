package tui

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joledev/ssh.joledev/internal/data"
)

//go:embed ascii_data.txt
var brailleArt string

type Section int

const (
	SectionHome Section = iota
	SectionAbout
	SectionBlog
	SectionSong
)

type tickMsg time.Time

type Model struct {
	Lang        Lang
	Section     Section
	Width       int
	Height      int
	Songs       []data.Song
	Posts       []data.Post
	PostsDir    string
	BlogCursor  int
	ReadingPost bool
	Quitting    bool
	Frame       int
}

func NewModel(songsPath, postsDir string) Model {
	songs, _ := data.LoadSongs(songsPath)
	posts, _ := data.LoadPosts(postsDir, "es")

	return Model{
		Lang:     LangES,
		Songs:    songs,
		Posts:    posts,
		PostsDir: postsDir,
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*400, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.Frame++
		return m, tickCmd()

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit

		case "esc":
			if m.ReadingPost {
				m.ReadingPost = false
				return m, nil
			}
			m.Quitting = true
			return m, tea.Quit

		case "left":
			if !m.ReadingPost {
				if m.Section > 0 {
					m.Section--
				} else {
					m.Section = SectionSong
				}
				m.BlogCursor = 0
				m.ReadingPost = false
			}
			return m, nil

		case "right":
			if !m.ReadingPost {
				if m.Section < SectionSong {
					m.Section++
				} else {
					m.Section = SectionHome
				}
				m.BlogCursor = 0
				m.ReadingPost = false
			}
			return m, nil

		case "l", "L":
			if m.Lang == LangES {
				m.Lang = LangEN
			} else {
				m.Lang = LangES
			}
			posts, _ := data.LoadPosts(m.PostsDir, string(m.Lang))
			m.Posts = posts
			return m, nil

		case "up", "k":
			if m.Section == SectionBlog && !m.ReadingPost && m.BlogCursor > 0 {
				m.BlogCursor--
			}
			return m, nil

		case "down", "j":
			if m.Section == SectionBlog && !m.ReadingPost && m.BlogCursor < len(m.Posts)-1 {
				m.BlogCursor++
			}
			return m, nil

		case "enter":
			if m.Section == SectionBlog && !m.ReadingPost && len(m.Posts) > 0 {
				m.ReadingPost = true
			}
			return m, nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		t := T(m.Lang)
		if m.Lang == LangES {
			return "\n " + DimStyle.Render("Hasta luego. -- "+t.Name) + "\n\n"
		}
		return "\n " + DimStyle.Render("See you around. -- "+t.Name) + "\n\n"
	}

	t := T(m.Lang)
	var view string

	switch m.Section {
	case SectionHome:
		view = m.viewHome(t)
	case SectionAbout:
		view = m.viewAbout(t)
	case SectionBlog:
		view = m.viewBlog(t)
	case SectionSong:
		view = m.viewSong(t)
	}

	// Navigation bar at bottom
	nav := m.renderNav(t)
	help := HelpStyle.Render("[<- -> to select · enter to open · L lang · q to quit]")

	return fmt.Sprintf("%s\n\n %s\n\n %s\n", view, nav, help)
}

func (m Model) renderNav(t Translations) string {
	sections := []struct {
		name    string
		section Section
	}{
		{t.TabHome, SectionHome},
		{t.TabAbout, SectionAbout},
		{t.TabBlog, SectionBlog},
		{t.TabSong, SectionSong},
	}

	var parts []string
	for _, s := range sections {
		marker := "  "
		style := InactiveNavStyle
		if s.section == m.Section {
			marker = AccentStyle.Render("* ")
			style = ActiveNavStyle
		}
		parts = append(parts, marker+style.Render(s.name))
	}

	return " " + strings.Join(parts, "    ")
}

func (m Model) viewHome(t Translations) string {
	artLines := strings.Split(strings.TrimSpace(brailleArt), "\n")
	artWidth := 0
	for _, line := range artLines {
		w := visualWidth(line)
		if w > artWidth {
			artWidth = w
		}
	}

	// smkeyboard figlet "JoleDev"
	logo := []string{
		" ____ ____ ____ ____ ____ ____ ____ ",
		"||J |||o |||l |||e |||D |||e |||v ||",
		"||__|||__|||__|||__|||__|||__|||__||",
		"|/__\\|/__\\|/__\\|/__\\|/__\\|/__\\|/__\\|",
	}

	// Animated sparkles
	sparkles := []string{"·", "+", "*", "✦", "·", "+"}
	sparkleColors := []lipgloss.Color{Magenta, Cyan, Pink, Magenta, Cyan, Pink}

	s1 := lipgloss.NewStyle().Foreground(sparkleColors[m.Frame%len(sparkleColors)]).Render(sparkles[m.Frame%len(sparkles)])
	s2 := lipgloss.NewStyle().Foreground(sparkleColors[(m.Frame+2)%len(sparkleColors)]).Render(sparkles[(m.Frame+1)%len(sparkles)])
	s3 := lipgloss.NewStyle().Foreground(sparkleColors[(m.Frame+4)%len(sparkleColors)]).Render(sparkles[(m.Frame+3)%len(sparkles)])

	rightLines := []string{
		"",
		"  " + s1,
	}
	for i, line := range logo {
		prefix := "  "
		if i == 0 {
			prefix = s2 + " "
		}
		if i == len(logo)-1 {
			line = line + " " + s3
		}
		rightLines = append(rightLines, NameStyle.Render(prefix+line))
	}
	rightLines = append(rightLines,
		"",
		SubtitleStyle.Render("  "+t.Role),
		"",
		BodyStyle.Render("  "+t.Welcome+"."),
		"",
		DimStyle.Render("  "+t.Contact),
		DimStyle.Render("  "+t.Website),
		"",
		"",
		"",
		"",
		"",
		"",
	)

	return buildSideBySide(artLines, rightLines, artWidth)
}

func (m Model) viewAbout(t Translations) string {
	artLines := strings.Split(strings.TrimSpace(brailleArt), "\n")
	artWidth := 0
	for _, line := range artLines {
		w := visualWidth(line)
		if w > artWidth {
			artWidth = w
		}
	}

	techs := []struct{ name, bg string }{
		{"PHP", "#777BB4"},
		{"Laravel", "#FF2D20"},
		{"Java", "#ED8B00"},
		{"Go", "#00ADD8"},
		{"TypeScript", "#3178C6"},
		{"Flutter", "#02569B"},
		{"React Native", "#61DAFB"},
	}

	var badges []string
	for _, tech := range techs {
		badges = append(badges, BadgeStyle(tech.bg).Render(tech.name))
	}

	rightLines := []string{
		"",
		"",
		"",
		NameStyle.Render("  " + t.Name),
		SubtitleStyle.Render("  " + t.Role),
		"",
		DimStyle.Render("  " + t.Contact),
		LinkStyle.Render("  " + t.Website),
		"",
		"",
		AccentStyle.Render("  " + t.TechStackTitle),
		"",
		"  " + strings.Join(badges[:4], " "),
		"  " + strings.Join(badges[4:], " "),
		"",
		"",
		BodyStyle.Render("  " + t.CoverTitle),
		"",
	}

	explainLines := strings.Split(t.CoverExplain, "\n")
	for _, line := range explainLines {
		rightLines = append(rightLines, DimStyle.Render("  "+line))
	}

	return buildSideBySide(artLines, rightLines, artWidth)
}

func (m Model) viewBlog(t Translations) string {
	var s strings.Builder
	s.WriteString("\n")

	if m.ReadingPost && m.BlogCursor < len(m.Posts) {
		post := m.Posts[m.BlogCursor]
		s.WriteString(" " + PostTitleStyle.Render(post.Title) + "\n")
		s.WriteString(" " + PostDateStyle.Render(post.Date.Format("2006-01-02")) + "\n\n")

		for _, line := range strings.Split(post.Content, "\n") {
			s.WriteString(" " + PostBodyStyle.Render(line) + "\n")
		}
		s.WriteString("\n " + DimStyle.Render("[esc "+t.BackToList+"]"))
		return s.String()
	}

	s.WriteString(" " + TitleStyle.Render(t.TabBlog) + "\n\n")

	if len(m.Posts) == 0 {
		s.WriteString(" " + DimStyle.Render(t.NoPosts))
		return s.String()
	}

	for i, post := range m.Posts {
		cursor := "  "
		style := InactiveNavStyle
		if i == m.BlogCursor {
			cursor = CursorStyle.Render("> ")
			style = ActiveNavStyle
		}
		date := PostDateStyle.Render(post.Date.Format("2006-01-02"))
		s.WriteString(fmt.Sprintf(" %s%s  %s\n", cursor, style.Render(post.Title), date))
	}

	s.WriteString("\n " + DimStyle.Render("[enter "+t.ReadMore+"]"))
	return s.String()
}

func (m Model) viewSong(t Translations) string {
	song := data.SongOfTheDay(m.Songs)

	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(" " + TitleStyle.Render(t.SongOfTheDay) + "\n\n")

	card := fmt.Sprintf(
		" %s  %s\n %s  %s\n %s  %s",
		AccentStyle.Render(t.SongTitle+":"),
		lipgloss.NewStyle().Foreground(White).Render(song.Title),
		AccentStyle.Render(t.SongArtist+":"),
		lipgloss.NewStyle().Foreground(White).Render(song.Artist),
		AccentStyle.Render(t.SongAlbum+":"),
		lipgloss.NewStyle().Foreground(White).Render(song.Album),
	)

	s.WriteString(SongCardStyle.Render(card))
	s.WriteString("\n\n")

	if song.URL != "" {
		s.WriteString(" " + DimStyle.Render(t.SongListenAt) + "\n")
		s.WriteString(" " + LinkStyle.Render(song.URL) + "\n")
	}

	return s.String()
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripAnsi(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func visualWidth(s string) int {
	return len([]rune(stripAnsi(s)))
}

func buildSideBySide(artLines []string, rightLines []string, artVisualWidth int) string {
	maxLines := len(artLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	var rows []string
	for i := 0; i < maxLines; i++ {
		left := ""
		leftVisual := 0
		if i < len(artLines) {
			left = artLines[i]
			leftVisual = visualWidth(artLines[i])
		}

		padding := ""
		if leftVisual < artVisualWidth {
			padding = strings.Repeat(" ", artVisualWidth-leftVisual)
		}

		right := ""
		if i < len(rightLines) {
			right = rightLines[i]
		}

		rows = append(rows, " "+left+padding+"  "+right)
	}

	return strings.Join(rows, "\n")
}
