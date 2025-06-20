package models

import "strings"

// Lyrics holds lyrics fetched by a provider.
type Lyrics struct {
	parts  []string
	synced map[string]string
}

func NewLyrics(parts []string, synced map[string]string) Lyrics {
	return Lyrics{
		parts:  parts,
		synced: synced,
	}
}

// String returns the lyrics' lines as one single string, separated by line feed.
func (l *Lyrics) String() string {
	return strings.TrimSpace(strings.Join(l.parts, "\n"))
}

// Parts returns the lines of the lyrics.
func (l *Lyrics) Parts() []string {
	return l.parts
}

// Synced similar to [Lyrics.Parts] but instead of plain lines, it has time stamps syncs for the line.
func (l *Lyrics) Synced() map[string]string {
	return l.synced
}
