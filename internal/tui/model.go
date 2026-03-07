package tui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joledev/ssh.joledev/internal/art"
	"github.com/joledev/ssh.joledev/internal/data"
)

type Section int

const (
	SectionSong Section = iota
	SectionAbout
	SectionBlog
)

const coverArtWidth = 55

type tickMsg time.Time

type coverMsg struct {
	mono  string
	color string
	err   error
}

type Model struct {
	Lang         Lang
	Section      Section
	Width        int
	Height       int
	Songs        []data.Song
	Posts        []data.Post
	PostsDir     string
	BlogCursor   int
	ReadingPost  bool
	Quitting     bool
	Frame        int
	TodaySong    data.Song
	CoverMono    string
	CoverColor   string
	ColorMode    bool
	CoverLoading bool
	QRCode       string
}

func NewModel(songsPath, postsDir string) Model {
	songs, _ := data.LoadSongs(songsPath)
	posts, _ := data.LoadPosts(postsDir, "es")
	todaySong := data.SongOfTheDay(songs)

	qrStr, _ := art.GenerateQR(todaySong.URL)

	return Model{
		Lang:         LangES,
		Songs:        songs,
		Posts:        posts,
		PostsDir:     postsDir,
		TodaySong:    todaySong,
		CoverLoading: true,
		QRCode:       qrStr,
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*400, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchCoverCmd(trackURL string) tea.Cmd {
	return func() tea.Msg {
		img, err := art.FetchCoverImage(trackURL)
		if err != nil {
			return coverMsg{err: err}
		}
		mono := art.ImageToBraille(img, coverArtWidth, false)
		color := art.ImageToBraille(img, coverArtWidth, true)
		return coverMsg{mono: mono, color: color}
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), fetchCoverCmd(m.TodaySong.URL))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.Frame++
		return m, tickCmd()

	case coverMsg:
		m.CoverLoading = false
		if msg.err == nil {
			m.CoverMono = msg.mono
			m.CoverColor = msg.color
		}
		return m, nil

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
					m.Section = SectionBlog
				}
				m.BlogCursor = 0
				m.ReadingPost = false
			}
			return m, nil

		case "right":
			if !m.ReadingPost {
				if m.Section < SectionBlog {
					m.Section++
				} else {
					m.Section = SectionSong
				}
				m.BlogCursor = 0
				m.ReadingPost = false
			}
			return m, nil

		case "c", "C":
			m.ColorMode = !m.ColorMode
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
	case SectionSong:
		view = m.viewSong(t)
	case SectionAbout:
		view = m.viewAbout(t)
	case SectionBlog:
		view = m.viewBlog(t)
	}

	// Show animated logo on all views except when reading a post
	showLogo := !(m.Section == SectionBlog && m.ReadingPost)

	var header string
	if showLogo {
		header = m.renderLogo()
	}

	nav := m.renderNav(t)
	help := HelpStyle.Render("[<- -> nav · C color · L lang · q quit]")

	if header != "" {
		return fmt.Sprintf("%s\n%s\n\n %s\n\n %s\n", header, view, nav, help)
	}
	return fmt.Sprintf("%s\n\n %s\n\n %s\n", view, nav, help)
}

