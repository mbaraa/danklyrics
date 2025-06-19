package lyrics

import (
	"errors"
	"strings"

	"github.com/mbaraa/gonius"
	"github.com/mbaraa/lrclibgo"
)

// Lyrics holds lyrics fetched by a provider.
type Lyrics struct {
	parts  []string
	synced map[string]string
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

// SearchParams holds the search criteria to find a song from a provider.
type SearchParams struct {
	SongName   string
	ArtistName string
	AlbumName  string
}

// Provider fetches lyrics for the given song in the search params.
type Provider interface {
	GetSongLyrics(s SearchParams) (Lyrics, error)
}

// ProviderName represents lyrics finding providers to choose from when doing a lyrics search.
type ProviderName string

const (
	// ProviderGenius pass this to [GetSongLyrics] to use Genius as a lyrics provider.
	ProviderGenius ProviderName = "genius"
	// ProviderLyricFind pass this to [GetSongLyrics] to use LyricFind as a lyrics provider.
	ProviderLyricFind ProviderName = "lrc"
)

// Finder finds lyrics for a song using the enabled providers.
type Finder struct {
	providers     map[ProviderName]Provider
	providerNames []ProviderName
}

// FinderConfig holds [Finder] configs.
type FinderConfig struct {
	GeniusClientId     string
	GeniusClientSecret string
	Providers          []ProviderName
}

// New creates a new [Finder] with the selected configs.
func New(config FinderConfig) (*Finder, error) {
	if len(config.Providers) == 0 {
		return nil, errors.New("must specify at least one lyrics provider")
	}

	l := &Finder{
		providers:     make(map[ProviderName]Provider),
		providerNames: config.Providers,
	}

	for _, providerName := range config.Providers {
		switch providerName {
		case ProviderGenius:
			l.providers[ProviderGenius] = &geniusProvider{
				client: gonius.NewClient(config.GeniusClientId, config.GeniusClientSecret),
			}
		case ProviderLyricFind:
			l.providers[ProviderLyricFind] = &lyricFindProvider{
				client: lrclibgo.NewClient(),
			}
		}
	}

	return l, nil
}

// GetSongLyrics search for song's lyrics using the enabled providers list,
// where using a provider depends on the provider's order in that list.
//
// returns [Lyrics] and an occurring [error]
func (l *Finder) GetSongLyrics(s SearchParams) (Lyrics, error) {
	lyrics, err := l.providers[l.providerNames[0]].GetSongLyrics(s)
	if err != nil {
		if len(l.providers) <= 1 {
			return Lyrics{}, err
		}

		for _, provider := range l.providerNames[1:] {
			lyrics, err = l.providers[provider].GetSongLyrics(s)
			if err == nil {
				break
			}
		}
	}

	return lyrics, nil
}
