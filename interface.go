package main

import (
	"fmt"
	"image/color"
	"os/exec"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func drawInterface(fileName string, paragraph string, randomParagraph CleanParagraph, paragraphArray []CleanParagraph) {
	screenWidth := 1920
	screenHeight := 1080

	// Creates and runs the user interface.
	windowsApp := app.New()
	windowsApp.Settings().SetTheme(customTheme{})
	window := windowsApp.NewWindow("Paragraph of the Day")

	// Adds two horizontal lines to separate the different sections of the interface.
	line1 := canvas.NewLine(color.White)
	line1.StrokeWidth = 2
	line2 := canvas.NewLine(color.White)
	line2.StrokeWidth = 2

	// Adds the title of the text file.
	fileNameLabel := widget.NewLabel(fileName)
	fileNameLabel.Alignment = fyne.TextAlignCenter
	containerFileName := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), fileNameLabel, layout.NewSpacer())

	// Adds the paragraph text.
	paragraphLabel := widget.NewLabel(paragraph)
	containerParagraph := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		layout.NewSpacer(), paragraphLabel, layout.NewSpacer())

	// Adds button to go to the target file.
	documentButton := widget.NewButton("Open", func() {
		cmd := exec.Command("\\Program Files\\Microsoft Office\\root\\Office16\\winword.exe", randomParagraph.FilePath)
		err := cmd.Run()
		fmt.Print(err)
	})
	containerDocumentButton := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		layout.NewSpacer(), documentButton, layout.NewSpacer())

	// Adds button to obtain a new paragraph.
	refreshButton := widget.NewButton("Refresh", func() {
		randomParagraph = getParagraph(paragraphArray)
		paragraphText := formatText(randomParagraph.Text, 11)
		paragraphLabel.SetText(paragraphText)
		fileNameLabel.SetText(randomParagraph.FileName)
	})
	containerRefreshButton := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		layout.NewSpacer(), refreshButton, layout.NewSpacer())

	// Adds button to exit the program.
	exitButton := widget.NewButton("Exit", func() {
		windowsApp.Quit()
	})
	containerExitButton := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		layout.NewSpacer(), exitButton, layout.NewSpacer())

	// Creates a window containing all of the elements.
	window.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		containerFileName, line1,
		containerParagraph, line2,
		containerDocumentButton,
		containerRefreshButton,
		containerExitButton))

	// Resizes and shows the interface.
	window.Resize(fyne.NewSize(screenWidth, screenHeight))
	window.CenterOnScreen()
	window.ShowAndRun()
}
