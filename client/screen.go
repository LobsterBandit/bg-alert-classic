package client

import (
	"image"

	"github.com/kbinani/screenshot"
)

func CaptureScreenArea(area image.Rectangle) *image.RGBA {
	img, err := screenshot.CaptureRect(area)
	if err != nil {
		panic(err)
	}

	return img
}

func GetScreenBounds(display int) image.Rectangle {
	return screenshot.GetDisplayBounds(display)
}

func PrimaryScreenBounds() image.Rectangle {
	return screenshot.GetDisplayBounds(0)
}

func NumActiveDisplays() int {
	return screenshot.NumActiveDisplays()
}
