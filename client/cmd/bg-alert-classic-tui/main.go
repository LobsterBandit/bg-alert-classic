package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	currentCapture := tview.NewTextView().SetText(`Current image of screen capure area.
Cache image and retake every X seconds or when capture area changes.`)
	currentCapture.Box.SetBorder(true).SetTitle("Latest Screen Capture")

	captureOverlay := tview.NewTextView().
		SetText(`Show overlay of current capture area over total screen area.
Simple wireframe display of two rectangles of different transparencies.`)
	captureOverlay.Box.SetBorder(true).SetTitle("Capture Overlay")

	captureAreaForm := tview.NewTextView().SetText("Form to set the capture area")
	captureAreaForm.Box.SetBorder(true).SetTitle("Capture Area")

	captureArea := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(captureOverlay, 0, 1, false).
		AddItem(captureAreaForm, 0, 1, false)

	ocrResults := tview.NewTextView().SetText(`Last preprocessed image and parsed result`)
	ocrResults.Box.SetBorder(true).SetTitle("OCR Results")

	commands := tview.NewTextView().SetText(`Streaming commands and form to set command parameters`)
	commands.Box.SetBorder(true).SetTitle("Commands")

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(currentCapture, 0, 1, false).
			AddItem(captureArea, 0, 1, false), 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(ocrResults, 0, 1, false).
			AddItem(commands, 0, 1, false), 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
