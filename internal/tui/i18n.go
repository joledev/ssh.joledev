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
	AboutMe        string
	AboutProject   string
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
		AboutMe: `Desarrollo SaaS para empresas pequenas y medianas.
Me interesa acercar la tecnologia a las personas de
forma limpia y organica, sin invadir su vida ni crear
necesidades falsas. Priorizo la interaccion real de
los usuarios y resolver problemas que existen.

Me gusta que las cosas se vean bien y funcionen bien.
Trabajo en un entorno minimalista: codigo limpio,
ordenado, que no crezca al infinito y me permita
crecer con el sin caer en la sobreingenieria.`,
		AboutProject: `Este proyecto nacio de un reel de Instagram de
@morilliu y decidi darle un enfoque mas personal:
recomendar canciones, compartir un poco de lo que
me gusta -- la musica, el arte, la filosofia y la
comida. Dejar un poco de mi en algo minimalista.`,
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
		AboutMe: `I build SaaS for small and medium businesses.
I care about bringing technology to people in a
clean, organic way -- without invading their lives
or creating false needs. I prioritize real user
interaction and solving problems that actually exist.

I like things that look good and work well.
I work in a minimalist environment: clean code,
organized, that doesn't grow to infinity and lets
me grow alongside it without over-engineering.`,
		AboutProject: `This project was born from an Instagram reel by
@morilliu and I decided to give it a more personal
spin: recommending songs, sharing a bit of what I
enjoy -- music, art, philosophy, and food.
Leaving a piece of myself in something minimalist.`,
	},
}

func T(lang Lang) Translations {
	if t, ok := I18n[lang]; ok {
		return t
	}
	return I18n[LangES]
}
