package models

import "strings"

// Lyrics holds lyrics fetched by a provider.
type Lyrics struct {
	// Parts holds the lines of the lyrics.
	Parts []string `json:"parts"`
	// Synced similar to [Lyrics.Parts] but instead of plain lines, it has time stamps syncs for the line.
	Synced map[string]string `json:"synced,omitempty"`
}

// String returns the lyrics' lines as one single string, separated by line feed.
func (l *Lyrics) String() string {
	return strings.TrimSpace(strings.Join(l.Parts, "\n"))
}
