package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
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
	ocrURL       string
	save         bool
	webhookID    string
	webhookToken string

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
	// TODO: add output flag for file save location
	captureCmd.Flags().StringVarP(&file, "file", "f", "", "path to a local image file")
	captureCmd.Flags().BoolVarP(&save, "save", "s", false, "save captured image to file")
	captureCmd.Flags().StringVarP(&webhookID, "id", "i", "", "discord webhook id")
	captureCmd.Flags().StringVarP(&webhookToken, "token", "t", "", "discord webhook token")

	rootCmd.AddCommand(captureCmd)
}

func runCapture(cmd *cobra.Command, args []string) (err error) {
	fmt.Printf("\nwowclassic-bg-ocr-client %v\n\tWoW Classic BG timer screen capture and analysis\n\n", version)

	// exit early if required arg combinations are not met
	if discord && (webhookID == "" || webhookToken == "") {
		return fmt.Errorf("missing required arguments to send discord webhooks: %w", ErrWebhookMissingRequiredArgs)
	}

	if analyze && ocrURL == "" {
		return fmt.Errorf("url is required if analyze is set: %w", ErrCaptureIncompatibleFlagset)
	}

	var imageFile *img.File
	if file != "" {
		// open local file and set capture to that
		imageFile, err = loadFromFile(file)
	} else {
		// capture screen
		imageFile, err = captureScreen()
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

	if save {
		err = imageFile.Save()
	}

	fmt.Println("\nComplete!")

	return err
}

func loadFromFile(path string) (*img.File, error) {
	infile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	imageFile, err := png.Decode(infile)
	if err != nil {
		return nil, err
	}

	stats, err := infile.Stat()
	if err != nil {
		return nil, err
	}

	return &img.File{
		Name:      infile.Name(),
		Timestamp: img.Timestamp(stats.ModTime().Unix()),
		Image:     imageFile.(*image.RGBA),
	}, nil
}

func captureScreen() (*img.File, error) {
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
	if err != nil {
		return nil, err
	}

	captureTime := time.Now()
	fileName := fmt.Sprintf("%dx%d_%s_%s_%d.png",
		capture.Rect.Dx(), capture.Rect.Dy(), bgArea.Min, bgArea.Max, captureTime.Unix())

	fmt.Printf("Captured screen area: %v\n\tTimestamp: %s\n\tFilename: %q\n\n", bgArea, captureTime, fileName)

	return &img.File{
		Name:      fileName,
		Timestamp: img.Timestamp(captureTime.Unix()),
		Image:     capture,
	}, nil
}
