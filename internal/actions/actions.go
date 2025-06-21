package actions

type Actions struct {
	repo Repository
}

func New(repo Repository) *Actions {
	return &Actions{
		repo: repo,
	}
}
