package genius

import (
	"danklyrics/pkg/models"
	"danklyrics/pkg/provider"
	"errors"
	"fmt"
	"strconv"

	"github.com/mbaraa/gonius"
)

type geniusProvider struct {
	client *gonius.Client
}

func New(clientId, clientSecret string) provider.Service {
	return &geniusProvider{
		client: gonius.NewClient(clientId, clientSecret),
	}
}

func (g *geniusProvider) GetSongLyrics(s provider.SearchParams) (models.Lyrics, error) {
	var hits []gonius.Hit
	var err error

	okArtist := s.ArtistName != ""
	okAlbum := s.AlbumName != ""
	okSong := s.SongName != ""

	switch {
	case !okArtist && !okAlbum && okSong:
		hits, err = g.client.Search.Get(s.SongName)
		if err != nil {
			return models.Lyrics{}, err
		}
	case okArtist && !okAlbum && okSong:
		hits, err = g.client.Search.Get(fmt.Sprintf("%s %s", s.SongName, s.ArtistName))
		if err != nil {
			return models.Lyrics{}, err
		}
	case !okArtist && okAlbum && okSong:
		hits, err = g.client.Search.Get(fmt.Sprintf("%s %s", s.SongName, s.AlbumName))
		if err != nil {
			return models.Lyrics{}, err
		}
	case okArtist && okAlbum && okSong:
		hits, err = g.client.Search.Get(fmt.Sprintf("%s %s %s", s.SongName, s.AlbumName, s.ArtistName))
		if err != nil {
			return models.Lyrics{}, err
		}
	}

	if len(hits) == 0 {
		return models.Lyrics{}, errors.New("no results were found")
	}

	song, err := g.client.Songs.Get(strconv.Itoa(hits[0].Result.Id))
	if err != nil {
		return models.Lyrics{}, err
	}

	lyrics, err := g.client.Lyrics.FindForSong(hits[0].Result.URL)
	if err != nil {
		return models.Lyrics{}, err
	}

	artistName := ""
	if song.PrimaryArtist != nil {
		artistName = song.PrimaryArtist.Name
	}
	albumTitle := ""
	if song.Album != nil {
		albumTitle = song.Album.FullTitle
	}

	return models.Lyrics{
		SongName:   song.Title,
		ArtistName: artistName,
		AlbumName:  albumTitle,
		Parts:      lyrics.Parts(),
	}, nil
}
