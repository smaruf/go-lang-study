package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create a new application
	myApp := app.New()
	
	// Create a new window
	myWindow := myApp.NewWindow("Hello Fyne")
	
	// Create a label
	hello := widget.NewLabel("Hello, Fyne!")
	
	// Create a button
	button := widget.NewButton("Click Me!", func() {
		hello.SetText("Button Clicked!")
	})
	
	// Create a container with vertical box layout
	content := container.NewVBox(
		hello,
		button,
		widget.NewButton("Reset", func() {
			hello.SetText("Hello, Fyne!")
		}),
	)
	
	// Set the content and show the window
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 200))
	myWindow.ShowAndRun()
}
