package client

import (
	"fmt"
	"image"

	"github.com/kbinani/screenshot"
)

func CaptureScreenArea(area image.Rectangle) (*image.RGBA, error) {
	img, err := screenshot.CaptureRect(area)
	if err != nil {
		err = fmt.Errorf("screen capture error: %w", err)
	}

	return img, err
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
