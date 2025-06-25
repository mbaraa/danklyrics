package mariadb

import (
	"fmt"

	"github.com/mbaraa/danklyrics/internal/actions"
	"github.com/mbaraa/danklyrics/internal/models"
)

func Migrate() error {
	dbConn, err := dbConnector()
	if err != nil {
		return err
	}

	err = dbConn.Debug().AutoMigrate(
		new(models.Lyrics),
		new(models.LyricsPart),
		new(models.LyricsSyncedPart),
		new(models.LyricsRequest),
		new(models.LyricsRequestPart),
		new(models.LyricsRequestSyncedPart),
		new(models.Admin),
	)
	if err != nil {
		return err
	}

	for _, tableName := range []string{
		"lyrics", "lyrics_parts", "lyrics_synced_parts",
		"lyrics_requests", "lyrics_request_parts", "lyrics_request_synced_parts",
	} {
		err = dbConn.Exec("ALTER TABLE " + tableName + " CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error
		if err != nil {
			return err
		}
	}

	for _, columnName := range []string{
		"song_title", "artist_name", "album_title",
	} {
		err = dbConn.Exec(fmt.Sprintf("ALTER TABLE lyrics CHANGE COLUMN `%s` `%s` TEXT CHARACTER SET 'utf8' COLLATE 'utf8_general_ci';", columnName, columnName)).Error
		if err != nil {
			return err
		}
	}

	var l []models.Lyrics
	err = dbConn.Model(new(models.Lyrics)).Find(&l).Error
	if err != nil {
		panic(err)
	}

	fmt.Println("len", len(l))

	duplicatesCount := make(map[string]struct {
		count int
		ids   []uint
	})

	for _, lyrics := range l {
		if lyrics.AlbumTitle != "" {
			lyrics.PublicId = actions.Slugify(fmt.Sprintf("%s-%s-%s", lyrics.ArtistName, lyrics.AlbumTitle, lyrics.SongTitle))
		} else {
			lyrics.PublicId = actions.Slugify(fmt.Sprintf("%s-%s", lyrics.ArtistName, lyrics.SongTitle))
		}
		duplicatesCount[lyrics.PublicId] = struct {
			count int
			ids   []uint
		}{
			count: duplicatesCount[lyrics.PublicId].count + 1,
			ids:   append(duplicatesCount[lyrics.PublicId].ids, lyrics.Id),
		}

		err := dbConn.Model(&lyrics).UpdateColumn("public_id", lyrics.PublicId).Error
		if err != nil {
			panic(err)
		}
	}

	deleted := 0

	for _, ding := range duplicatesCount {
		if ding.count > 1 {
			deleted += ding.count
			err := dbConn.Exec("delete from lyrics_parts where lyrics_id in ?", ding.ids).Error
			if err != nil {
				panic(err)
			}
			err = dbConn.Exec("delete from lyrics_synced_parts where lyrics_id in ?", ding.ids).Error
			if err != nil {
				panic(err)
			}
			err = dbConn.Exec("delete from lyrics where id in ?", ding.ids).Error
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("deleted", deleted)

	return nil
}
