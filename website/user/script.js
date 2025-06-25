"use strict";

const form = document.getElementById("dank-form");
const lyricsAuthForm = document.getElementById("dank-submit-auth-form");
const lyricsForm = document.getElementById("dank-submit-form");
const dankLyricsOutput = document.getElementById("dank-lyrics-content");
const lyricsSubmitOutput = document.getElementById(
  "dank-lyrics-submit-content",
);
const authOutput = document.getElementById("dank-lyrics-auth-content");
const apiEndpointTab = document.getElementById("api-endpoint");

const loadingMsgs = [
  "Wait",
  "Loading",
  "Finding lyrics",
  "Doing the thing",
  "Crunching the numbers",
  "Loading terrain",
  "Just wait",
  "AAAAAAAAAAAAAAAAAAA",
  "Ring dingdong",
  "Ring ding ding ding dingdong",
];

class DankFormData {
  song_name;
  album_name;
  artist_name;
  constructor(data) {
    this.song_name = data.songName.value;
    this.album_name = data.albumName.value;
    this.artist_name = data.artistName.value;
  }
}

class LyricsFormData extends DankFormData {
  plain_lyrics;
  synced_lyrics;
  constructor(data) {
    super(data);
    this.plain_lyrics = data.lyrics.value;
    // this.syncedLyrics = data.syncedLyrics.value;
  }
}

class LyricsAuthFormData {
  email;
  constructor(data) {
    this.email = data.email.value;
  }
}

function getCookie(key) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${key}=`);
  if (parts.length === 2) return parts.pop().split(";").shift();
}

async function findLyics(dankFormData) {
  let apiLink = `https://api.danklyrics.com/lyrics?`;
  let dankApiLink = `https://api.danklyrics.com/dank/lyrics?`;
  let searchQuery = "";
  searchQuery += `song=${encodeURIComponent(dankFormData.song_name)}`;
  if (dankFormData.artist_name) {
    searchQuery += `&artist=${encodeURIComponent(dankFormData.artist_name)}`;
  }
  if (dankFormData.album_name) {
    searchQuery += `&album=${encodeURIComponent(dankFormData.album_name)}`;
  }
  apiEndpointTab.innerHTML = `<ul><li>Dank API: <code style="user-select: all;">${dankApiLink + searchQuery}</code></li>`;
  apiEndpointTab.innerHTML += `<li>Providers API: <code style="user-select: all;">${apiLink + searchQuery}&providers=dank&providers=lrc</code></li></ul>`;

  const intervalId = setInterval(() => {
    dankLyricsOutput.innerText =
      loadingMsgs[Math.floor(Math.random() * loadingMsgs.length)] + "...";
  }, 1200);

  await fetch(
    `/api/lyrics?song=${dankFormData.song_name}&artist=${dankFormData.artist_name}&album=${dankFormData.album_name}`,
  )
    .then((res) => res.text())
    .then((data) => {
      dankLyricsOutput.innerText = data;
      let urlParams = {
        song: dankFormData.song_name,
      };
      if (dankFormData.artist_name) {
        urlParams.artist = dankFormData.artist_name;
      }
      if (dankFormData.album_name) {
        urlParams.album = dankFormData.album_name;
      }
      window.history.replaceState(
        null,
        document.title,
        `${window.location.protocol}//${window.location.host}?${new URLSearchParams(urlParams)}`,
      );
      clearInterval(intervalId);
    })
    .catch((err) => {
      dankLyricsOutput.innerText = `Something went wrong, ${err}`;
      clearInterval(intervalId);
    });
}

async function handleFindLyricsFormSubmission(e) {
  e.preventDefault();
  const dankFormData = new DankFormData(form);

  return await findLyics(dankFormData);
}

async function handleSubmitLyricsFormSubmission(e) {
  e.preventDefault();
  const dankFormData = new LyricsFormData(lyricsForm);

  const intervalId = setInterval(() => {
    lyricsSubmitOutput.innerText =
      loadingMsgs[Math.floor(Math.random() * loadingMsgs.length)] + "...";
  }, 1200);

  await fetch(`/api/lyrics`, {
    method: "POST",
    body: JSON.stringify(dankFormData),
  })
    .then((resp) => {
      if (!resp.ok) {
        lyricsSubmitOutput.innerText = "Something went wrong";
        alert("Something went wrong");
      }
      clearInterval(intervalId);
      lyricsSubmitOutput.innerText =
        "Done, you'll receive an email when the lyrics is approved!";
      lyricsForm.reset();
    })
    .catch((err) => {
      lyricsSubmitOutput.innerText = `Something went wrong, ${err}`;
      clearInterval(intervalId);
    });
}

async function handleAuthenticateLyricsFormSubmission(e) {
  e.preventDefault();
  const dankFormData = new LyricsAuthFormData(lyricsAuthForm);

  const intervalId = setInterval(() => {
    authOutput.innerText =
      loadingMsgs[Math.floor(Math.random() * loadingMsgs.length)] + "...";
  }, 1200);

  await fetch(`/api/auth`, {
    method: "POST",
    body: JSON.stringify({
      email: dankFormData.email,
    }),
  })
    .then((resp) => {
      if (!resp.ok) {
        authOutput.innerText = "Something went wrong";
        alert("Something went wrong");
      }
      clearInterval(intervalId);
      authOutput.innerText = "Done, check your inbox...";
    })
    .catch((err) => {
      authOutput.innerText = `Something went wrong, ${err}`;
      clearInterval(intervalId);
    });
}

function checkAuth() {
  const cookieToken = getCookie("token");
  if (cookieToken) {
    lyricsAuthForm.style.display = "none";
    lyricsForm.style.display = "flex";
  }
}

async function getSongFromURLQuery() {
  const queryParams = new URL(document.location.toString()).searchParams;
  const songName = queryParams.get("song") ?? "";
  const artistName = queryParams.get("artist") ?? "";
  const albumName = queryParams.get("album") ?? "";

  if (!songName) {
    return;
  }
  if (!songName && !artistName && !albumName) {
    return;
  }

  await findLyics({
    song_name: songName,
    artist_name: artistName,
    album_name: albumName,
  });
}

function init() {
  checkAuth();
  form.addEventListener("submit", handleFindLyricsFormSubmission);
  lyricsForm.addEventListener("submit", handleSubmitLyricsFormSubmission);
  lyricsAuthForm.addEventListener(
    "submit",
    handleAuthenticateLyricsFormSubmission,
  );
  getSongFromURLQuery();

  console.log("helloooo 👋");
}

// init();
