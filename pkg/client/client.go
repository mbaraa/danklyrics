package client

import (
	"danklyrics/internal/providers/genius"
	"danklyrics/internal/providers/lyricfind"
	"danklyrics/pkg/finder"
	"danklyrics/pkg/models"
	"danklyrics/pkg/provider"
)

// Lyricser is the dank lyrics finding client that uses [finder.Service] to find lyrics using the enabled providers.
type Lyricser struct {
	finder *finder.Service
}

// LyricserConfig holds the configs needed to initialize [Lyricser].
type LyricserConfig struct {
	GeniusClientId     string
	GeniusClientSecret string
	Providers          []provider.Name
}

// New initializes a new [Lyricser] instance with the given configs.
func New(c LyricserConfig) (*Lyricser, error) {
	providers := make([]provider.Service, 0, len(c.Providers))

	for _, providerName := range c.Providers {
		switch providerName {
		case provider.Genius:
			providers = append(providers, genius.New(c.GeniusClientId, c.GeniusClientSecret))
		case provider.LyricFind:
			providers = append(providers, lyricfind.New())
		}
	}

	finder, err := finder.New(providers)
	if err != nil {
		return nil, err
	}

	return &Lyricser{
		finder: finder,
	}, nil
}

// GetSongLyrics search for song's lyrics using the enabled providers list,
// where using a provider depends on the provider's order in that list.
//
// returns [Lyrics] and an occurring [error]
func (c *Lyricser) GetSongLyrics(s provider.SearchParams) (models.Lyrics, error) {
	return c.finder.GetSongLyrics(s)
}
