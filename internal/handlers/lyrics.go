package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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
		{"name": "DankLyrics", "id": "dank"},
		{"name": "LyricFind", "id": "lrc"},
	})
}

func (l *lyricsFinderApi) HandleGetSongLyrics(w http.ResponseWriter, r *http.Request) {
	providers := r.URL.Query()["providers"]

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
		Providers: providersConfig,
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
