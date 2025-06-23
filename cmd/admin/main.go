package main

import (
	"embed"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/mbaraa/danklyrics/internal/config"
	"github.com/mbaraa/danklyrics/pkg/client"
	"github.com/mbaraa/danklyrics/pkg/provider"
	website "github.com/mbaraa/danklyrics/website/admin"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	mjson "github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

var (
	lyricser *client.Http

	publicFiles embed.FS
	minifyer    *minify.M
)

func init() {
	publicFiles = website.FS()

	minifyer = minify.New()
	minifyer.AddFunc("text/css", css.Minify)
	minifyer.AddFunc("text/html", html.Minify)
	minifyer.AddFunc("image/svg+xml", svg.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]json$"), mjson.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

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

func handleIndex(w http.ResponseWriter, r *http.Request) {
	content, err := publicFiles.Open("index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, _ = io.Copy(w, content)
	_ = content.Close()
}

func main() {
	pagesHandler := http.NewServeMux()
	pagesHandler.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.FS(publicFiles))))
	pagesHandler.HandleFunc("/", handleIndex)

	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("GET /lyrics/requests", handleGetSongLyrics)
	apisHandler.HandleFunc("GET /lyrics/request/{id}", handleGetSongLyrics)
	apisHandler.HandleFunc("POST /lyrics/request/approve/{id}", handleGetSongLyrics)
	apisHandler.HandleFunc("POST /auth", handleAuthSubmitLyrics)

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", minifyer.Middleware(pagesHandler))
	applicationHandler.Handle("/api/", http.StripPrefix("/api", apisHandler))

	log.Printf("Starting web server at port %s", config.Env().AdminPort)
	log.Fatalln(http.ListenAndServe(":"+config.Env().AdminPort, applicationHandler))
}
