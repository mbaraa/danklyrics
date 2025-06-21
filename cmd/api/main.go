package main

import (
	"log"
	"net/http"

	"github.com/mbaraa/danklyrics/internal/actions"
	"github.com/mbaraa/danklyrics/internal/config"
	"github.com/mbaraa/danklyrics/internal/handlers"
	"github.com/mbaraa/danklyrics/internal/mariadb"
)

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

func main() {
	apiHandler := http.NewServeMux()

	lyricsApi := handlers.NewLyricsFinderApi(usecases)
	dankLyricsApi := handlers.NewDankLyricsApi(usecases)

	apiHandler.HandleFunc("/", lyricsApi.HandleIndex)
	apiHandler.HandleFunc("GET /providers", lyricsApi.HandleListProviders)
	apiHandler.HandleFunc("GET /lyrics", lyricsApi.HandleGetSongLyrics)
	apiHandler.HandleFunc("GET /dank/lyrics", dankLyricsApi.HandleGetSongLyrics)

	log.Printf("Starting web server at port %s", config.Env().ApiPort)
	log.Fatalln(http.ListenAndServe(":"+config.Env().ApiPort, apiHandler))
}
