package handlers

import (
	"net/http"
	"strconv"

	"github.com/mbaraa/danklyrics/internal/actions"
)

type adminApi struct {
	usecases *actions.Actions
}

func NewAdminApi(usecases *actions.Actions) *adminApi {
	return &adminApi{
		usecases: usecases,
	}
}

func (a *adminApi) HandleAuthenticate(w http.ResponseWriter, r *http.Request) {

}

func (a *adminApi) HandleListLyricsRequests(w http.ResponseWriter, r *http.Request) {

}

func (a *adminApi) HandleGetLyricsRequest(w http.ResponseWriter, r *http.Request) {
	lyricsRequestId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid lyrics request id"))
		return
	}

	_ = lyricsRequestId

}
