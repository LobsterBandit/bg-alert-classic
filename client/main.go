// +build windows

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/kbinani/screenshot"
)

type BgTimer struct {
	Bg      string
	Hours   string
	Minutes string
	Seconds string
	Ready   bool
}

func main() {
	version := "v1.0.0-alpha1"
	upload := true
	save := true

	fmt.Printf("\nwowclassic-bg-ocr-client %v\n\tWoW Classic BG timer screen capture and analysis\n\n", version)

	bgArea, img := captureBgArea()

	captureTime := time.Now()
	fileName := fmt.Sprintf("%dx%d_%s_%s_%d.png",
		img.Rect.Dx(), img.Rect.Dy(), bgArea.Min, bgArea.Max, captureTime.Unix())

	fmt.Printf("Captured screen area: %v\n\tTimestamp: %s\n\tFilename: %q\n\n", bgArea, captureTime, fileName)

	if upload {
		timers := postImage("http://192.168.1.14:3003", img, fileName)
		fmt.Printf("\nTimer Results:\n%v\n", timers)
	}

	if save {
		saveImage(img, fileName)
	}
}

func captureBgArea() (image.Rectangle, *image.RGBA) {
	pt0 := image.Point{
		X: 575,
		Y: 1275,
	}
	pt1 := image.Point{
		X: 800,
		Y: 1400,
	}
	bgArea := image.Rectangle{pt0, pt1}

	img, err := screenshot.CaptureRect(bgArea)
	if err != nil {
		panic(err)
	}

	return bgArea, img
}

func postImage(server string, image *image.RGBA, fileName string) []BgTimer {
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

func saveImage(image *image.RGBA, fileName string) {
	file, _ := os.Create(fileName)
	defer file.Close()

	err := png.Encode(file, image)
	if err != nil {
		panic(err)
	}
}
