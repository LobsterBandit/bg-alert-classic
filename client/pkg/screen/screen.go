// Package screen provides access to display devices
package screen

import (
	"fmt"
	"image"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/lobsterbandit/wowclassic-bg-ocr/client/pkg/img"
)

// Capture the given screen area
func Capture(area image.Rectangle) (*img.File, error) {
	capture, err := captureArea(area)
	if err != nil {
		return nil, err
	}

	captureTime := time.Now()
	fileName := fmt.Sprintf("%dx%d_%s_%s_%d.png",
		capture.Rect.Dx(), capture.Rect.Dy(), area.Min, area.Max, captureTime.Unix())

	fmt.Printf("Captured screen area: %v\n\tTimestamp: %s\n\tFilename: %q\n\n", area, captureTime, fileName)

	return &img.File{
		Name:      fileName,
		Timestamp: img.Timestamp(captureTime.Unix()),
		Image:     capture,
	}, nil
}

// CapturePrimary is a helper func to capture the entire primary display
func CapturePrimary() (*img.File, error) {
	return Capture(PrimaryScreenBounds())
}

func captureArea(area image.Rectangle) (*image.RGBA, error) {
	img, err := screenshot.CaptureRect(area)
	if err != nil {
		return nil, fmt.Errorf("screen capture error: %w", err)
	}

	return img, nil
}

func GetBounds(display int) image.Rectangle {
	return screenshot.GetDisplayBounds(display)
}

func PrimaryScreenBounds() image.Rectangle {
	return GetBounds(0)
}

func NumActiveDisplays() int {
	return screenshot.NumActiveDisplays()
}
