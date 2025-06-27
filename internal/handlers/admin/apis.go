package admin

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/mbaraa/danklyrics/internal/actions"
)

var (
	//go:embed admin_lyrics_requests.html
	lyricsRequestsTemplate embed.FS
	//go:embed admin_lyrics_request.html
	lyricsRequestTemplate embed.FS
)

type api struct {
	usecases *actions.Actions
}

func NewAdminApi(usecases *actions.Actions) *api {
	return &api{
		usecases: usecases,
	}
}

func (a *api) HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
	var reqBody actions.AuthenticateAdminParams
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
		return
	}

	payload, err := a.usecases.AuthenticateAdmin(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("incorrect email or password"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "admin-token",
		Value:   payload.SessionToken,
		Path:    "/",
		Expires: time.Now().UTC().Add(time.Hour * 2),
	})

	w.Header().Set("HX-Redirect", "/")
}

func (a *api) HandleListLyricsRequests(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("admin-token")
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("excuse me?"))
		return
	}

	requests, err := a.usecases.ListLyricsRequests(sessionToken.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("excuse me?"))
		return
	}

	t := template.Must(template.ParseFS(lyricsRequestsTemplate, "admin_lyrics_requests.html"))
	err = t.Execute(w, requests)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("excuse me?"))
		return
	}
}

func (a *api) HandleGetLyricsRequest(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("admin-token")
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("excuse me?"))
		return
	}

	lyricsRequestId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid lyrics request id"))
		return
	}

	lyricsRequest, err := a.usecases.GetLyricsRequest(sessionToken.Value, uint(lyricsRequestId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("excuse me?"))
		return
	}

	t := template.Must(template.ParseFS(lyricsRequestTemplate, "admin_lyrics_request.html"))
	err = t.Execute(w, lyricsRequest.Parts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("excuse me?"))
		return
	}
}

func (a *api) HandleApproveLyricsRequest(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("admin-token")
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("excuse me?"))
		return
	}

	lyricsRequestId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid lyrics request id"))
		return
	}

	err = a.usecases.ApproveLyricsRequest(sessionToken.Value, uint(lyricsRequestId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("excuse me?"))
		return
	}
}

func (a *api) HandleRejectLyricsRequest(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("admin-token")
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("excuse me?"))
		return
	}

	lyricsRequestId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid lyrics request id"))
		return
	}

	err = a.usecases.RejectLyricsRequest(sessionToken.Value, uint(lyricsRequestId), "Lyrics Exists")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("excuse me?"))
		return
	}
}
