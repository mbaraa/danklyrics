"use strict";

const form = document.getElementById("dank-form");
const dankLyricsOutput = document.getElementById("dank-lyrics-content");

form.addEventListener("submit", async (e) => {
  e.preventDefault();
  const dankFormData = new DankFormData(form);

  await fetch(
    `/lyrics?song=${dankFormData.songName}&artist=${dankFormData.artistName}&album=${dankFormData.albumName}`,
  )
    .then((res) => res.text())
    .then((data) => {
      dankLyricsOutput.innerText = data;
    })
    .catch((err) => {
      dankLyricsOutput.innerText = `Something went wrong, ${err}`;
    });
});

class DankFormData {
  songName;
  albumName;
  artistName;
  constructor(data) {
    this.songName = data.songName.value;
    this.albumName = data.albumName.value;
    this.artistName = data.artistName.value;
  }
}
