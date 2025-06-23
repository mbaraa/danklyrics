package actions

type AuthenticateAdminParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticateAdminPayload struct {
	SessionToken string `json:"session_token"`
}

func (a *Actions) AuthenticateAdmin(params AuthenticateAdminParams) (AuthenticateAdminPayload, error) {
}

type LyricsRequest struct {
	Id           uint              `json:"id"`
	SongTitle    string            `json:"song_title"`
	ArtistName   string            `json:"artist_name"`
	AlbumTitle   string            `json:"album_title"`
	LyricsPlain  []string          `json:"lyrics_plain,omitempty"`
	LyricsSynced map[string]string `json:"lyrics_synced,omitempty"`
}

func (a *Actions) ListLyricsRequests(adminToken string) ([]LyricsRequest, error) {
}

func (a *Actions) GetLyricsRequest(adminToken string, id uint) (LyricsRequest, error) {
}

func (a *Actions) ApproveLyricsRequest(adminToken string, id uint) error {
}
