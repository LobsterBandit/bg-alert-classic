// +build windows

package main

import (
	"fmt"
	"image"
	"time"

	"github.com/lobsterbandit/wowclassic-bg-ocr/client"
)

func main() {
	version := "v1.0.0-alpha1"
	upload := true
	save := true

	fmt.Printf("\nwowclassic-bg-ocr-client %v\n\tWoW Classic BG timer screen capture and analysis\n\n", version)

	pt0 := image.Point{
		X: 575,
		Y: 1275,
	}
	pt1 := image.Point{
		X: 800,
		Y: 1400,
	}
	bgArea := image.Rectangle{pt0, pt1}

	img := client.CaptureScreenArea(bgArea)

	captureTime := time.Now()
	fileName := fmt.Sprintf("%dx%d_%s_%s_%d.png",
		img.Rect.Dx(), img.Rect.Dy(), bgArea.Min, bgArea.Max, captureTime.Unix())

	fmt.Printf("Captured screen area: %v\n\tTimestamp: %s\n\tFilename: %q\n\n", bgArea, captureTime, fileName)

	if upload {
		timers := client.PostImage("http://192.168.1.14:3003", img, fileName)
		fmt.Printf("\nTimer Results:\n%v\n", timers)
	}

	if save {
		client.SaveImage(img, fileName)
	}
}
