const express = require("express");
const { initWorker, parseResults, preprocessImage } = require("./ocr");

const app = express();

let worker;

app.get("/ocr", async (req, res, next) => {
  if (!worker) {
    try {
      worker = await initWorker();
    } catch (error) {
      console.error(error);
    }
  }
  next();
});

app.get("/", (req, res) => {
  res.send("Hello World!");
});

app.get("/ocr", async (req, res) => {
  console.log(req.query);
  if (!req.query.imageName) {
    res.status(400).send("Bad Request");
  }

  try {
    const buffer = await preprocessImage(req.query.imageName);
    const {
      data: { text },
    } = await worker.recognize(buffer);

    const results = parseResults(text);

    res.send(results);
  } catch (error) {
    console.error(error);
    res.status(500).send(error);
  }
});

// async function shutdownHandler() {
//   console.log("SIGTERM signal received: closing HTTP server");
//   await worker.terminate();
//   app.close(() => {
//     console.log("HTTP server closed");
//   });
// }

// process.on("SIGTERM", shutdownHandler);
// process.on("SIGINT", shutdownHandler);

module.exports = app;
