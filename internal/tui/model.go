package tui

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joledev/ssh.joledev/internal/data"
)

//go:embed ascii_data.txt
var asciiArt string

type Tab int

const (
	TabHome Tab = iota
	TabAbout
	TabBlog
	TabSong
)

var tabOrder = []Tab{TabHome, TabAbout, TabBlog, TabSong}

type Model struct {
	Lang       Lang
	ActiveTab  Tab
	Width      int
	Height     int
	Songs      []data.Song
	Posts      []data.Post
	PostsDir   string
	BlogCursor int
	ReadingPost bool
	Quitting   bool
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

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			if m.ReadingPost {
				m.ReadingPost = false
				return m, nil
			}
			m.Quitting = true
			return m, tea.Quit

		case "tab", "right":
			if !m.ReadingPost {
				idx := int(m.ActiveTab)
				idx = (idx + 1) % len(tabOrder)
				m.ActiveTab = tabOrder[idx]
				m.BlogCursor = 0
				m.ReadingPost = false
			}
			return m, nil

		case "shift+tab", "left":
			if !m.ReadingPost {
				idx := int(m.ActiveTab)
				idx = (idx - 1 + len(tabOrder)) % len(tabOrder)
				m.ActiveTab = tabOrder[idx]
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
			// Reload posts for new language
			posts, _ := data.LoadPosts(m.PostsDir, string(m.Lang))
			m.Posts = posts
			return m, nil

		case "up", "k":
			if m.ActiveTab == TabBlog && !m.ReadingPost && m.BlogCursor > 0 {
				m.BlogCursor--
			}
			return m, nil

		case "down", "j":
			if m.ActiveTab == TabBlog && !m.ReadingPost && m.BlogCursor < len(m.Posts)-1 {
				m.BlogCursor++
			}
			return m, nil

		case "enter":
			if m.ActiveTab == TabBlog && !m.ReadingPost && len(m.Posts) > 0 {
				m.ReadingPost = true
			}
			return m, nil

		case "1":
			m.ActiveTab = TabHome
			return m, nil
		case "2":
			m.ActiveTab = TabAbout
			return m, nil
		case "3":
			m.ActiveTab = TabBlog
			return m, nil
		case "4":
			m.ActiveTab = TabSong
			return m, nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		t := T(m.Lang)
		if m.Lang == LangES {
			return DimStyle.Render("Hasta luego. -- " + t.Name) + "\n"
		}
		return DimStyle.Render("See you around. -- " + t.Name) + "\n"
	}

	t := T(m.Lang)

	var s strings.Builder

	// Tab bar
	s.WriteString(m.renderTabs(t))
	s.WriteString("\n\n")

	// Content based on active tab
	contentHeight := m.Height - 6 // Reserve space for tabs and help
	_ = contentHeight

	switch m.ActiveTab {
	case TabHome:
		s.WriteString(m.renderHome(t))
	case TabAbout:
		s.WriteString(m.renderAbout(t))
	case TabBlog:
		s.WriteString(m.renderBlog(t))
	case TabSong:
		s.WriteString(m.renderSong(t))
	}

	// Help bar
	s.WriteString("\n")
	help := fmt.Sprintf(" %s | %s | %s | 1-4: tabs",
		t.NavHelp, t.LangToggle, t.QuitHelp)
	s.WriteString(HelpStyle.Render(help))

	return Container.Render(s.String())
}

func (m Model) renderTabs(t Translations) string {
	tabs := []string{t.TabHome, t.TabAbout, t.TabBlog, t.TabSong}
	var rendered []string
	for i, tab := range tabs {
		if Tab(i) == m.ActiveTab {
			rendered = append(rendered, ActiveTab.Render(tab))
		} else {
			rendered = append(rendered, InactiveTab.Render(tab))
		}
	}
	return TabBar.Render(strings.Join(rendered, " "))
}

