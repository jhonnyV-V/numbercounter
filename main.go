package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func displayError(err error, window fyne.Window) {
	errorDiaglog := dialog.NewError(err, window)
	errorDiaglog.Show()
}

func initialBytes(uri fyne.URI) []byte {
	b := []byte(strings.Replace(uri.Extension(), uri.Name(), "", 1))
	b = append(b, '\n')
	b = append(b, 0)
	return b
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Box Layout")

	defaultPath := "C:/Users/user/Documents/examplefolder"
	fileLabel := widget.NewRichTextFromMarkdown(fmt.Sprintf("## %s", defaultPath))
	//TODO: "remember" last open URI
	var filePath fyne.ListableURI

	chooseFolderDialog := dialog.NewFolderOpen(
		func(lu fyne.ListableURI, err error) {
			if err != nil {
				displayError(err, myWindow)
				return
			}
			if lu == nil {
				return
			}
			fileLabel.ParseMarkdown(fmt.Sprintf("## %s", lu.Path()))
			filePath = lu
			childs, err := filePath.List()
			if err != nil {
				displayError(err, myWindow)
				fmt.Printf("failed to list files: %#v\n", err)
				return
			}
			for _, uri := range childs {
				fmt.Printf("files: %#v\n", uri)
			}
		},
		myWindow,
	)
	chooseFolderButton := widget.NewButton("Choose folder", func() {
		chooseFolderDialog.Show()
	})
	content := container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		chooseFolderButton,
		fileLabel,
		layout.NewSpacer(),
	)

	createCounterDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			displayError(err, myWindow)
			fmt.Printf("failed to save file %#v\n", err)
			return
		}
		if writer == nil {
			return
		}
		writer.Write(initialBytes(writer.URI()))
	}, myWindow)
	createCounterButton := widget.NewButton("Create new counter", func() {
		createCounterDialog.SetLocation(filePath)
		createCounterDialog.SetFileName("counter.txt")
		createCounterDialog.Show()
	})
	centered := container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		createCounterButton,
		layout.NewSpacer(),
	)

	countersContainer := container.NewVScroll(
		widget.NewCard(
			"counterName", "",
			widget.NewEntry(),
		),
	)
	buttonContainer := container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		countersContainer,
		container.New(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)

	mainContainer := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		content,
		centered,
		layout.NewSpacer(),
		buttonContainer,
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
	)

	myWindow.SetContent(mainContainer)
	myWindow.ShowAndRun()
}
