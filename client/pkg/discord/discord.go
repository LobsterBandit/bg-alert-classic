// Package discord provides types and functions for sending
// webhooks to discord
package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/lobsterbandit/wowclassic-bg-ocr/client"
)

const (
	webhookID    string = "758496757217755167"
	webhookToken string = "CWZS8f3gABLgyq2XMWqmq5F3-I4hgW37CHvzeX_nk8nf0Z4ldfxbrG9HqZYm2DrPuERF"
)

var discordWebhookURL string

func init() {
	discordWebhookURL = fmt.Sprintf("https://discordapp.com/api/webhooks/%s/%s", webhookID, webhookToken)
}

type WebhookParams struct {
	Content  string           `json:"content,omitempty"`
	Username string           `json:"username,omitempty"`
	Images   []WebhookImage   `json:"images,omitempty"`
	Embeds   []*MessageEmbed  `json:"embeds,omitempty"`
	Timers   []client.BgTimer `json:"timers,omitempty"`
}

type WebhookImage struct {
	Name  string
	Image *image.RGBA
}

func PostDiscordMessage(images []WebhookImage, timers []client.BgTimer) (err error) {
	fmt.Println("\nSending webhook to discord...")

	webhookParams := &WebhookParams{
		Content: "BG Timer Message",
		Images:  images,
		Timers:  timers,
	}

	// fmt.Printf("Webhook payload: %v\n", webhookParams)

	msg, err := executeWebhookMultipart(false, webhookParams)
	if err != nil {
		return
	}

	response, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return
	}

	fmt.Printf("Webhook response: %s\n", string(response))

	return
}

func executeWebhookMultipart(wait bool, data *WebhookParams) (response *Message, err error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	// add image form fields
	err = addImageParts(w, data.Images)
	if err != nil {
		return
	}

	// add other content in form field "payload_json"
	if data.Content != "" {
		err = addPayloadJSON(w, data.Content)
		if err != nil {
			return
		}
	}

	w.Close()

	// req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, discordWebhookURL, body)
	// if err != nil {
	// 	return
	// }

	// if body != nil {
	// 	req.Header.Set("Content-Type", w.FormDataContentType())
	// }

	// resp, err := http.DefaultClient.Do(req)

	url := discordWebhookURL
	if wait {
		url += "?wait=true"
	}

	fmt.Printf("Issuing webhook to %s\n", url)

	resp, err := http.Post(url, w.FormDataContentType(), body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Println("discord posting response: ", resp.Status)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return
	}

	return response, nil
}

func addImageParts(w *multipart.Writer, images []WebhookImage) error {
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

func addPayloadJSON(w *multipart.Writer, content string) error {
	jsonPayload, err := json.Marshal(map[string]string{
		"content": content,
	})
	if err != nil {
		return err
	}
	// var tmp struct {
	// 	Content string `json:"content"`
	// }

	// tmp.Content = content
	// fmt.Println(string(jsonPayload))

	fw, err := w.CreateFormField("payload_json")
	if err != nil {
		return err
	}

	_, err = fw.Write(jsonPayload)
	if err != nil {
		return err
	}

	return nil
}
