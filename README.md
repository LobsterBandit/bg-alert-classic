# wowclassic-bg-ocr

A client/server tool to continuously capture a specific area of the screen containing BG timers for OCR analysis and alerting.

## [client](client/README.md)

Go client to capture screenshots and upload to the server for analysis

```bash
# starts both server and client containers in development
# specify server|client service name to control individually
docker-compose -f docker-compose.dev.yml up -d --build

# attach terminal to container
# or attach via vs code remote container to install extensions
docker-compose -f docker-compose.dev.yml exec client bash

# build executable inside container
# client.exe output to client folder
./build_client_windows.sh
```

## [server](server/README.md)

node.js express server providing the OCR capabilities via a single POST endpoint

```bash
# starts both server and client containers in development
# specify server|client service name to control individually
docker-compose -f docker-compose.dev.yml up -d --build

# run production build of ocr server
docker-compose up -d --build
```
