package main

import (
	"embed"
	"io"
	"log"
	"net/http"

	"github.com/mbaraa/danklyrics/internal/actions"
	"github.com/mbaraa/danklyrics/internal/config"
	"github.com/mbaraa/danklyrics/internal/handlers"
	"github.com/mbaraa/danklyrics/internal/jwt"
	"github.com/mbaraa/danklyrics/internal/mailer"
	"github.com/mbaraa/danklyrics/internal/mariadb"
	website "github.com/mbaraa/danklyrics/website/admin"
)

var (
	usecases    *actions.Actions
	publicFiles embed.FS
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

	publicFiles = website.FS()
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

func handleRobots(w http.ResponseWriter, r *http.Request) {
	content, err := publicFiles.Open("robots.txt")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = io.Copy(w, content)
	_ = content.Close()
}

func main() {
	pagesHandler := http.NewServeMux()
	pagesHandler.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.FS(publicFiles))))
	pagesHandler.HandleFunc("/", handleIndex)
	pagesHandler.HandleFunc("/robots.txt", handleRobots)

	adminApi := handlers.NewAdminApi(usecases)

	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("GET /lyrics/requests", adminApi.HandleListLyricsRequests)
	apisHandler.HandleFunc("GET /lyrics/request/{id}", adminApi.HandleGetLyricsRequest)
	apisHandler.HandleFunc("POST /lyrics/request/approve/{id}", adminApi.HandleApproveLyricsRequest)
	apisHandler.HandleFunc("POST /lyrics/request/reject/{id}", adminApi.HandleRejectLyricsRequest)
	apisHandler.HandleFunc("POST /auth", adminApi.HandleAuthenticate)

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", pagesHandler)
	applicationHandler.Handle("/api/", http.StripPrefix("/api", apisHandler))

	log.Printf("Starting web server at port %s", config.Env().AdminPort)
	log.Fatalln(http.ListenAndServe(":"+config.Env().AdminPort, applicationHandler))
}
