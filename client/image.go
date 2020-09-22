package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
)

type BgTimer struct {
	Bg      string
	Hours   string
	Minutes string
	Seconds string
	Ready   bool
}

func PostImage(server string, image *image.RGBA, fileName string) []BgTimer {
	fmt.Printf("Posting image %q to OCR server %s\n", fileName, server)

	buf := new(bytes.Buffer)

	w := multipart.NewWriter(buf)

	fw, err := w.CreateFormFile("image", fileName)
	if err != nil {
		panic(err)
	}

	err = png.Encode(fw, image)
	if err != nil {
		panic(err)
	}

	w.Close()

	req, err := http.NewRequest(http.MethodPost, server, buf)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var timers []BgTimer

	err = json.NewDecoder(resp.Body).Decode(&timers)
	if err != nil {
		panic(err)
	}

	return timers
}

func SaveImage(image *image.RGBA, fileName string) {
	file, _ := os.Create(fileName)
	defer file.Close()

	err := png.Encode(file, image)
	if err != nil {
		panic(err)
	}
}
