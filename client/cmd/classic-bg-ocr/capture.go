package main

import (
	"errors"
	"fmt"
	"image"
	"time"

	d "github.com/lobsterbandit/wowclassic-bg-ocr/client/pkg/discord"
	"github.com/lobsterbandit/wowclassic-bg-ocr/client/pkg/img"
	"github.com/lobsterbandit/wowclassic-bg-ocr/client/pkg/screen"
	"github.com/spf13/cobra"
)

var (
	analyze      bool
	discord      bool
	file         string
	save         bool
	webhookID    string
	webhookToken string

	captureCmd = &cobra.Command{
		Use:   "capture",
		Short: "capture the screen",
		Long:  "capture the full screen or a given rectangle",
		RunE:  runCapture,
	}

	ErrWebhookMissingRequiredArgs = errors.New("webhookID and webhookToken are required")
)

func init() {
	// TODO: add ocr server flag for when analyze is set
	captureCmd.Flags().BoolVarP(&analyze, "analyze", "a", false, "analyze image for bg timers")
	captureCmd.Flags().BoolVarP(
		&discord, "discord", "d", false, "send screen capture and analysis via webhook to a discord channel")
	// TODO: add output flag for file save location
	captureCmd.Flags().StringVarP(&file, "file", "f", "", "path to a local image file")
	captureCmd.Flags().BoolVarP(&save, "save", "s", false, "save captured image to file")
	captureCmd.Flags().StringVarP(&webhookID, "id", "i", "", "discord webhook id")
	captureCmd.Flags().StringVarP(&webhookToken, "token", "t", "", "discord webhook token")

	rootCmd.AddCommand(captureCmd)
}

func runCapture(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nwowclassic-bg-ocr-client %v\n\tWoW Classic BG timer screen capture and analysis\n\n", version)

	// exit early if required arg combinations are not met
	if discord && (webhookID == "" || webhookToken == "") {
		return fmt.Errorf("missing required arguments to send discord webhooks: %w", ErrWebhookMissingRequiredArgs)
	}

	// get points from flags
	// otherwise default to primary screen bounds
	pt0 := image.Point{
		X: 575,
		Y: 1275,
	}
	pt1 := image.Point{
		X: 800,
		Y: 1400,
	}
	bgArea := image.Rectangle{pt0, pt1}

	capture, err := screen.CaptureScreenArea(bgArea)

	captureTime := time.Now()
	fileName := fmt.Sprintf("%dx%d_%s_%s_%d.png",
		capture.Rect.Dx(), capture.Rect.Dy(), bgArea.Min, bgArea.Max, captureTime.Unix())

	fmt.Printf("Captured screen area: %v\n\tTimestamp: %s\n\tFilename: %q\n\n", bgArea, captureTime, fileName)

	imageFile := &img.File{
		Name:  fileName,
		Image: capture,
	}

	var timers []img.BgTimer
	if analyze {
		timers, err = imageFile.Post("http://192.168.1.14:3003")
		if err != nil {
			return err
		}

		fmt.Printf("\nTimer Results:\n%v\n", timers)
	}

	if discord {
		webhook := &d.Webhook{
			ID:    webhookID,
			Token: webhookToken,
			Params: &d.WebhookParams{
				Content: "BG Timer Alert",
				Images:  []*img.File{imageFile},
				Timers:  timers,
			},
		}

		// webhook to post discord channel message
		err = webhook.PostDiscordMessage()
		if err != nil {
			return err
		}
	}

	if save {
		err = imageFile.Save()
	}

	fmt.Println("\nComplete!")

	return err
}
