package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/mbaraa/danklyrics/internal/actions"
	"github.com/mbaraa/danklyrics/pkg/client"
	"github.com/mbaraa/danklyrics/pkg/models"
	"github.com/mbaraa/danklyrics/pkg/provider"
)

type lyricsFinderApi struct {
	usecases *actions.Actions
}

func NewLyricsFinderApi(usecases *actions.Actions) *lyricsFinderApi {
	return &lyricsFinderApi{
		usecases: usecases,
	}
}

func (l *lyricsFinderApi) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("refer to (" + docsLink + ") for API docs!"))
}

func (l *lyricsFinderApi) HandleListProviders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode([]map[string]string{
		{"name": "LyricFind", "id": "lrc"},
		{"name": "Genius", "id": "genius"},
	})
}

func (l *lyricsFinderApi) HandleGetSongLyrics(w http.ResponseWriter, r *http.Request) {
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
		lyricses, _ = l.usecases.GetLyricsBySongTitle(songName[0])
		if err != nil {
			break
		}
	case okArtist && !okAlbum && okSong:
		lyricses, _ = l.usecases.GetLyricsBySongTitleAndArtistName(songName[0], artistName[0])
		if err != nil {
			break
		}
	case !okArtist && okAlbum && okSong:
		lyricses, _ = l.usecases.GetLyricsBySongTitleAndArtistName(songName[0], albumName[0])
		if err != nil {
			break
		}
	case okArtist && okAlbum && okSong:
		lyricses, _ = l.usecases.GetLyricsBySongTitleArtistNameAndAlbumTitle(songName[0], artistName[0], albumName[0])
		if err != nil {
			break
		}
	}
	if len(lyricses) == 0 && len(lyrics.Parts) > 0 {
		_, _ = l.usecases.CreateLyrics(lyrics)
	}

	_ = json.NewEncoder(w).Encode(lyrics)
}

func (l *lyricsFinderApi) HandleSubmitSongLyrics(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Header["Authorization"]
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "No no no, you need to authenticate first!",
		})
		return
	}

	var lyrics struct {
		SongName   string `json:"song_name"`
		ArtistName string `json:"artist_name"`
		AlbumName  string `json:"album_name"`
		Plain      string `json:"plain_lyrics"`
	}
	err := json.NewDecoder(r.Body).Decode(&lyrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "invalid request body",
		})
		return
	}

	if lyrics.SongName == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "missing required field `song_name`",
		})
		return
	}
	if lyrics.ArtistName == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "missing required field `artist_name`",
		})
		return
	}
	if lyrics.AlbumName == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "missing required field `album_name`",
		})
		return
	}
	if lyrics.Plain == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "missing required field `plain_lyrics`",
		})
		return
	}

	err = l.usecases.CreateLyricsRequest(token[0], models.Lyrics{
		SongName:   lyrics.SongName,
		ArtistName: lyrics.ArtistName,
		AlbumName:  lyrics.AlbumName,
		Parts:      strings.Split(lyrics.Plain, "\n"),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "Something went wrong",
		})
		return
	}
}
