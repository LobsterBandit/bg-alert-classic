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
	"path/filepath"
	"strconv"
)

type Timestamp int64

func (t Timestamp) ToInt64() int64 {
	return int64(t)
}

func (t Timestamp) ToString() string {
	return strconv.FormatInt(t.ToInt64(), 10)
}

type BgTimer struct {
	Timestamp Timestamp `json:"timestamp"`
	Bg        string    `json:"bg"`
	Hours     string    `json:"hours"`
	Minutes   string    `json:"minutes"`
	Seconds   string    `json:"seconds"`
	Ready     bool      `json:"ready"`
}

type File struct {
	Name      string
	Timestamp Timestamp
	Image     *image.RGBA
}

func (f *File) Post(server string) (results []BgTimer, err error) {
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

	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return results, err
}

func (f *File) Save() error {
	file, err := os.Create(f.Name)
	if err != nil {
		return fmt.Errorf("image saving error: %w", err)
	}
	defer file.Close()

	err = png.Encode(file, f.Image)
	if err != nil {
		return fmt.Errorf("image saving error: %w", err)
	}

	return nil
}

func (f *File) Write(path string) error {
	filePath := ""

	if ext := filepath.Ext(path); ext != "" {
		// extension included, assume full filename given
		filePath = filepath.FromSlash(path)
	} else {
		// only dir given use original filename
		filePath = filepath.Join(path, f.Name)
	}

	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, f.Image)
	if err != nil {
		return fmt.Errorf("image saving error: %w", err)
	}

	return nil
}

func AddToMultipartForm(w *multipart.Writer, images []*File) error {
	for _, image := range images {
		fw, err := w.CreateFormFile("image", image.Name)
		if err != nil {
			return err
		}

		err = png.Encode(fw, image.Image)
		if err != nil {
			return err
		}

		fw, err = w.CreateFormField("timestamp")
		if err != nil {
			return err
		}

		_, err = fw.Write([]byte(image.Timestamp.ToString()))
		if err != nil {
			return err
		}
	}

	return nil
}
