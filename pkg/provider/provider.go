package provider

import "danklyrics/pkg/models"

// SearchParams holds the search criteria to find a song from a provider.
type SearchParams struct {
	SongName   string
	ArtistName string
	AlbumName  string
}

// Service fetches lyrics for the given song in the search params.
type Service interface {
	// GetSongLyrics searches for a song's lyrics using the given search params.
	GetSongLyrics(s SearchParams) (models.Lyrics, error)
	// GetSongsLyrics same as [GetSongLyrics] but returns all the songs in the search results with their lyrics.
	// GetSongsLyrics(s SearchParams) ([]models.Song, error)
}

// Name represents lyrics finding providers to choose from when doing a lyrics search.
type Name string

const (
	// Genius pass this to [GetSongLyrics] to use Genius as a lyrics provider.
	Genius Name = "genius"
	// LyricFind pass this to [GetSongLyrics] to use LyricFind as a lyrics provider.
	LyricFind Name = "lrc"
)
