package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/mbaraa/danklyrics/internal/config"
	"github.com/mbaraa/danklyrics/internal/handlers/web"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	mjson "github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

var (
	minifyer *minify.M
)

func init() {
	minifyer = minify.New()
	minifyer.AddFunc("text/css", css.Minify)
	minifyer.AddFunc("text/html", html.Minify)
	minifyer.AddFunc("image/svg+xml", svg.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]json$"), mjson.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
}

func main() {
	pages := web.NewPages(nil)
	pagesHandler := http.NewServeMux()
	pagesHandler.Handle("/static/", pages.StaticFiles())
	pagesHandler.HandleFunc("/", pages.HandleIndex)

	apis := web.NewApi(nil)
	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("GET /lyrics", apis.HandleGetSongLyrics)
	apisHandler.HandleFunc("POST /lyrics", apis.HandleSubmitLyrics)
	apisHandler.HandleFunc("POST /auth", apis.HandleAuthSubmitLyrics)
	apisHandler.HandleFunc("GET /auth/confirm", apis.HandleConfirmAuthSubmitLyrics)

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", minifyer.Middleware(pagesHandler))
	applicationHandler.Handle("/api/", http.StripPrefix("/api", apisHandler))

	log.Printf("Starting web server at port %s", config.Env().Port)
	log.Fatalln(http.ListenAndServe(":"+config.Env().Port, applicationHandler))
}
