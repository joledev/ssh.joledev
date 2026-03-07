package tui

type Lang string

const (
	LangES Lang = "es"
	LangEN Lang = "en"
)

type Translations struct {
	Welcome         string
	TabHome         string
	TabAbout        string
	TabBlog         string
	TabSong         string
	Name            string
	Role            string
	Website         string
	TechStackTitle  string
	NavHelp         string
	LangToggle      string
	Contact         string
	QuitHelp        string
	SongOfTheDay    string
	SongListenAt    string
	SongAlbum       string
	SongArtist      string
	BlogEmpty       string
	BlogBack        string
	SongTitle       string
	CoverTitle      string
	CoverExplain    string
	NoPosts         string
	ReadMore        string
	BackToList      string
}

var I18n = map[Lang]Translations{
	LangES: {
		Welcome:        "Bienvenido a mi rincon digital",
		TabHome:        "Inicio",
		TabAbout:       "Sobre mi",
		TabBlog:        "Blog",
		TabSong:        "Cancion del dia",
		Name:           "Joel Ernesto Lopez Verdugo",
		Role:           "SaaS Fullstack Developer",
		Website:        "joledev.com",
		TechStackTitle: "Tech Stack",
		NavHelp:        "tab: navegar | enter: seleccionar",
		LangToggle:     "L: cambiar idioma",
		Contact:        "contacto@joledev.com",
		QuitHelp:       "q/esc: salir",
		SongOfTheDay:   "Cancion del Dia",
		SongListenAt:   "Escuchala aqui:",
		SongAlbum:      "Album",
		SongArtist:     "Artista",
		BlogEmpty:      "No hay publicaciones aun...",
		BlogBack:       "esc: volver",
		SongTitle:      "Cancion",
		CoverTitle:     "Shinseiki No Love Song - Asian Kung-Fu Generation",
		CoverExplain: `La cancion trata del amor como el unico elemento verdaderamente
humano en un mundo que sigue indiferente ante la muerte, el terror
y el arrepentimiento. No romantiza el amor -- lo describe como algo
incierto, imperfecto y dificil de nombrar, pero precisamente por
eso es lo que nos separa de ser solo animales biologicos.

Es una "cancion de amor" en el sentido mas filosofico y honesto
posible. Adoro todo lo que significa, la cantidad de pasos y
escenarios que cubre desde la primera estrofa -- el peso del pasado
personal -- hasta el coro final: un adios a la edad de piedra.

Como es que con lagrimas en los ojos seguimos adelante.`,
		NoPosts:    "No hay posts todavia...",
		ReadMore:   "enter: leer",
		BackToList: "esc: volver a la lista",
	},
	LangEN: {
		Welcome:        "Welcome to my digital corner",
		TabHome:        "Home",
		TabAbout:       "About",
		TabBlog:        "Blog",
		TabSong:        "Song of the Day",
		Name:           "Joel Ernesto Lopez Verdugo",
		Role:           "SaaS Fullstack Developer",
		Website:        "joledev.com",
		TechStackTitle: "Tech Stack",
		NavHelp:        "tab: navigate | enter: select",
		LangToggle:     "L: switch language",
		Contact:        "contacto@joledev.com",
		QuitHelp:       "q/esc: quit",
		SongOfTheDay:   "Song of the Day",
		SongListenAt:   "Listen here:",
		SongAlbum:      "Album",
		SongArtist:     "Artist",
		BlogEmpty:      "No posts yet...",
		BlogBack:       "esc: back",
		SongTitle:      "Song",
		CoverTitle:     "Shinseiki No Love Song - Asian Kung-Fu Generation",
		CoverExplain: `The song speaks of love as the only truly human element in a world
that remains indifferent to death, terror, and regret. It doesn't
romanticize love -- it describes it as something uncertain, imperfect,
and hard to name, but precisely because of that, it's what separates
us from being just biological animals.

It's a "love song" in the most philosophical and honest sense possible.
I love everything it stands for, the amount of steps and scenarios it
covers from the first verse -- the weight of one's personal past --
to the final chorus: a farewell to the stone age.

How is it that with tears in our eyes, we keep moving forward.`,
		NoPosts:    "No posts yet...",
		ReadMore:   "enter: read",
		BackToList: "esc: back to list",
	},
}

func T(lang Lang) Translations {
	if t, ok := I18n[lang]; ok {
		return t
	}
	return I18n[LangES]
}
