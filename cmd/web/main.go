package main

import (
	"danklyrics/pkg/lyrics"
	"danklyrics/website"
	"embed"
	"log"
	"net/http"
	"os"
)

var (
	port               = os.Getenv("PORT")
	geniusClientId     = os.Getenv("GENIUS_CLIENT_ID")
	geniusClientSecret = os.Getenv("GENIUS_CLIENT_SECRET")

	lyricser *lyrics.Finder

	publicFiles embed.FS
)

func init() {
	publicFiles = website.FS()

	var err error
	lyricser, err = lyrics.New(lyrics.FinderConfig{
		GeniusClientId:     geniusClientId,
		GeniusClientSecret: geniusClientSecret,
		Providers:          []lyrics.ProviderName{lyrics.ProviderLyricFind, lyrics.ProviderGenius},
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

		searchInput := lyrics.SearchParams{
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
			w.Write([]byte("No results were found"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(lyricsText.String()))
	})

	log.Printf("Starting web server at port %s", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
