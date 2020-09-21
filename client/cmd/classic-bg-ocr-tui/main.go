package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	topLeft := tview.NewTextView().SetText(`Current image of screen capure area`)
	topLeft.Box.SetBorder(true).SetTitle("Top Left")

	captureOverlay := tview.NewTextView().
		SetText("Show overlay of current capture area over total screen area")
	captureOverlay.Box.SetBorder(true).SetTitle("Capture Overlay")

	captureAreaForm := tview.NewTextView().SetText("Form to set the capture area")
	captureAreaForm.Box.SetBorder(true).SetTitle("Capture Area")

	bottomLeft := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(captureOverlay, 0, 1, false).
		AddItem(captureAreaForm, 0, 1, false)

	topRight := tview.NewTextView().SetText(`Last preprocessed image and parsed result`)
	topRight.Box.SetBorder(true).SetTitle("Top Right")

	bottomRight := tview.NewTextView().SetText(`Streaming commands and form to set command parameters`)
	bottomRight.Box.SetBorder(true).SetTitle("Bottom Right")

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(topLeft, 0, 1, false).
			AddItem(bottomLeft, 0, 1, false), 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(topRight, 0, 1, false).
			AddItem(bottomRight, 0, 1, false), 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
