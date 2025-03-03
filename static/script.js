"use strict";

const form = document.getElementById("dank-form");
const dankLyricsOutput = document.getElementById("dank-lyrics-content");
const loadingMsgs = [
  "Wait",
  "Loading",
  "Finding lyrics",
  "Doing the thing",
  "Crunching the numbers",
  "Loading terrain",
  "Just wait",
  "AAAAAAAAAAAAAAAAAAA",
];

form.addEventListener("submit", async (e) => {
  e.preventDefault();
  const dankFormData = new DankFormData(form);

  const intervalId = setInterval(() => {
    dankLyricsOutput.innerText =
      loadingMsgs[Math.floor(Math.random() * loadingMsgs.length)] + "...";
  }, 1200);

  await fetch(
    `/lyrics?song=${dankFormData.songName}&artist=${dankFormData.artistName}&album=${dankFormData.albumName}`,
  )
    .then((res) => res.text())
    .then((data) => {
      dankLyricsOutput.innerText = data;
      clearInterval(intervalId);
    })
    .catch((err) => {
      dankLyricsOutput.innerText = `Something went wrong, ${err}`;
      clearInterval(intervalId);
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
