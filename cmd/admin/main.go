package main

import (
	"log"
	"net/http"

	"github.com/mbaraa/danklyrics/internal/actions"
	"github.com/mbaraa/danklyrics/internal/config"
	"github.com/mbaraa/danklyrics/internal/handlers/admin"
	"github.com/mbaraa/danklyrics/internal/jwt"
	"github.com/mbaraa/danklyrics/internal/mailer"
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

	mailUtil := mailer.New()
	jwtUtil := jwt.New[actions.TokenPayload]()
	usecases = actions.New(repo, mailUtil, jwtUtil)
}

func main() {
	pages := admin.NewAdminPages(usecases)
	pagesHandler := http.NewServeMux()
	pagesHandler.Handle("/static/", pages.StaticFiles())
	pagesHandler.HandleFunc("/", pages.HandleIndex)
	pagesHandler.HandleFunc("/robots.txt", pages.HandleIndex)

	apis := admin.NewAdminApi(usecases)
	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("GET /lyrics/requests", apis.HandleListLyricsRequests)
	apisHandler.HandleFunc("GET /lyrics/request/{id}", apis.HandleGetLyricsRequest)
	apisHandler.HandleFunc("POST /lyrics/request/approve/{id}", apis.HandleApproveLyricsRequest)
	apisHandler.HandleFunc("POST /lyrics/request/reject/{id}", apis.HandleRejectLyricsRequest)
	apisHandler.HandleFunc("POST /auth", apis.HandleAuthenticate)

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", pagesHandler)
	applicationHandler.Handle("/api/", http.StripPrefix("/api", apisHandler))

	log.Printf("Starting web server at port %s", config.Env().AdminPort)
	log.Fatalln(http.ListenAndServe(":"+config.Env().AdminPort, applicationHandler))
}
