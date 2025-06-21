<div align="center">
  <a href="https://danklyrics.com" target="_blank"><img src="https://danklyrics.com/static/favicon.png" width="150" /></a>

  <h1>DankLyrics</h1>
  <p>
    <strong>A Genius and LyricFind powered lyrics finder with the legendary cs1.6 theme.</strong>
  </p>
  <p>
    <a href="https://github.com/mbaraa/danklyrics/actions/workflows/rex-deploy.yml"><img alt="rex-deployment" src="https://github.com/mbaraa/danklyrics/actions/workflows/rex-deploy.yml/badge.svg"/></a>
  </p>
</div>

## About

**DankLyrics:** A lyrics finder API, Website and Go package!

# REST API Docs

- **`GET /`**: displays this message.

```
refer to (https://github.com/mbaraa/danklyrics) for API docs!
```

- **`GET /providers`**: returns a list of the current supported lyrics providers.

```json
[
  {
    "name": "string: Name of the provider",
    "id": "string: id to specify the provider to use in /lyrics"
  }
]
```

- **`GET /lyrics`**:

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

---

A [DankStuff <img height="16" width="16" src="https://dankstuff.net/assets/favicon.ico" />](https://dankstuff.net) product!

Made with 🧉 by [Baraa Al-Masri](https://mbaraa.com).
