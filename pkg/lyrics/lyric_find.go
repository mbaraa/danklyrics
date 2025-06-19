package lyrics

import (
	"errors"

	"github.com/mbaraa/lrclibgo"
)

type lyricFindProvider struct {
	client *lrclibgo.Client
}

func (l *lyricFindProvider) GetSongLyrics(s SearchParams) (Lyrics, error) {
	lrcSearch := lrclibgo.SearchParams{
		TrackName:  s.SongName,
		ArtistName: s.ArtistName,
		AlbumName:  s.AlbumName,
		Limit:      0,
	}

	hits, err := l.client.Search.Get(lrcSearch)
	if err != nil {
		return Lyrics{}, err
	}

	if len(hits) == 0 {
		return Lyrics{}, errors.New("no results were found")
	}

	lyrics := hits[0].Lyrics()
	return Lyrics{
		parts:  lyrics.Parts(),
		synced: lyrics.Synced(),
	}, nil
}
