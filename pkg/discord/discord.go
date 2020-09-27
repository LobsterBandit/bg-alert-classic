// Package discord provides types and functions for sending
// webhooks to discord
package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/lobsterbandit/bg-alert-classic/pkg/img"
)

var webhookBaseURL string = "https://discordapp.com/api/webhooks/"

type Webhook struct {
	ID     string
	Token  string
	Params *WebhookParams
}

type WebhookParams struct {
	Content  string          `json:"content,omitempty"`
	Username string          `json:"username,omitempty"`
	Images   []*img.File     `json:"images,omitempty"`
	Embeds   []*MessageEmbed `json:"embeds,omitempty"`
	Timers   []img.BgTimer   `json:"timers,omitempty"`
}

func (w *Webhook) PostDiscordMessage() (err error) {
	fmt.Println("\nSending webhook to discord...")

	msg, err := w.executeMultipart(false)
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

func (w *Webhook) URL() string {
	return webhookBaseURL + w.ID + "/" + w.Token
}

func (w *Webhook) executeMultipart(wait bool) (response *Message, err error) {
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)

	// add image form fields
	err = img.AddToMultipartForm(mw, w.Params.Images)
	if err != nil {
		return
	}

	// add other content in form field "payload_json"
	if w.Params.Content != "" {
		err = addPayloadJSON(mw, w.Params.Content)
		if err != nil {
			return
		}
	}

	mw.Close()

	url := w.URL()
	if wait {
		url += "?wait=true"
	}

	fmt.Printf("Issuing webhook to %s\n", url)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, body)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", mw.FormDataContentType())

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

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return
	}

	return response, nil
}

func addPayloadJSON(w *multipart.Writer, content string) error {
	jsonPayload, err := json.Marshal(map[string]string{
		"content": content,
	})
	if err != nil {
		return err
	}

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
