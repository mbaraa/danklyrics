package actions

import (
	"errors"

	intmodels "github.com/mbaraa/danklyrics/internal/models"
	"github.com/mbaraa/danklyrics/pkg/models"
)

func (a *Actions) GetLyricsByPublicId(id string) (models.Lyrics, error) {
	intLyrics, err := a.repo.GetLyricsByPublicId(id)
	if err != nil {
		return models.Lyrics{}, err
	}

	return models.Lyrics{
		PublicId:   intLyrics.PublicId,
		SongName:   intLyrics.SongTitle,
		ArtistName: intLyrics.ArtistName,
		AlbumName:  intLyrics.AlbumTitle,
		Parts:      intLyrics.LyricsPlain,
		Synced:     intLyrics.LyricsSynced,
	}, nil
}

func (a *Actions) GetLyricsBySongTitle(title string) ([]models.Lyrics, error) {
	intLyricses, err := a.repo.GetLyricsBySongTitle(title)
	if err != nil {
		return nil, err
	}

	lyricses := make([]models.Lyrics, 0, len(intLyricses))
	for _, intLyrics := range intLyricses {
		lyricses = append(lyricses, models.Lyrics{
			SongName:   intLyrics.SongTitle,
			ArtistName: intLyrics.ArtistName,
			AlbumName:  intLyrics.AlbumTitle,
			Parts:      intLyrics.LyricsPlain,
			Synced:     intLyrics.LyricsSynced,
		})
	}

	return lyricses, nil
}

func (a *Actions) GetLyricsBySongTitleAndArtistName(title, artistName string) ([]models.Lyrics, error) {
	intLyricses, err := a.repo.GetLyricsBySongTitleAndArtistName(title, artistName)
	if err != nil {
		return nil, err
	}

	lyricses := make([]models.Lyrics, 0, len(intLyricses))
	for _, intLyrics := range intLyricses {
		lyricses = append(lyricses, models.Lyrics{
			SongName:   intLyrics.SongTitle,
			ArtistName: intLyrics.ArtistName,
			AlbumName:  intLyrics.AlbumTitle,
			Parts:      intLyrics.LyricsPlain,
			Synced:     intLyrics.LyricsSynced,
		})
	}

	return lyricses, nil
}

func (a *Actions) GetLyricsBySongTitleAndAlbumTitle(title, albumTitle string) ([]models.Lyrics, error) {
	intLyricses, err := a.repo.GetLyricsBySongAndAlbumTitle(title, albumTitle)
	if err != nil {
		return nil, err
	}

	lyricses := make([]models.Lyrics, 0, len(intLyricses))
	for _, intLyrics := range intLyricses {
		lyricses = append(lyricses, models.Lyrics{
			SongName:   intLyrics.SongTitle,
			ArtistName: intLyrics.ArtistName,
			AlbumName:  intLyrics.AlbumTitle,
			Parts:      intLyrics.LyricsPlain,
			Synced:     intLyrics.LyricsSynced,
		})
	}

	return lyricses, nil
}

func (a *Actions) GetLyricsBySongTitleArtistNameAndAlbumTitle(title, artistName, albumTitle string) ([]models.Lyrics, error) {
	intLyricses, err := a.repo.GetLyricsBySongTitleArtistNameAndAlbumTitle(title, artistName, albumTitle)
	if err != nil {
		return nil, err
	}

	lyricses := make([]models.Lyrics, 0, len(intLyricses))
	for _, intLyrics := range intLyricses {
		lyricses = append(lyricses, models.Lyrics{
			SongName:   intLyrics.SongTitle,
			ArtistName: intLyrics.ArtistName,
			AlbumName:  intLyrics.AlbumTitle,
			Parts:      intLyrics.LyricsPlain,
			Synced:     intLyrics.LyricsSynced,
		})
	}

	return lyricses, nil
}

func (a *Actions) CreateLyrics(l models.Lyrics) (models.Lyrics, error) {
	if l.SongName == "" {
		return models.Lyrics{}, errors.New("missing song name")
	}

	intLyrics := intmodels.Lyrics{
		SongTitle:    l.SongName,
		ArtistName:   l.ArtistName,
		AlbumTitle:   l.AlbumName,
		LyricsPlain:  l.Parts,
		LyricsSynced: l.Synced,
	}

	_, err := a.repo.CreateLyrics(intLyrics)
	if err != nil {
		return models.Lyrics{}, err
	}

	return models.Lyrics{}, nil
}

func (a *Actions) CreateLyricsRequest(token string, l models.Lyrics) error {
	tokenDecoded, err := a.jwt.Decode(token, JwtSessionToken)
	if err != nil {
		return err
	}

	if l.SongName == "" {
		return errors.New("missing song name")
	}

	intLyrics := intmodels.LyricsRequest{
		SongTitle:      l.SongName,
		ArtistName:     l.ArtistName,
		AlbumTitle:     l.AlbumName,
		LyricsPlain:    l.Parts,
		LyricsSynced:   l.Synced,
		RequesterEmail: tokenDecoded.Payload.Email,
	}

	_, err = a.repo.CreateLyricsRequest(intLyrics)
	if err != nil {
		return err
	}

	return nil
}
