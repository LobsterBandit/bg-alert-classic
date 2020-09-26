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

func PostImage(server string, image *image.RGBA, fileName string) ([]BgTimer, error) {
	fmt.Printf("Posting image %q to OCR server %s\n", fileName, server)

	buf, contentType, err := createMultipartImage(image, fileName)
	if err != nil {
		return nil, fmt.Errorf("image posting error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, server, buf)
	if err != nil {
		return nil, fmt.Errorf("image posting error: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("image posting error: %w", err)
	}
	defer resp.Body.Close()

	var timers []BgTimer

	err = json.NewDecoder(resp.Body).Decode(&timers)
	if err != nil {
		return nil, fmt.Errorf("image posting error: %w", err)
	}

	return timers, err
}

func SaveImage(image *image.RGBA, fileName string) error {
	file, _ := os.Create(fileName)
	defer file.Close()

	err := png.Encode(file, image)
	if err != nil {
		return fmt.Errorf("image saving error: %w", err)
	}

	return nil
}

func createMultipartImage(image *image.RGBA, fileName string) (*bytes.Buffer, string, error) {
	buf := new(bytes.Buffer)

	w := multipart.NewWriter(buf)
	defer w.Close()

	fw, err := w.CreateFormFile("image", fileName)
	if err != nil {
		return buf, "", fmt.Errorf("create multipart error: %w", err)
	}

	err = png.Encode(fw, image)
	if err != nil {
		return buf, "", fmt.Errorf("create multipart error: %w", err)
	}

	return buf, w.FormDataContentType(), nil
}
