package tui

type Lang string

const (
	LangES Lang = "es"
	LangEN Lang = "en"
)

type Translations struct {
	TabAbout       string
	TabBlog        string
	TabSong        string
	Name           string
	Role           string
	Website        string
	TechStackTitle string
	Contact        string
	SongOfTheDay   string
	SongListenAt   string
	SongAlbum      string
	SongArtist     string
	SongTitle      string
	NoPosts        string
	ReadMore       string
	BackToList     string
}

var I18n = map[Lang]Translations{
	LangES: {
		TabAbout:       "Sobre mi",
		TabBlog:        "Blog",
		TabSong:        "Cancion del dia",
		Name:           "Joel Ernesto Lopez Verdugo",
		Role:           "SaaS Fullstack Developer",
		Website:        "joledev.com",
		TechStackTitle: "Tech Stack",
		Contact:        "contacto@joledev.com",
		SongOfTheDay:   "Cancion del Dia",
		SongListenAt:   "Escuchala aqui:",
		SongAlbum:      "Album",
		SongArtist:     "Artista",
		SongTitle:      "Cancion",
		NoPosts:        "No hay posts todavia...",
		ReadMore:       "enter: leer",
		BackToList:     "esc: volver a la lista",
	},
	LangEN: {
		TabAbout:       "About",
		TabBlog:        "Blog",
		TabSong:        "Song of the Day",
		Name:           "Joel Ernesto Lopez Verdugo",
		Role:           "SaaS Fullstack Developer",
		Website:        "joledev.com",
		TechStackTitle: "Tech Stack",
		Contact:        "contacto@joledev.com",
		SongOfTheDay:   "Song of the Day",
		SongListenAt:   "Listen here:",
		SongAlbum:      "Album",
		SongArtist:     "Artist",
		SongTitle:      "Song",
		NoPosts:        "No posts yet...",
		ReadMore:       "enter: read",
		BackToList:     "esc: back to list",
	},
}

func T(lang Lang) Translations {
	if t, ok := I18n[lang]; ok {
		return t
	}
	return I18n[LangES]
}
