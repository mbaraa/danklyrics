<div align="center">
  <a href="https://danklyrics.com" target="_blank"><img src="https://danklyrics.com/static/favicon.png" width="150" /></a>

  <h1>DankLyrics</h1>
  <p>
    <strong>A Genius and LyricFind powered lyrics finder with the legendary cs1.6 theme.</strong>
  </p>
  <p>
    <a href="https://goreportcard.com/report/github.com/mbaraa/danklyrics"><img alt="rex-deployment" src="https://goreportcard.com/badge/github.com/mbaraa/danklyrics"/></a>
    <a href="https://godoc.org/github.com/mbaraa/danklyrics"><img alt="rex-deployment" src="https://godoc.org/github.com/mbaraa/danklyrics?status.png"/></a>
    <a href="https://github.com/mbaraa/danklyrics/actions/workflows/rex-deploy.yml"><img alt="rex-deployment" src="https://github.com/mbaraa/danklyrics/actions/workflows/rex-deploy.yml/badge.svg"/></a>
  </p>
</div>

## About

**DankLyrics:** A lyrics finder API, Website and Go package!

# Go Package Docs

DankLyrics provides a Go package, since the project is written in Go lol.

Here's a sample usage, it's pretty straight forward, as the client only has one method :)

```go
package main

import (
	"github.com/mbaraa/danklyrics/pkg/client"
	"github.com/mbaraa/danklyrics/pkg/provider"
)

func main() {
	lyricser, err := client.NewHttp(client.Config{
        // available providers are the following in addition to Genius, and you need to provide client id and token to use it.
		Providers: []provider.Name{provider.Dank, provider.LyricFind},
	})
	if err != nil {
		panic(err)
	}

	searchInput := provider.SearchParams{
		SongName: "sos",
        ArtistName: "abba",
	}
    lyrics, err := lyricser.GetSongLyrics(searchInput)
	if err != nil {
		panic(err)
	}

    fmt.Println(lyrics.String())
    fmt.Println(lyrics.Parts)
    fmt.Println(lyrics.Synced)
}
```

# REST API Docs

_Rest API is available at [api.danklyrics.com](https://api.danklyrics.com)_

- **`GET /`**:

_Displays this message_

```
refer to (https://github.com/mbaraa/danklyrics) for API docs!
```

- **`GET /providers`**:

_Returns a list of the current supported lyrics providers_

```json
[
  {
    "name": "string: Name of the provider",
    "id": "string: id to specify the provider to use in /lyrics"
  }
]
```

- **`GET /lyrics`**:

_Finds lyrics for a song using the specified providers_

Query parameters

| name                   | required              | description                                                                        |
| ---------------------- | --------------------- | ---------------------------------------------------------------------------------- |
| `providers`            | required (at least 1) | to specify which lyrics provider(s) to use, list is fetched from `GET /providers`. |
| `genius_client_id`     | conditional           | to set genius' client id when `genius` is used as a provider.                      |
| `genius_client_secret` | conditional           | same as `genius_client_id` but for client secret.                                  |
| `song`                 | required              | song's name to search for.                                                         |
| `artist`               | optional              | artist's name to search for, if the song's name isn't enough.                      |
| `album`                | optional              | album's name to search for, if the song's name isn't enough.                       |

```json
{
  "parts": ["string lyrics parts of the song"],
  "synced": { "time": "lyrics part" }
}
```

- **`GET /dank/lyrics`**:

_Find lyrics from DankLyrics' database, equivalent to using the Go client with `provider.Dank` set_

Query parameters

| name     | required | description                                                   |
| -------- | -------- | ------------------------------------------------------------- |
| `song`   | required | song's name to search for.                                    |
| `artist` | optional | artist's name to search for, if the song's name isn't enough. |
| `album`  | optional | album's name to search for, if the song's name isn't enough.  |

```json
[
  {
    "song_name": "string: represents the song's name",
    "artist_name": "string: represents the song artist's name",
    "album_name": "string: represents the song album's name",
    "parts": ["string lyrics parts of the song"],
    "synced": { "time": "lyrics part" }
  }
]
```

---

A [DankStuff <img height="16" width="16" src="https://dankstuff.net/assets/favicon.ico" />](https://dankstuff.net) product!

Made with 🧉 by [Baraa Al-Masri](https://mbaraa.com).
