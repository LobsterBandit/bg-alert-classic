const app = require("./server");

const port = process.env.PORT || 3000;

app.listen(port, () => {
  console.log(`wowclassic-bg-ocr listening at http://localhost:${port}`);
});
