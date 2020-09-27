package main

import (
	"errors"
	"fmt"
	"image"

	d "github.com/lobsterbandit/wowclassic-bg-ocr/client/pkg/discord"
	"github.com/lobsterbandit/wowclassic-bg-ocr/client/pkg/img"
	"github.com/lobsterbandit/wowclassic-bg-ocr/client/pkg/screen"
	"github.com/spf13/cobra"
)

var (
	// send image for OCR analysis
	analyze bool
	// OCR analysis endpoint
	ocrURL string
	// send image results as message to discord channel
	discord bool
	// discord webhook id
	webhookID string
	// discord webhook token
	webhookToken string
	// path to local image
	file string
	// write captured image to file
	write bool
	// path to write captured image
	outFile string

	captureCmd = &cobra.Command{
		Use:   "capture",
		Short: "capture the screen",
		Long:  "capture the full screen or a given rectangle",
		RunE:  runCapture,
	}

	ErrCaptureIncompatibleFlagset = errors.New("incompatible set of flags")
	ErrWebhookMissingRequiredArgs = errors.New("webhookID and webhookToken are required")
)

func init() {
	captureCmd.Flags().BoolVarP(&analyze, "analyze", "a", false, "analyze image for bg timers")
	captureCmd.Flags().StringVarP(&ocrURL, "url", "u", "", "remote ocr analysis endpoint, required if analyze is set")

	captureCmd.Flags().BoolVarP(
		&discord, "discord", "d", false, "send screen capture and analysis via webhook to a discord channel")
	captureCmd.Flags().StringVarP(&webhookID, "id", "i", "", "discord webhook id")
	captureCmd.Flags().StringVarP(&webhookToken, "token", "t", "", "discord webhook token")

	captureCmd.Flags().BoolVarP(&write, "write", "w", false, "write captured image to file")
	captureCmd.Flags().StringVarP(&outFile, "out", "o", "", "path to write captured image")

	captureCmd.Flags().StringVarP(&file, "file", "f", "", "path to a local image file")

	rootCmd.AddCommand(captureCmd)
}

func runCapture(cmd *cobra.Command, args []string) (err error) {
	// exit early if required arg combinations are not met
	if discord && (webhookID == "" || webhookToken == "") {
		return fmt.Errorf("missing required arguments to send discord webhooks: %w", ErrWebhookMissingRequiredArgs)
	}

	if analyze && ocrURL == "" {
		return fmt.Errorf("url is required if analyze is set: %w", ErrCaptureIncompatibleFlagset)
	}

	var imageFile *img.File
	if file != "" {
		imageFile, err = img.FromFile(file)
	} else {
		// TODO: get points from flags
		pt0 := image.Point{
			X: 575,
			Y: 1275,
		}
		pt1 := image.Point{
			X: 800,
			Y: 1400,
		}
		imageFile, err = screen.Capture(image.Rectangle{pt0, pt1})
	}

	var results []img.BgTimer
	if analyze {
		results, err = imageFile.Post(ocrURL)
		if err != nil {
			return err
		}

		fmt.Printf("\nTimer Results:\n%v\n", results)
	}

	if discord {
		webhook := &d.Webhook{
			ID:    webhookID,
			Token: webhookToken,
			Params: &d.WebhookParams{
				Content: "BG Timer Alert",
				Images:  []*img.File{imageFile},
				Timers:  results,
			},
		}

		// webhook to post discord channel message
		err = webhook.PostDiscordMessage()
		if err != nil {
			return err
		}
	}

	if write {
		if outFile != "" {
			// write file to given path
			err = imageFile.Write(outFile)
		} else {
			// save to current directory
			err = imageFile.Save()
		}
	}

	fmt.Println("\nComplete!")

	return err
}
