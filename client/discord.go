package client

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
)

type EmbedType string

const (
	webhookID    string = "758496757217755167"
	webhookToken string = "CWZS8f3gABLgyq2XMWqmq5F3-I4hgW37CHvzeX_nk8nf0Z4ldfxbrG9HqZYm2DrPuERF"

	EmbedTypeRich    EmbedType = "rich"
	EmbedTypeImage   EmbedType = "image"
	EmbedTypeVideo   EmbedType = "video"
	EmbedTypeGifv    EmbedType = "gifv"
	EmbedTypeArticle EmbedType = "article"
	EmbedTypeLink    EmbedType = "link"
)

var discordWebhookURL string

func init() {
	discordWebhookURL = fmt.Sprintf("https://discordapp.com/api/webhooks/%s/%s", webhookID, webhookToken)
}

type WebhookParams struct {
	Content  string          `json:"content,omitempty"`
	Username string          `json:"username,omitempty"`
	Images   []WebhookImage  `json:"images,omitempty"`
	Embeds   []*MessageEmbed `json:"embeds,omitempty"`
	Timers   []BgTimer       `json:"timers,omitempty"`
}

type MessageEmbed struct {
	Type        string `json:"type,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
	Color       int    `json:"color,omitempty"`
}

type Message struct {
	Body string
}

type WebhookImage struct {
	Name  string
	Image *image.RGBA
}

func PostDiscordMessage(images []WebhookImage, timers []BgTimer) (err error) {
	fmt.Printf("\nSending webhook to %s\n", discordWebhookURL)

	webhookParams := &WebhookParams{
		Content: "BG Timer Message",
		Images:  images,
		Timers:  timers,
	}

	// fmt.Printf("Webhook payload: %v\n", webhookParams)

	msg, err := executeWebhookMultipart(true, webhookParams)
	if err != nil {
		return
	}

	fmt.Printf("Webhook response: %v\n", msg)

	return
}

func executeWebhookMultipart(wait bool, data *WebhookParams) (st *Message, err error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	// add image form fields
	addImageParts(w, data.Images)

	// add other content in form field "payload_json"
	if data.Content != "" {
		addPayloadJson(w, data.Content)
	}

	w.Close()

	req, err := http.NewRequest(http.MethodPost, discordWebhookURL, body)
	if err != nil {
		return
	}

	if body != nil {
		req.Header.Set("Content-Type", w.FormDataContentType())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Println("discord posting response: ", resp.Status)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// fmt.Println("discord posting response: ", string(respBody))

	return &Message{Body: string(respBody)}, nil
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

func addPayloadJson(w *multipart.Writer, content string) error {
	var tmp struct {
		Content string `json:"content"`
	}
	tmp.Content = content

	jsonPayload, err := json.Marshal(&tmp)
	if err != nil {
		return err
	}
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
