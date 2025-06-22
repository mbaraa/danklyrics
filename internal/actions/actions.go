package actions

type Actions struct {
	repo   Repository
	mailer Mailer
	jwt    JwtManager[TokenPayload]
}

func New(repo Repository, mailer Mailer, jwt JwtManager[TokenPayload]) *Actions {
	return &Actions{
		repo:   repo,
		mailer: mailer,
		jwt:    jwt,
	}
}
