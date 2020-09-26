# classic-bg-ocr

A Windows Go client to screenshot an area of the screen and
- save as png
- run OCR analysis for battleground timers
- post message to discord channel via webhook

## Development

```bash
# mount source into a dev container with go and supporting tools installed
docker run --rm -it --mount type=bind,source="$(pwd)",targer=/usr/src lobsterbandit/dev-golang:edge
```

## Build

```bash
# build with windows/amd64 target
# screenshot code is os specific
GOOS=windows GOARCH=amd64 go build -trimpath -o bin cmd/classic-bg-ocr
```
