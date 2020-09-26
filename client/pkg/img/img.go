// Package img provides the ability to work with images
package img

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type BgTimer struct {
	Bg      string
	Hours   string
	Minutes string
	Seconds string
	Ready   bool
}

type File struct {
	Name  string
	Image *image.RGBA
}

func (f *File) Post(server string) (timers []BgTimer, err error) {
	fmt.Printf("Posting image %q to OCR server %s\n", f.Name, server)

	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	err = AddToMultipartForm(w, []*File{f})
	if err != nil {
		return nil, err
	}

	w.Close()

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, server, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&timers)
	if err != nil {
		return nil, err
	}

	return timers, err
}

func (f *File) Save() error {
	file, _ := os.Create(f.Name)
	defer file.Close()

	err := png.Encode(file, f.Image)
	if err != nil {
		return fmt.Errorf("image saving error: %w", err)
	}

	return nil
}

func AddToMultipartForm(w *multipart.Writer, images []*File) error {
	for i, image := range images {
		fw, err := w.CreateFormFile("image"+strconv.Itoa(i), image.Name)
		if err != nil {
			return err
		}

		err = png.Encode(fw, image.Image)
		if err != nil {
			return err
		}
	}

	return nil
}
