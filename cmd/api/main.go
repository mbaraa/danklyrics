package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/mbaraa/danklyrics/internal/actions"
	"github.com/mbaraa/danklyrics/internal/config"
	"github.com/mbaraa/danklyrics/internal/mariadb"
	"github.com/mbaraa/danklyrics/pkg/client"
	"github.com/mbaraa/danklyrics/pkg/models"
	"github.com/mbaraa/danklyrics/pkg/provider"
)

const docsLink = "https://github.com/mbaraa/danklyrics"

var (
	usecases *actions.Actions
)

func init() {
	repo, err := mariadb.New()
	if err != nil {
		panic(err)
	}

	err = mariadb.Migrate()
	if err != nil {
		panic(err)
	}

	usecases = actions.New(repo)
}

type errorResponse struct {
	Message         string `json:"message"`
	SuggestedAction string `json:"suggested_action,omitempty"`
	DocsLink        string `json:"docs_link,omitempty"`
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("refer to (" + docsLink + ") for API docs!"))
}

func handleListProviders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode([]map[string]string{
		{"name": "LyricFind", "id": "lrc"},
		{"name": "Genius", "id": "genius"},
	})
}

func handleGetSongLyrics(w http.ResponseWriter, r *http.Request) {
	providers := r.URL.Query()["providers"]
	geniusClientId := r.URL.Query().Get("genius_client_id")
	geniusClientSecret := r.URL.Query().Get("genius_client_secret")

	artistName, okArtist := r.URL.Query()["artist"]
	albumName, okAlbum := r.URL.Query()["album"]
	songName, okSong := r.URL.Query()["song"]

	w.Header().Set("Content-Type", "application/json")

	if len(providers) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:         "You must specify at least one provider",
			SuggestedAction: "Check the `GET /providers` endpoint.",
			DocsLink:        docsLink,
		})
		return
	}
	if slices.Contains(providers, "genius") && (geniusClientId == "" || geniusClientSecret == "") {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:         "You must specify genius' client id and secret when using genius as a provider.",
			SuggestedAction: "Visit https://docs.genius.com/",
			DocsLink:        docsLink,
		})
		return
	}

	if !okArtist && !okAlbum && !okSong {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:  "Missing all query parameters `artist`, `album` and `song`",
			DocsLink: docsLink,
		})
		return
	}

	if !okSong {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:  "Missing required query parameter `song`",
			DocsLink: docsLink,
		})
		return
	}

	searchInput := provider.SearchParams{
		SongName: songName[0],
	}
	if okAlbum {
		searchInput.AlbumName = albumName[0]
	}
	if okArtist {
		searchInput.ArtistName = artistName[0]
	}

	providersConfig := make([]provider.Name, 0, len(providers))
	for _, p := range providers {
		providersConfig = append(providersConfig, provider.Name(p))
	}
	lyricser, err := client.New(client.Config{
		GeniusClientId:     geniusClientId,
		GeniusClientSecret: geniusClientSecret,
		Providers:          providersConfig,
	})

	lyrics, err := lyricser.GetSongLyrics(searchInput)
	if err != nil {
		log.Println("oppsie doopsie some shit happened", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "No results were found",
		})
		return
	}

	var lyricses []models.Lyrics
	switch {
	case !okArtist && !okAlbum && okSong:
		lyricses, err = usecases.GetLyricsBySongTitle(songName[0])
		if err != nil {
			break
		}
	case okArtist && !okAlbum && okSong:
		lyricses, err = usecases.GetLyricsBySongTitleAndArtistName(songName[0], artistName[0])
		if err != nil {
			break
		}
	case !okArtist && okAlbum && okSong:
		lyricses, err = usecases.GetLyricsBySongTitleAndArtistName(songName[0], albumName[0])
		if err != nil {
			break
		}
	case okArtist && okAlbum && okSong:
		lyricses, err = usecases.GetLyricsBySongTitleArtistNameAndAlbumTitle(songName[0], artistName[0], albumName[0])
		if err != nil {
			break
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:         "Something went wrong",
			SuggestedAction: "Check the docs, or contact admin (baraa@dankstuff.net)",
			DocsLink:        docsLink,
		})
		return
	}
	if len(lyricses) == 0 && len(lyrics.Parts) > 0 {
		_, _ = usecases.CreateLyrics(lyrics)
	}

	_ = json.NewEncoder(w).Encode(lyrics)
}

func handleGetSongDankLyrics(w http.ResponseWriter, r *http.Request) {
	artistName, okArtist := r.URL.Query()["artist"]
	albumName, okAlbum := r.URL.Query()["album"]
	songName, okSong := r.URL.Query()["song"]

	w.Header().Set("Content-Type", "application/json")

	if !okArtist && !okAlbum && !okSong {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:  "Missing all query parameters `artist`, `album` and `song`",
			DocsLink: docsLink,
		})
		return
	}

	if !okSong {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:  "Missing required query parameter `song`",
			DocsLink: docsLink,
		})
		return
	}

	var lyricses []models.Lyrics
	var err error

	switch {
	case !okArtist && !okAlbum && okSong:
		lyricses, err = usecases.GetLyricsBySongTitle(songName[0])
		if err != nil {
			break
		}
	case okArtist && !okAlbum && okSong:
		lyricses, err = usecases.GetLyricsBySongTitleAndArtistName(songName[0], artistName[0])
		if err != nil {
			break
		}
	case !okArtist && okAlbum && okSong:
		lyricses, err = usecases.GetLyricsBySongTitleAndArtistName(songName[0], albumName[0])
		if err != nil {
			break
		}
	case okArtist && okAlbum && okSong:
		lyricses, err = usecases.GetLyricsBySongTitleArtistNameAndAlbumTitle(songName[0], artistName[0], albumName[0])
		if err != nil {
			break
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:         "Something went wrong",
			SuggestedAction: "Check the docs, or contact admin (baraa@dankstuff.net)",
			DocsLink:        docsLink,
		})
		return
	}

	if len(lyricses) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "No results were found",
		})
		return
	}

	_ = json.NewEncoder(w).Encode(lyricses)
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("GET /providers", handleListProviders)
	http.HandleFunc("GET /lyrics", handleGetSongLyrics)
	http.HandleFunc("GET /dank/lyrics", handleGetSongDankLyrics)

	log.Printf("Starting web server at port %s", config.Env().ApiPort)
	log.Fatalln(http.ListenAndServe(":"+config.Env().ApiPort, nil))
}
