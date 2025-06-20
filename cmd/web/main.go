package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/mbaraa/danklyrics/internal/config"
	"github.com/mbaraa/danklyrics/pkg/client"
	"github.com/mbaraa/danklyrics/pkg/provider"
	"github.com/mbaraa/danklyrics/website"
)

var (
	lyricser *client.Http

	publicFiles embed.FS
)

func init() {
	publicFiles = website.FS()

	var err error
	lyricser, err = client.NewHttp(client.Config{
		GeniusClientId:     config.Env().GeniusClientId,
		GeniusClientSecret: config.Env().GeniusClientSecret,
		Providers:          []provider.Name{provider.Dank, provider.LyricFind, provider.Genius},
		ApiAddress:         config.Env().ApiAddress,
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.FS(publicFiles))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := publicFiles.ReadFile("index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
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

		searchInput := provider.SearchParams{
			SongName: songName[0],
		}
		if okAlbum {
			searchInput.AlbumName = albumName[0]
		}
		if okArtist {
			searchInput.ArtistName = artistName[0]
		}

		lyricsText, err := lyricser.GetSongLyrics(searchInput)
		if err != nil {
			log.Println("oppsie doopsie some shit happened", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("No results were found"))
			return
		}

		w.Write([]byte(lyricsText.String()))
	})

	log.Printf("Starting web server at port %s", config.Env().Port)
	log.Fatalln(http.ListenAndServe(":"+config.Env().Port, nil))
}
