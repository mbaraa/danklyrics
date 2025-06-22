package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/mbaraa/danklyrics/internal/config"
)

//go:embed verification_email.html
var tmpl embed.FS

type SmtpMailer struct {
}

func New() *SmtpMailer {
	return &SmtpMailer{}
}

type verificationEmailData struct {
	ConfirmationLink string
}

func (s *SmtpMailer) SendVerificationEmail(token, email string) error {
	t := template.Must(template.ParseFS(tmpl, "verification_email.html"))

	emailBuf := bytes.NewBuffer(nil)
	err := t.Execute(emailBuf, verificationEmailData{
		ConfirmationLink: "https://danklyrics.com/api/auth/confirm?token=" + token,
	})
	if err != nil {
		return err
	}

	return sendEmail("Submit Lyrics Authentication", emailBuf.String(), email)
}

func sendEmail(subject, content, to string) error {
	receiver := []string{to}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	_subject := "Subject: " + subject
	_to := "To: " + to
	_from := fmt.Sprintf("From: Baraa from DankLyrics <%s>", config.Env().Smtp.Username)
	body := fmt.Appendf([]byte{}, "%s\n%s\n%s\n%s\n%s", _from, _to, _subject, mime, content)

	addr := config.Env().Smtp.Host + ":" + config.Env().Smtp.Port
	auth := smtp.PlainAuth("", config.Env().Smtp.Username, config.Env().Smtp.Password, config.Env().Smtp.Host)

	return smtp.SendMail(addr, auth, config.Env().Smtp.Username, receiver, body)
}
