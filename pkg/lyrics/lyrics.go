package lyrics

import (
	"errors"
	"fmt"
	"os"

	"github.com/mbaraa/gonius"
	"github.com/mbaraa/lrclibgo"
)

var (
	geniusToken = os.Getenv("GENIUS_TOKEN")

	genius *gonius.Client
	lrclib *lrclibgo.Client
)

func init() {
	genius = gonius.NewClient(geniusToken)
	lrclib = lrclibgo.NewClient()
}

type SearchInput struct {
	SongName   string
	ArtistName string
	AlbumName  string
}

func GetSongLyrics(s SearchInput) (string, error) {
	lyrics, err := getSongLyricsGenius(s)
	if err != nil {
		return getSongLyricsLyricFind(s)
	}

	return lyrics, nil
}

func getSongLyricsGenius(s SearchInput) (string, error) {
	var hits []gonius.Hit
	var err error

	okArtist := s.ArtistName != ""
	okAlbum := s.AlbumName != ""
	okSong := s.SongName != ""

	switch {
	case !okArtist && !okAlbum && okSong:
		hits, err = genius.Search.Get(s.SongName)
		if err != nil {
			return "", err
		}
	case okArtist && !okAlbum && okSong:
		hits, err = genius.Search.Get(fmt.Sprintf("%s %s", s.SongName, s.ArtistName))
		if err != nil {
			return "", err
		}
	case !okArtist && okAlbum && okSong:
		hits, err = genius.Search.Get(fmt.Sprintf("%s %s", s.SongName, s.AlbumName))
		if err != nil {
			return "", err
		}
	case okArtist && okAlbum && okSong:
		hits, err = genius.Search.Get(fmt.Sprintf("%s %s %s", s.SongName, s.AlbumName, s.ArtistName))
		if err != nil {
			return "", err
		}
	}

	if len(hits) == 0 {
		return "", errors.New("no results were found")
	}

	lyrics, err := genius.Lyrics.FindForSong(hits[0].Result.URL)
	if err != nil {
		return "", err
	}

	return lyrics.String(), nil
}

func getSongLyricsLyricFind(s SearchInput) (string, error) {
	lrcSearch := lrclibgo.SearchParams{
		TrackName:  s.SongName,
		ArtistName: s.ArtistName,
		AlbumName:  s.AlbumName,
		Limit:      0,
	}

	hits, err := lrclib.Search.Get(lrcSearch)
	if err != nil {
		return "", err
	}

	if len(hits) == 0 {
		return "", errors.New("no results were found")
	}

	lyrics := hits[0].Lyrics()
	return lyrics.String(), nil
}