func (m Model) renderHome(t Translations) string {
	var s strings.Builder

	// ASCII art
	s.WriteString(AsciiStyle.Render(asciiArt))
	s.WriteString("\n\n")

	// Title
	s.WriteString(Title.Render(t.CoverTitle))
	s.WriteString("\n\n")

	// Song explanation
	s.WriteString(TextStyle.Render(t.CoverExplain))
	s.WriteString("\n")

	return s.String()
}

func (m Model) renderAbout(t Translations) string {
	var s strings.Builder

	s.WriteString(Title.Render(t.Name))
	s.WriteString("\n")
	s.WriteString(Subtitle.Render(t.Role))
	s.WriteString("\n\n")

	s.WriteString(LinkStyle.Render(t.Website))
	s.WriteString("  ")
	s.WriteString(DimStyle.Render(t.Contact))
	s.WriteString("\n\n")

	s.WriteString(AccentStyle.Render(t.TechStackTitle))
	s.WriteString("\n\n")

	techStack := []struct{ name, color string }{
		{"PHP", "#777BB4"},
		{"Laravel", "#FF2D20"},
		{"Java", "#ED8B00"},
		{"Go", "#00ADD8"},
		{"TypeScript", "#3178C6"},
		{"Flutter", "#02569B"},
		{"React Native", "#61DAFB"},
		{"React", "#61DAFB"},
	}

	var badges []string
	for _, tech := range techStack {
		badge := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color(tech.color)).
			Padding(0, 1).
			Render(tech.name)
		badges = append(badges, badge)
	}

	// Render badges in rows
	line := ""
	for i, b := range badges {
		if i > 0 {
			line += " "
		}
		line += b
	}
	s.WriteString(line)
	s.WriteString("\n")

	return s.String()
}

func (m Model) renderBlog(t Translations) string {
	var s strings.Builder

	s.WriteString(Title.Render(t.TabBlog))
	s.WriteString("\n\n")

	if len(m.Posts) == 0 {
		s.WriteString(DimStyle.Render(t.NoPosts))
		return s.String()
	}

	if m.ReadingPost && m.BlogCursor < len(m.Posts) {
		post := m.Posts[m.BlogCursor]
		s.WriteString(PostTitle.Render(post.Title))
		s.WriteString("\n")
		s.WriteString(PostDate.Render(post.Date.Format("2006-01-02")))
		s.WriteString("\n")
		s.WriteString(PostBody.Render(post.Content))
		s.WriteString("\n\n")
		s.WriteString(DimStyle.Render(t.BackToList))
		return s.String()
	}

	for i, post := range m.Posts {
		cursor := "  "
		titleStyle := PostTitle
		if i == m.BlogCursor {
			cursor = CursorStyle.Render("> ")
			titleStyle = titleStyle.Foreground(Magenta)
		}
		date := PostDate.Render(post.Date.Format("2006-01-02"))
		s.WriteString(fmt.Sprintf("%s%s  %s\n", cursor, titleStyle.Render(post.Title), date))
	}

	s.WriteString("\n")
	s.WriteString(DimStyle.Render(t.ReadMore))
	return s.String()
}

func (m Model) renderSong(t Translations) string {
	var s strings.Builder

	song := data.SongOfTheDay(m.Songs)

	s.WriteString(Title.Render(t.SongOfTheDay))
	s.WriteString("\n\n")

	card := fmt.Sprintf(
		"%s  %s\n%s  %s\n%s  %s",
		AccentStyle.Render(t.SongTitle+":"),
		TextStyle.Render(song.Title),
		AccentStyle.Render(t.SongArtist+":"),
		TextStyle.Render(song.Artist),
		AccentStyle.Render(t.SongAlbum+":"),
		TextStyle.Render(song.Album),
	)

	s.WriteString(SongCard.Render(card))
	s.WriteString("\n\n")

	if song.URL != "" {
		s.WriteString(DimStyle.Render(t.SongListenAt))
		s.WriteString("\n")
		s.WriteString(LinkStyle.Render(song.URL))
	}

	return s.String()
}
