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
  console.log(req.file);
  console.log(req.body);
  console.log(
    new Date(parseInt(req.body["timestamp"], 10) * 1000).toLocaleString()
  );

  try {
    const buffer = await preprocessImage(
      req.file.buffer,
      process.env.SAVE_PREPROCESSED,
      req.file.originalname
    );
    const {
      data: { text },
    } = await worker.recognize(buffer);

    const results = parseResults(text);

    res.send(
      results.map((result) => {
        result["timestamp"] = parseInt(req.body["timestamp"], 10);
        return result;
      })
    );
  } catch (error) {
    console.error(error);
    res.status(500).send(error);
  }
});

async function shutdown(server) {
  if (worker) {
    await worker.terminate();
  }
  server.close((err) => {
    if (err) {
      console.error(err);
      process.exitCode = 1;
    }
    console.log("HTTP server closed");
    process.exit();
  });
}

module.exports = { app, shutdown };