func (m Model) renderLogo() string {
	logo := []string{
		" ____ ____ ____ ____ ____ ____ ____ ",
		"||J |||o |||l |||e |||D |||e |||v ||",
		"||__|||__|||__|||__|||__|||__|||__||",
		"|/__\\|/__\\|/__\\|/__\\|/__\\|/__\\|/__\\|",
	}

	sparkles := []string{"·", "+", "*", "✦", "·", "+"}
	sparkleColors := []lipgloss.Color{Magenta, Cyan, Pink, Magenta, Cyan, Pink}

	s1 := lipgloss.NewStyle().Foreground(sparkleColors[m.Frame%len(sparkleColors)]).Render(sparkles[m.Frame%len(sparkles)])
	s2 := lipgloss.NewStyle().Foreground(sparkleColors[(m.Frame+2)%len(sparkleColors)]).Render(sparkles[(m.Frame+1)%len(sparkles)])
	s3 := lipgloss.NewStyle().Foreground(sparkleColors[(m.Frame+4)%len(sparkleColors)]).Render(sparkles[(m.Frame+3)%len(sparkles)])

	var lines []string
	lines = append(lines, "  "+s1)
	for i, line := range logo {
		prefix := "  "
		if i == 0 {
			prefix = s2 + " "
		}
		if i == len(logo)-1 {
			line = line + " " + s3
		}
		lines = append(lines, NameStyle.Render(prefix+line))
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderNav(t Translations) string {
	sections := []struct {
		name    string
		section Section
	}{
		{t.TabSong, SectionSong},
		{t.TabAbout, SectionAbout},
		{t.TabBlog, SectionBlog},
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

func (m Model) viewSong(t Translations) string {
	song := m.TodaySong

	rightLines := []string{
		"",
		SubtitleStyle.Render("  "+t.Role),
		"",
		"  " + AccentStyle.Render(t.SongOfTheDay),
		"",
		fmt.Sprintf("  %s  %s",
			AccentStyle.Render(t.SongTitle+":"),
			lipgloss.NewStyle().Foreground(White).Render(song.Title)),
		fmt.Sprintf("  %s  %s",
			AccentStyle.Render(t.SongArtist+":"),
			lipgloss.NewStyle().Foreground(White).Render(song.Artist)),
		fmt.Sprintf("  %s  %s",
			AccentStyle.Render(t.SongAlbum+":"),
			lipgloss.NewStyle().Foreground(White).Render(song.Album)),
		"",
	}

	if m.QRCode != "" {
		rightLines = append(rightLines,
			"  "+DimStyle.Render(t.SongListenAt),
			"",
		)
		for _, qrLine := range strings.Split(m.QRCode, "\n") {
			rightLines = append(rightLines, "  "+qrLine)
		}
		rightLines = append(rightLines, "")
	}

	rightLines = append(rightLines,
		"  "+DimStyle.Render(t.Contact),
		"  "+DimStyle.Render(t.Website),
	)

	if m.CoverLoading {
		loadingLines := []string{"", "", DimStyle.Render("  Loading cover art...")}
		for len(loadingLines) < 27 {
			loadingLines = append(loadingLines, "")
		}
		return buildSideBySide(loadingLines, rightLines, coverArtWidth+2)
	}

	coverArt := m.CoverMono
	if m.ColorMode {
		coverArt = m.CoverColor
	}

	if coverArt == "" {
		var s strings.Builder
		s.WriteString("\n")
		for _, line := range rightLines {
			s.WriteString(line + "\n")
		}
		return s.String()
	}

	artLines := strings.Split(coverArt, "\n")
	framedLines := m.addAnimatedFrame(artLines)
	frameWidth := 0
	for _, line := range framedLines {
		w := visualWidth(line)
		if w > frameWidth {
			frameWidth = w
		}
	}

	return buildSideBySide(framedLines, rightLines, frameWidth)
}

func (m Model) addAnimatedFrame(artLines []string) []string {
	borderColors := []lipgloss.Color{Magenta, Cyan, Pink}
	c1 := borderColors[m.Frame%len(borderColors)]
	c2 := borderColors[(m.Frame+1)%len(borderColors)]
	c3 := borderColors[(m.Frame+2)%len(borderColors)]

	s1 := lipgloss.NewStyle().Foreground(c1)
	s2 := lipgloss.NewStyle().Foreground(c2)
	s3 := lipgloss.NewStyle().Foreground(c3)

	artWidth := 0
	for _, line := range artLines {
		w := visualWidth(line)
		if w > artWidth {
			artWidth = w
		}
	}

	hBar := strings.Repeat("─", artWidth)
	top := s1.Render("╭") + s2.Render(hBar) + s3.Render("╮")
	bottom := s3.Render("╰") + s2.Render(hBar) + s1.Render("╯")

	var framed []string
	framed = append(framed, top)
	for _, line := range artLines {
		pad := artWidth - visualWidth(line)
		framed = append(framed, s1.Render("│")+line+strings.Repeat(" ", pad)+s3.Render("│"))
	}
	framed = append(framed, bottom)
	return framed
}

func (m Model) viewAbout(t Translations) string {
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

	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(" " + NameStyle.Render(t.Name) + "\n")
	s.WriteString(" " + SubtitleStyle.Render(t.Role) + "\n\n")
	s.WriteString(" " + DimStyle.Render(t.Contact) + "\n")
	s.WriteString(" " + LinkStyle.Render(t.Website) + "\n\n")

	for _, line := range strings.Split(t.AboutMe, "\n") {
		s.WriteString(" " + BodyStyle.Render(line) + "\n")
	}

	s.WriteString("\n " + AccentStyle.Render(t.TechStackTitle) + "\n\n")
	s.WriteString(" " + strings.Join(badges[:4], " ") + "\n")
	s.WriteString(" " + strings.Join(badges[4:], " ") + "\n\n")

	for _, line := range strings.Split(t.AboutProject, "\n") {
		s.WriteString(" " + DimStyle.Render(line) + "\n")
	}

	return s.String()
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
