# wowclassic-bg-ocr client

A Windows Go client to screenshot an area of the screen and save a png or upload to an OCR server for analysis.

## Development

```bash
# mount source into a dev container with go and supporting tools installed
docker run --rm -it --mount type=bind,source="$(pwd)",targer=/usr/src lobsterbandit/dev-golang:edge
```

## Build

```bash
# build with windows/amd64 target
# screenshot code is os specific
GOOS=windows GOARCH=amd64 go build
```
