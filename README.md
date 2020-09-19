# wowclassic-bg-ocr

A client/server tool to continuously capture a specific area of the screen containing BG timers for OCR analysis and alerting.

## [client](client/README.md)

Go client to capture screenshots and upload to the server for analysis

## [server](server/README.md)

node.js express server providing the OCR capabilities via a single POST endpoint

```bash
# run development build of ocr server with live reload
docker-compose -f docker-compose.dev.yml up -d --build

# run production build of ocr server
docker-compose up -d --build
```
