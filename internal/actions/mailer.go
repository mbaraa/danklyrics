package actions

import "github.com/mbaraa/danklyrics/internal/models"

type Mailer interface {
	SendVerificationEmail(token, email string) error
	SendLyricsApprovedEmail(lyrics models.Lyrics, email string) error
	SendLyricsRejectedEmail(reason, email string) error
}
