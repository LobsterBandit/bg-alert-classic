// +build windows

package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/kbinani/screenshot"
)

func main() {
	upload := true
	save := true

	bgArea, img := captureBgArea()

	fileName := fmt.Sprintf("%dx%d_%s_%s_%d.png",
		img.Rect.Dx(), img.Rect.Dy(), bgArea.Min, bgArea.Max, time.Now().Unix())

	fmt.Printf("%v \"%s\"\n", bgArea, fileName)

	if upload {
		postImage(img, fileName)
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

func postImage(image *image.RGBA, fileName string) {
	fmt.Printf("Posting image %s to OCR server", fileName)
}

func saveImage(image *image.RGBA, fileName string) {
	file, _ := os.Create(fileName)
	defer file.Close()

	err := png.Encode(file, image)
	if err != nil {
		panic(err)
	}
}
