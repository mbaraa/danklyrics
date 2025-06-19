package lyrics

import (
	"errors"
	"fmt"

	"github.com/mbaraa/gonius"
)

type geniusProvider struct {
	client *gonius.Client
}

func (g *geniusProvider) GetSongLyrics(s SearchParams) (Lyrics, error) {
	var hits []gonius.Hit
	var err error

	okArtist := s.ArtistName != ""
	okAlbum := s.AlbumName != ""
	okSong := s.SongName != ""

	switch {
	case !okArtist && !okAlbum && okSong:
		hits, err = g.client.Search.Get(s.SongName)
		if err != nil {
			return Lyrics{}, err
		}
	case okArtist && !okAlbum && okSong:
		hits, err = g.client.Search.Get(fmt.Sprintf("%s %s", s.SongName, s.ArtistName))
		if err != nil {
			return Lyrics{}, err
		}
	case !okArtist && okAlbum && okSong:
		hits, err = g.client.Search.Get(fmt.Sprintf("%s %s", s.SongName, s.AlbumName))
		if err != nil {
			return Lyrics{}, err
		}
	case okArtist && okAlbum && okSong:
		hits, err = g.client.Search.Get(fmt.Sprintf("%s %s %s", s.SongName, s.AlbumName, s.ArtistName))
		if err != nil {
			return Lyrics{}, err
		}
	}

	if len(hits) == 0 {
		return Lyrics{}, errors.New("no results were found")
	}

	lyrics, err := g.client.Lyrics.FindForSong(hits[0].Result.URL)
	if err != nil {
		return Lyrics{}, err
	}

	return Lyrics{
		parts: lyrics.Parts(),
	}, nil
}
