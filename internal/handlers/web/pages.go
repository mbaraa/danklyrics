package web

import (
	"embed"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/mbaraa/danklyrics/internal/actions"
	website "github.com/mbaraa/danklyrics/website/user"
)

var publicFiles embed.FS

func init() {
	publicFiles = website.FS()
}

type pages struct {
	usecases *actions.Actions
}

func NewPages(usecases *actions.Actions) *pages {
	return &pages{
		usecases: usecases,
	}
}

func (p *pages) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".go") || strings.Contains(r.URL.Path, "admin") {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
		return
	}

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

func (p *pages) StaticFiles() http.Handler {
	return http.StripPrefix("/static", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, ".go") || strings.Contains(r.URL.Path, "admin") {
			http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusTemporaryRedirect)
			return
		}

		http.FileServer(http.FS(publicFiles)).ServeHTTP(w, r)
	}))
}
