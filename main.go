package main

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mbaraa/gonius"
)

var (
	port        = os.Getenv("PORT")
	geniusToken = os.Getenv("GENIUS_TOKEN")
)

var (
	//go:embed static/*
	publicFiles embed.FS

	genius *gonius.Client
)

func init() {
	genius = gonius.NewClient(geniusToken)
}

type SearchInput struct {
	SongName   string
	ArtistName string
	AlbumName  string
}

func getSongLyrics(s SearchInput) (string, error) {
	var hits []gonius.Hit
	var err error

	okArtist := s.ArtistName != ""
	okAlbum := s.AlbumName != ""
	okSong := s.SongName != ""

	switch {
	case !okArtist && !okAlbum && okSong:
		fmt.Println("just song")
		hits, err = genius.Search.Get(s.SongName)
		if err != nil {
			return "", err
		}
	case okArtist && !okAlbum && okSong:
		fmt.Println("artist & song")
		hits, err = genius.Search.Get(fmt.Sprintf("%s %s", s.SongName, s.ArtistName))
		if err != nil {
			return "", err
		}
	case !okArtist && okAlbum && okSong:
		fmt.Println("album & song")
		hits, err = genius.Search.Get(fmt.Sprintf("%s %s", s.SongName, s.AlbumName))
		if err != nil {
			return "", err
		}
	case okArtist && okAlbum && okSong:
		log.Println("artist, album and song")
		hits, err = genius.Search.Get(fmt.Sprintf("%s %s %s", s.SongName, s.AlbumName, s.ArtistName))
		if err != nil {
			return "", err
		}
	}

	if len(hits) == 0 {
		return "", errors.New("no results were found")
	}

	var songUrl string
	for _, hit := range hits {
		if hit.Type == "song" && hit.Result != nil {
			songUrl = hit.Result.URL
		}
	}

	if len(songUrl) == 0 {
		return "", errors.New("no results were found")
	}

	lyrics, err := genius.Lyrics.FindForSong(hits[0].Result.URL)
	if err != nil {
		return "", err
	}

	return lyrics.String(), nil
}

func main() {
	http.Handle("/static/", http.FileServer(http.FS(publicFiles)))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := publicFiles.ReadFile("static/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	})

	http.HandleFunc("GET /lyrics", func(w http.ResponseWriter, r *http.Request) {
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

		searchInput := SearchInput{
			SongName: songName[0],
		}
		if okAlbum {
			searchInput.AlbumName = albumName[0]
		}
		if okArtist {
			searchInput.ArtistName = artistName[0]
		}
		lyricsText, err := getSongLyrics(searchInput)
		if err != nil {
			log.Println("oppsie doopsie some shit happened", err)
			w.Write([]byte("No results were found"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(lyricsText))
	})

	http.ListenAndServe(":"+port, nil)
}
