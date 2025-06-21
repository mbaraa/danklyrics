package main

import (
	"danklyrics/pkg/client"
	"danklyrics/pkg/provider"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slices"
)

const docsLink = "https://github.com/mbaraa/danklyrics"

var (
	port = os.Getenv("API_PORT")
)

type errorResponse struct {
	Message         string `json:"message"`
	SuggestedAction string `json:"suggested_action,omitempty"`
	DocsLink        string `json:"docs_link,omitempty"`
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("refer to (" + docsLink + ") for API docs!"))
}

func handleListProviders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode([]map[string]string{
		{"name": "LyricFind", "id": "lrc"},
		{"name": "Genius", "id": "genius"},
	})
}

func handleGetSongLyrics(w http.ResponseWriter, r *http.Request) {
	providers := r.URL.Query()["providers"]
	geniusClientId := r.URL.Query().Get("genius_client_id")
	geniusClientSecret := r.URL.Query().Get("genius_client_secret")

	artistName, okArtist := r.URL.Query()["artist"]
	albumName, okAlbum := r.URL.Query()["album"]
	songName, okSong := r.URL.Query()["song"]

	w.Header().Set("Content-Type", "application/json")

	if len(providers) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:         "You must specify at least one provider",
			SuggestedAction: "Check the `GET /providers` endpoint.",
			DocsLink:        docsLink,
		})
		return
	}
	if slices.Contains(providers, "genius") && (geniusClientId == "" || geniusClientSecret == "") {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:         "You must specify genius' client id and secret when using genius as a provider.",
			SuggestedAction: "Visit https://docs.genius.com/",
			DocsLink:        docsLink,
		})
		return
	}

	if !okArtist && !okAlbum && !okSong {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:  "Missing all query parameters `artist`, `album` and `song`",
			DocsLink: docsLink,
		})
		return
	}

	if !okSong {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:  "Missing required query parameter `song`",
			DocsLink: docsLink,
		})
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

	providersConfig := make([]provider.Name, 0, len(providers))
	for _, p := range providers {
		providersConfig = append(providersConfig, provider.Name(p))
	}
	lyricser, err := client.New(client.Config{
		GeniusClientId:     geniusClientId,
		GeniusClientSecret: geniusClientSecret,
		Providers:          providersConfig,
	})

	lyrics, err := lyricser.GetSongLyrics(searchInput)
	if err != nil {
		log.Println("oppsie doopsie some shit happened", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "No results were found",
		})
		return
	}

	_ = json.NewEncoder(w).Encode(lyrics)
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("GET /providers", handleListProviders)
	http.HandleFunc("GET /lyrics", handleGetSongLyrics)

	log.Printf("Starting web server at port %s", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
