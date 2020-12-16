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
	fileNameText := widget.NewLabel(fileName)
	fileNameText.Alignment = fyne.TextAlignCenter
	containerFileName := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), fileNameText, layout.NewSpacer())

	// Adds the paragraph text.
	paragraphText := widget.NewLabel(paragraph)
	containerParagraph := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		layout.NewSpacer(), paragraphText, layout.NewSpacer())

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
		newRandomParagraph := getParagraph(paragraphArray)
		newParagraphText := formatText(newRandomParagraph.Text, 11)
		paragraphText.SetText(newParagraphText)
		fileNameText.SetText(newRandomParagraph.FileName)
	})
	containerRefreshButton := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		layout.NewSpacer(), refreshButton, layout.NewSpacer())

	// Creates a window containing all of the elements.
	window.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		containerFileName, line1,
		containerParagraph, line2,
		containerDocumentButton,
		containerRefreshButton))

	// Resizes and shows the interface.
	window.Resize(fyne.NewSize(400, 0))
	window.CenterOnScreen()
	window.ShowAndRun()
}
