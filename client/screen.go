// +build windows

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
