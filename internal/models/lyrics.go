package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Lyrics struct {
	Id       uint   `gorm:"primaryKey;autoIncrement"`
	PublicId string `gorm:"index;unique;not null"`

	SongTitle  string `gorm:"index"`
	ArtistName string `gorm:"index"`
	AlbumTitle string `gorm:"index"`

	LyricsPlain  []string          `gorm:"-"`
	LyricsSynced map[string]string `gorm:"-"`

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`
}

func (l *Lyrics) AfterDelete(tx *gorm.DB) error {
	err := tx.
		Exec("DELETE FROM lyrics_parts WHERE lyrics_id = ?", l.Id).
		Error
	if err != nil {
		return err
	}

	_ = tx.
		Exec("DELETE FROM lyrics_synced_parts WHERE lyrics_id = ?", l.Id).
		Error

	return nil
}

func (l *Lyrics) AfterCreate(tx *gorm.DB) error {
	errs := make([]error, 0, len(l.LyricsPlain)+len(l.LyricsSynced))
	for _, part := range l.LyricsPlain {
		lp := &LyricsPart{
			LyricsId: l.Id,
			Text:     part,
		}

		err := tx.Model(new(LyricsPart)).Create(lp).Error
		errs = append(errs, err)
	}

	for tm, part := range l.LyricsSynced {
		lp := &LyricsSyncedPart{
			LyricsId: l.Id,
			Text:     part,
			Time:     tm,
		}

		err := tx.Model(new(LyricsSyncedPart)).Create(lp).Error
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

type LyricsPart struct {
	LyricsId uint
	Lyrics   Lyrics
	Text     string
}

type LyricsSyncedPart struct {
	LyricsId uint
	Lyrics   Lyrics
	Time     string
	Text     string
}
