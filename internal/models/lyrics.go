package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Lyrics struct {
	Id uint `gorm:"primaryKey;autoIncrement"`

	SongTitle  string `gorm:"index"`
	ArtistName string `gorm:"index"`
	AlbumTitle string `gorm:"index"`

	LyricsPlain  []string          `gorm:"-"`
	LyricsSynced map[string]string `gorm:"-"`

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`
}

func (l *Lyrics) AfterFind(tx *gorm.DB) error {
	parts := make([]LyricsPart, 0)
	err := tx.
		Model(new(LyricsPart)).
		Where("lyrics_id = ?", l.Id).
		Find(&parts).
		Error
	if err != nil {
		return err
	}

	l.LyricsPlain = make([]string, 0, len(parts))
	for _, part := range parts {
		l.LyricsPlain = append(l.LyricsPlain, part.Text)
	}

	synced := make([]LyricsSyncedPart, 0)
	err = tx.
		Model(new(LyricsSyncedPart)).
		Where("lyrics_id = ?", l.Id).
		Find(&synced).
		Error
	if err != nil {
		return err
	}

	l.LyricsSynced = make(map[string]string, 0)
	for _, part := range synced {
		l.LyricsSynced[part.Time] = part.Text
	}

	return nil
}

func (l *Lyrics) AfterDelete(tx *gorm.DB) error {
	err := tx.
		Exec("DELETE FROM lyrics_parts WHERE lyrics_id = ?", l.Id).
		Error
	if err != nil {
		return err
	}

	return tx.
		Exec("DELETE FROM lyrics_synced_parts WHERE lyrics_id = ?", l.Id).
		Error
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
