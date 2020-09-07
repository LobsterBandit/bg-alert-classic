const { createWorker } = require("tesseract.js");
const jimp = require("jimp");

const bgTimerRegex = /^(?<bg>[A-Z]{2}):\s(?<ready>READY!!!)?(?<hours>-?[0-9]{2}h)?\s?(?<minutes>-?[0-9]{2}m)?\s?(?<seconds>-?[0-9]{2}s)?$/gm;
const charWhitelist = "ABDERSVWYhms:0123456789 !-";

async function main(args) {
  if (!args[0]) {
    throw new Error("Must supply an input image");
  }

  const worker = await initWorker();

  const buffer = await preprocessImage(args[0]);

  const { data } = await worker.recognize(buffer);

  await worker.terminate();

  return data.text;
}

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

async function preprocessImage(imageName, saveIntermediate = false) {
  const image = await jimp.read(`testdata/${imageName}.png`);
  image.contrast(1).scale(1.5);
  if (saveIntermediate) {
    await image.writeAsync(`testdata/preprocessed${imageName}.png`);
  }
  const buffer = await image.getBufferAsync(jimp.MIME_PNG);
  return buffer;
}

module.exports = { initWorker, parseResults, preprocessImage };

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
