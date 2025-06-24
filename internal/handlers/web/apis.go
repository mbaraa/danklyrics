package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/mbaraa/danklyrics/internal/actions"
	"github.com/mbaraa/danklyrics/internal/config"
	"github.com/mbaraa/danklyrics/pkg/client"
	"github.com/mbaraa/danklyrics/pkg/provider"
)

type api struct {
	usecases *actions.Actions
	lyricser *client.Http
}

func NewApi(usecases *actions.Actions) *api {
	lyricser, err := client.NewHttp(client.Config{
		Providers:  []provider.Name{provider.Dank, provider.LyricFind},
		ApiAddress: config.Env().ApiAddress,
	})
	if err != nil {
		panic(err)
	}

	return &api{
		usecases: usecases,
		lyricser: lyricser,
	}
}

func (a *api) HandleGetSongLyrics(w http.ResponseWriter, r *http.Request) {
	artistName, okArtist := r.URL.Query()["artist"]
	albumName, okAlbum := r.URL.Query()["album"]
	songName, okSong := r.URL.Query()["song"]

	if !okArtist && !okAlbum && !okSong {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing all query parameters `artist`, `album` and `song`"))
		return
	}

	if !okSong {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing required query parameter `song`"))
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

	lyricsText, err := a.lyricser.GetSongLyrics(searchInput)
	if err != nil {
		log.Println("oppsie doopsie some shit happened", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("No results were found"))
		return
	}

	w.Write([]byte(lyricsText.String()))
}

func (a *api) HandleAuthSubmitLyrics(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid body"))
		return
	}

	err = makeApiPostRequest("/auth", "", reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
}

func (a *api) HandleConfirmAuthSubmitLyrics(w http.ResponseWriter, r *http.Request) {
	token, tokenExists := r.URL.Query()["token"]

	if !tokenExists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing token"))
		return
	}

	err := makeApiPostRequest("/auth/confirm", "", map[string]string{
		"token": token[0],
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token[0],
		Path:    "/",
		Expires: time.Now().UTC().Add(time.Hour),
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (a *api) HandleSubmitLyrics(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("No no no, you need to authenticate first!"))
		return
	}

	var lyrics struct {
		SongName   string `json:"song_name"`
		ArtistName string `json:"artist_name"`
		AlbumName  string `json:"album_name"`
		Plain      string `json:"plain_lyrics"`
	}
	err = json.NewDecoder(r.Body).Decode(&lyrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request body"))
		return
	}

	if lyrics.SongName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing required field `song_name`"))
		return
	}
	if lyrics.ArtistName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing required field `artist_name`"))
		return
	}
	if lyrics.AlbumName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing required field `album_name`"))
		return
	}
	if lyrics.Plain == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing required field `plain_lyrics`"))
		return
	}

	err = makeApiPostRequest("/dank/lyrics", sessionToken.Value, lyrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
}

func makeApiPostRequest[T any](path, token string, body T) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, config.Env().ApiAddress+path, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("api responded with a non 200 status")
	}

	return nil
}
