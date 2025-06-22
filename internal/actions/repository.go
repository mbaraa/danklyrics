package actions

import "github.com/mbaraa/danklyrics/internal/models"

type Repository interface {
	CreateLyrics(l models.Lyrics) (models.Lyrics, error)
	GetLyricsBySongTitle(Title string) ([]models.Lyrics, error)
	GetLyricsBySongTitleAndArtistName(songTitle, artistName string) ([]models.Lyrics, error)
	GetLyricsBySongAndAlbumTitle(songTitle, albumTitle string) ([]models.Lyrics, error)
	GetLyricsBySongTitleArtistNameAndAlbumTitle(songTitle, artistName, albumTitle string) ([]models.Lyrics, error)

	CreateLyricsRequest(l models.LyricsRequest) (models.LyricsRequest, error)
	GetLyricsRequestById(id uint) (models.LyricsRequest, error)
}
