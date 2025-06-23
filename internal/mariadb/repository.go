package mariadb

import (
	"fmt"

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

func (r *repository) GetLyricsBySongTitle(title string) ([]models.Lyrics, error) {
	lyricses := make([]models.Lyrics, 0)

	err := tryWrapDbError(
		r.client.
			Model(new(models.Lyrics)).
			Find(&lyricses, "LOWER(song_title) LIKE LOWER(?)", likeArg(title)).
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
			Find(&lyricses, "LOWER(song_title) LIKE LOWER(?) AND LOWER(artist_name) LIKE LOWER(?)", likeArg(songTitle), likeArg(artistName)).
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
			Find(&lyricses, "LOWER(song_title) LIKE LOWER(?) AND LOWER(album_title) LIKE LOWER(?)", likeArg(songTitle), likeArg(albumTitle)).
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
			Find(&lyricses, "LOWER(song_title) LIKE LOWER(?) AND LOWER(artist_name) LIKE LOWER(?) AND LOWER(album_title) LIKE LOWER(?)", likeArg(songTitle), likeArg(artistName), likeArg(albumTitle)).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return lyricses, nil
}

func (r *repository) CreateLyricsRequest(l models.LyricsRequest) (models.LyricsRequest, error) {
	err := tryWrapDbError(
		r.client.
			Model(new(models.LyricsRequest)).
			Create(&l).
			Error,
	)
	if err != nil {
		return models.LyricsRequest{}, err
	}

	return l, nil
}

func (r *repository) DeleteLyricsRequest(id uint) error {
	return tryWrapDbError(
		r.client.
			Exec("DELETE FROM lyrics_requests WHERE id = ?", id).
			Error,
	)
}

func (r *repository) GetLyricsRequestById(id uint) (models.LyricsRequest, error) {
	var lyrics models.LyricsRequest

	err := tryWrapDbError(
		r.client.
			Model(new(models.LyricsRequest)).
			First(&lyrics, "id = ?", id).
			Error,
	)
	if err != nil {
		return models.LyricsRequest{}, err
	}

	parts := make([]models.LyricsRequestPart, 0)
	err = tryWrapDbError(
		r.client.
			Model(new(models.LyricsRequestPart)).
			Where("lyrics_request_id = ?", id).
			Find(&parts).
			Error,
	)
	if err != nil {
		return models.LyricsRequest{}, err
	}

	lyrics.LyricsPlain = make([]string, 0, len(parts))
	for _, part := range parts {
		lyrics.LyricsPlain = append(lyrics.LyricsPlain, part.Text)
	}

	synced := make([]models.LyricsRequestSyncedPart, 0)
	err = tryWrapDbError(
		r.client.
			Model(new(models.LyricsRequestSyncedPart)).
			Where("lyrics_request_id = ?", id).
			Find(&synced).
			Error,
	)
	if err != nil {
		return models.LyricsRequest{}, err
	}

	lyrics.LyricsSynced = make(map[string]string, 0)
	for _, part := range synced {
		lyrics.LyricsSynced[part.Time] = part.Text
	}

	return lyrics, nil
}

func (r *repository) GetLyricsRequests() ([]models.LyricsRequest, error) {
	lyricsRequests := make([]models.LyricsRequest, 0)

	err := tryWrapDbError(
		r.client.
			Model(new(models.LyricsRequest)).
			Find(&lyricsRequests).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return lyricsRequests, nil
}

func (r *repository) GetAdminByUsername(username string) (models.Admin, error) {
	var admin models.Admin

	err := tryWrapDbError(
		r.client.
			Model(new(models.Admin)).
			First(&admin, "username = ?", username).
			Error,
	)
	if err != nil {
		return models.Admin{}, err
	}

	return admin, nil
}

func likeArg(arg string) string {
	return fmt.Sprintf("%%%s%%", arg)
}
