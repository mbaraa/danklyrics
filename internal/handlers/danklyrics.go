package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mbaraa/danklyrics/internal/actions"
	"github.com/mbaraa/danklyrics/pkg/models"
)

type dankLyricsApi struct {
	usecases *actions.Actions
}

func NewDankLyricsApi(usecases *actions.Actions) *dankLyricsApi {
	return &dankLyricsApi{
		usecases: usecases,
	}
}

func (d *dankLyricsApi) HandleGetSongLyrics(w http.ResponseWriter, r *http.Request) {
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
		lyricses, err = d.usecases.GetLyricsBySongTitle(songName[0])
		if err != nil {
			break
		}
	case okArtist && !okAlbum && okSong:
		lyricses, err = d.usecases.GetLyricsBySongTitleAndArtistName(songName[0], artistName[0])
		if err != nil {
			break
		}
	case !okArtist && okAlbum && okSong:
		lyricses, err = d.usecases.GetLyricsBySongTitleAndArtistName(songName[0], albumName[0])
		if err != nil {
			break
		}
	case okArtist && okAlbum && okSong:
		lyricses, err = d.usecases.GetLyricsBySongTitleArtistNameAndAlbumTitle(songName[0], artistName[0], albumName[0])
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
