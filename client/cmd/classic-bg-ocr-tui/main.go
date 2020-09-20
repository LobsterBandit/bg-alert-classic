package main

import (
	"github.com/rivo/tview"
)

const text = `The terminal ui
Showing screen bounds,
selected screen area within those bounds,
last screenshot,
last processed image from ocr server,
server parsed response,
etc`

func main() {
	app := tview.NewApplication()

	view := tview.NewModal().SetText(text)
	if err := app.SetRoot(view, false).Run(); err != nil {
		panic(err)
	}
}
