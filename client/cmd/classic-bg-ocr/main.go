package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/lobsterbandit/wowclassic-bg-ocr/client"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Fatal error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	version := "v1.0.0-alpha1"
	upload := true
	discord := true
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

	img, err := client.CaptureScreenArea(bgArea)

	captureTime := time.Now()
	fileName := fmt.Sprintf("%dx%d_%s_%s_%d.png",
		img.Rect.Dx(), img.Rect.Dy(), bgArea.Min, bgArea.Max, captureTime.Unix())

	fmt.Printf("Captured screen area: %v\n\tTimestamp: %s\n\tFilename: %q\n\n", bgArea, captureTime, fileName)

	var timers []client.BgTimer
	if upload {
		timers, err = client.PostImage("http://192.168.1.14:3003", img, fileName)
		if err != nil {
			return err
		}

		fmt.Printf("\nTimer Results:\n%v\n", timers)
	}

	if discord {
		// webhook to post discord channel message
		err = client.PostDiscordMessage([]client.WebhookImage{
			{Name: fileName, Image: img},
		}, timers)
		if err != nil {
			return err
		}
	}

	if save {
		err = client.SaveImage(img, fileName)
	}

	fmt.Println("\nComplete!")

	return err
}
