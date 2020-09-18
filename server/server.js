const express = require("express");
const multer = require("multer");
const { initWorker, parseResults, preprocessImage } = require("./ocr");

const app = express();

const storage = multer.memoryStorage();
const upload = multer({ storage });

let worker;

app.post("/", async (req, res, next) => {
  if (!worker) {
    try {
      worker = await initWorker();
    } catch (error) {
      console.error(error);
      next(error);
    }
  }
  next();
});

app.post("/", upload.single("image"), async (req, res) => {
  if (!req.file) {
    res.status(400).send("Missing image");
  }

  try {
    const buffer = await preprocessImage(req.file.buffer);
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
