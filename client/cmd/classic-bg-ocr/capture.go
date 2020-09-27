package main

import (
	"errors"
	"fmt"
	"image"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	d "github.com/lobsterbandit/wowclassic-bg-ocr/client/pkg/discord"
	"github.com/lobsterbandit/wowclassic-bg-ocr/client/pkg/img"
	"github.com/lobsterbandit/wowclassic-bg-ocr/client/pkg/screen"
	"github.com/spf13/cobra"
)

var (
	// points defining the capture rectangle
	points []string
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
	// run in daemon mode
	daemon bool
	// run command every N seconds
	every int

	captureCmd = &cobra.Command{
		Use:   "capture",
		Short: "capture the screen",
		Long:  "capture the full screen or a given rectangle",
		RunE:  runCapture,
	}

	ErrCaptureNotTwoPoints        = errors.New("capture requires exactly 2 points")
	ErrCapturePointMalformed      = errors.New("point value malformed, expecting x,y format")
	ErrCaptureIncompatibleFlagset = errors.New("incompatible set of flags")
	ErrWebhookMissingRequiredArgs = errors.New("webhookID and webhookToken are required")
)

func init() {
	addCaptureFlags()

	rootCmd.AddCommand(captureCmd)
}

func addCaptureFlags() {
	captureCmd.Flags().StringArrayVarP(
		&points, "point", "p", make([]string, 0, 2), "a point in X,Y format; 2 points define the capture rectangle")

	captureCmd.Flags().BoolVarP(&analyze, "analyze", "a", false, "analyze image for bg timers")
	captureCmd.Flags().StringVarP(&ocrURL, "url", "u", "", "remote ocr analysis endpoint, required if analyze is set")

	captureCmd.Flags().BoolVarP(
		&discord, "discord", "d", false, "send screen capture and analysis via webhook to a discord channel")
	captureCmd.Flags().StringVarP(&webhookID, "id", "i", "", "discord webhook id")
	captureCmd.Flags().StringVarP(&webhookToken, "token", "t", "", "discord webhook token")

	captureCmd.Flags().BoolVarP(&write, "write", "w", false, "write captured image to file")
	captureCmd.Flags().StringVarP(&outFile, "out", "o", "", "path to write captured image")

	captureCmd.Flags().StringVarP(&file, "file", "f", "", "path to a local image file")

	captureCmd.Flags().BoolVar(
		&daemon, "daemon", false, "in daemon mode command is run every N seconds controlled by the --every flag")
	captureCmd.Flags().IntVar(&every, "every", 15, "run command every N seconds; default=15")
}

func validateFlags() error {
	if n := len(points); n != 0 && n != 2 {
		return fmt.Errorf("%d points provided: %w", n, ErrCaptureNotTwoPoints)
	}

	if discord && (webhookID == "" || webhookToken == "") {
		return fmt.Errorf("missing required arguments to send discord webhooks: %w", ErrWebhookMissingRequiredArgs)
	}

	if analyze && ocrURL == "" {
		return fmt.Errorf("url is required if analyze is set: %w", ErrCaptureIncompatibleFlagset)
	}

	if daemon && file != "" {
		return fmt.Errorf("daemon and file flags cannot be set together: %w", ErrCaptureIncompatibleFlagset)
	}

	return nil
}

func executeCommand() (err error) {
	var imageFile *img.File
	if file != "" {
		imageFile, err = img.FromFile(file)
	} else {
		imageFile, err = captureScreen(points)
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

		err = webhook.PostDiscordMessage()
		if err != nil {
			return err
		}
	}

	if write {
		if outFile != "" {
			err = imageFile.Write(outFile)
		} else {
			err = imageFile.Save()
		}
	}

	fmt.Println("\nComplete!")

	return err
}

func runCapture(cmd *cobra.Command, args []string) (err error) {
	// exit early if required arg combinations are not met
	err = validateFlags()
	if err != nil {
		return err
	}

	if !daemon {
		return executeCommand()
	}

	return daemonMode()
}

func daemonMode() (err error) {
	fmt.Println("starting up daemon mode...")

	c := make(chan os.Signal, 1)
	commands := make(chan bool, 1)
	sleep := make(chan bool, 1)
	done := make(chan error, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-c
		fmt.Printf("received %v\n", sig)
		// cleanup
		done <- nil
	}()

	doCommand := func() {
		fmt.Println("executing capture command...")

		err = executeCommand()
		if err != nil {
			done <- err
		}

		sleep <- true
	}
	// fire off initial execution
	go doCommand()

	doSleep := func() {
		fmt.Printf("sleeping for %d seconds...\n", every)
		time.Sleep(time.Duration(every) * time.Second)

		commands <- true
	}

	for {
		select {
		case <-commands:
			go doCommand()
		case <-sleep:
			go doSleep()
		case err = <-done:
			fmt.Println("exiting...")

			return err
		}
	}
}

func captureScreen(points []string) (imageFile *img.File, err error) {
	if len(points) == 0 {
		return screen.CapturePrimary()
	}

	min := image.Point{X: math.MaxInt32, Y: math.MaxInt32}
	max := image.Point{X: 0, Y: 0}

	for _, point := range points {
		s := strings.Split(point, ",")

		// nolint:gomnd
		if len(s) != 2 {
			return nil, fmt.Errorf("error parsing points: %w", ErrCapturePointMalformed)
		}

		x, err := strconv.Atoi(s[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing X coord %v: %w", s[0], ErrCapturePointMalformed)
		}

		y, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing Y coord %v: %w", s[1], ErrCapturePointMalformed)
		}

		// ensure min and max arranged for rectangle creation
		if x < min.X {
			min.X = x
		}

		if x > max.X {
			max.X = x
		}

		if y < min.Y {
			min.Y = y
		}

		if y > max.Y {
			max.Y = y
		}
	}

	fmt.Printf("Capture Points\nMin: %v\nMax: %v\n\n", min, max)

	return screen.Capture(image.Rectangle{min, max})
}
