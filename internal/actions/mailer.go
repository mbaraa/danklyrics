package actions

type Mailer interface {
	SendVerificationEmail(token, email string) error
}
