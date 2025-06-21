package mariadb

import (
	"github.com/mbaraa/danklyrics/internal/models"

	"gorm.io/gorm"
)

type repository struct {
	client *gorm.DB
}

func New() (*repository, error) {
	conn, err := dbConnector()
	if err != nil {
		return nil, err
	}

	return &repository{
		client: conn,
	}, nil
}

func (r *repository) CreateLyrics(lyrics models.Lyrics) (models.Lyrics, error) {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Lyrics)).
			Create(&lyrics).
			Error,
	)
	if err != nil {
		return models.Lyrics{}, err
	}

	return lyrics, nil
}

func (r *repository) GetLyricsById(id uint) (models.Lyrics, error) {
	var lyrics models.Lyrics

	err := tryWrapDbError(
		r.client.
			Model(new(models.Lyrics)).
			First(&lyrics, "id = ?", id).
			Error,
	)
	if err != nil {
		return models.Lyrics{}, err
	}

	return lyrics, nil
}

func (r *repository) GetLyricsBySongTitle(title string) ([]models.Lyrics, error) {
	lyricses := make([]models.Lyrics, 0)

	err := tryWrapDbError(
		r.client.
			Model(new(models.Lyrics)).
			Find(&lyricses, "song_title LIKE '%"+title+"%'").
			Error,
	)
	if err != nil {
		return nil, err
	}

	return lyricses, nil
}

func (r *repository) GetLyricsBySongTitleAndArtistName(songTitle, artistName string) ([]models.Lyrics, error) {
	lyricses := make([]models.Lyrics, 0)

	err := tryWrapDbError(
		r.client.
			Model(new(models.Lyrics)).
			Find(&lyricses, "song_title LIKE '%"+songTitle+"%' AND artist_name LIKE '%"+artistName+"%'").
			Error,
	)
	if err != nil {
		return nil, err
	}

	return lyricses, nil
}

func (r *repository) GetLyricsBySongAndAlbumTitle(songTitle, albumTitle string) ([]models.Lyrics, error) {
	lyricses := make([]models.Lyrics, 0)

	err := tryWrapDbError(
		r.client.
			Model(new(models.Lyrics)).
			Find(&lyricses, "song_title LIKE '%"+songTitle+"%' AND album_title LIKE '%"+albumTitle+"%'").
			Error,
	)
	if err != nil {
		return nil, err
	}

	return lyricses, nil
}

func (r *repository) GetLyricsBySongTitleArtistNameAndAlbumTitle(songTitle, artistName, albumTitle string) ([]models.Lyrics, error) {
	lyricses := make([]models.Lyrics, 0)

	err := tryWrapDbError(
		r.client.
			Model(new(models.Lyrics)).
			Find(&lyricses, "song_title LIKE '%"+songTitle+"%' AND artist_name LIKE '%"+artistName+"%' AND album_title LIKE '%"+albumTitle+"%'").
			Error,
	)
	if err != nil {
		return nil, err
	}

	return lyricses, nil
}
