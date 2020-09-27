#!/usr/bin/env node

const { app, shutdown } = require("./server");

const port = process.env.PORT || 3003;

const server = app.listen(port, () => {
  console.log(`bg-alert-classic listening at http://localhost:${port}`);
});

// quit on ctrl-c when running docker in terminal
process.on("SIGINT", function onSigint() {
  console.info(
    "Got SIGINT (aka ctrl-c in docker). Graceful shutdown ",
    new Date().toISOString()
  );
  shutdown(server);
});

// quit properly on docker stop
process.on("SIGTERM", function onSigterm() {
  console.info(
    "Got SIGTERM (docker container stop). Graceful shutdown ",
    new Date().toISOString()
  );
  shutdown(server);
});
