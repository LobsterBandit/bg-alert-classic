const fs = require("fs/promises");
const { createWorker } = require("tesseract.js");
const jimp = require("jimp");

const bgTimerRegex = /^(?<bg>[ABSWV]{2}):\s(?<ready>READY!!!)?(?<hours>-?[0-9]{2}h)?\s?(?<minutes>-?[0-9]{2}m)?\s?(?<seconds>-?[0-9]{2}s)?$/gm;
const charWhitelist = "ABDERSVWYhms:0123456789 !-";

function parseResults(text) {
  console.log(text);

  return Array.from(text.matchAll(bgTimerRegex), (match) => ({
    bg: match.groups.bg,
    hours: match.groups.hours && match.groups.hours.replace("h", ""),
    minutes: match.groups.minutes && match.groups.minutes.replace("m", ""),
    seconds: match.groups.seconds && match.groups.seconds.replace("s", ""),
    ready: !!match.groups.ready,
  }));
}

async function initWorker() {
  const worker = createWorker({
    logger: (m) => console.log(m),
    errorHandler: (m) => console.log(m),
  });

  await worker.load();
  await worker.loadLanguage("eng");
  await worker.initialize("eng");
  await worker.setParameters({
    tessedit_char_whitelist: charWhitelist,
  });

  return worker;
}

async function preprocessImage(
  inputBuffer,
  saveIntermediate = false,
  imageName = `preprocessed-${Date.now()}`
) {
  const image = await jimp.read(inputBuffer);
  image.contrast(1).scale(1.5);
  if (saveIntermediate) {
    await image.writeAsync(
      imageName.endsWith(".png")
        ? `testdata/${imageName}`
        : `testdata/${imageName}.png`
    );
  }
  const buffer = await image.getBufferAsync(jimp.MIME_PNG);
  return buffer;
}

module.exports = { initWorker, parseResults, preprocessImage };

async function main(args) {
  if (!args[0]) {
    throw new Error("Must supply an input image");
  }

  const worker = await initWorker();

  const inputBuffer = await fs.readFile(args[0]);
  const buffer = await preprocessImage(inputBuffer, true);

  const { data } = await worker.recognize(buffer);

  await worker.terminate();

  return data.text;
}

if (!module.parent) {
  main(process.argv.slice(2))
    .then(parseResults)
    .then((timers) => {
      console.log(timers);
    })
    .catch((e) => {
      console.log(e);
    });
}
