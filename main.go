package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/jhonnyV-V/gio-x/explorer"
)

var (
	folderPath string = ""
	didReadFromCache = false
	counters   []*Counter
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Size(unit.Dp(800), unit.Dp(700)))
		app.Title("Counter")
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(window *app.Window) error {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	var list layout.List
	var createCounterButton widget.Clickable
	var selectFolderButton widget.Clickable
	var reloadCountersButton widget.Clickable

	for {
		event := window.Event()
		switch eventType := event.(type) {
		case app.DestroyEvent:
			return eventType.Err
		case app.FrameEvent:
			context := app.NewContext(&ops, eventType)
			xmargin := unit.Dp(context.Constraints.Max.X) / 5
			ymargin := unit.Dp(context.Constraints.Max.Y) / 10
			list = layout.List{
				Axis: layout.Vertical,
			}

			layoutMargin := layout.Inset{
				Left:   xmargin,
				Right:  xmargin,
				Top:    ymargin,
				Bottom: ymargin,
			}

			if selectFolderButton.Clicked(context) {
				ex := explorer.NewExplorer(window)
				selectedFolder, err := ex.ChooseFolder(nil)
				if err != nil {
					fmt.Printf("failed to open folder: %s \n", err)
				} else {
					folderPath = selectedFolder
					readFolderForCounters()
					cachePath := getOrCreateConfigDir()
					writeToCache(cachePath, folderPath)
				}
			}

			if createCounterButton.Clicked(context) && folderPath != "" {
				ex := explorer.NewExplorer(window)
				//TODO:modify to take the folderPath as a default directory
				//also return filename or path
				w, err := ex.CreateFile("counter.txt")
				if err != nil {
					fmt.Printf("Failed to create file %#v\n", err)
				}

				_, err = w.Write([]byte{'0'})
				if err != nil {
					fmt.Printf("Failed to write to file %#v\n", err)
				}

				err = w.Close()
				if err != nil {
					fmt.Printf("Failed to close file %#v\n", err)
				}

				readFolderForCounters()
			}

			for _, counter := range counters {
				if counter.decrementButton.Clicked(context) {
					err := counter.decrement()
					if err != nil {
						fmt.Printf("Failed to decrement %#v\n", err)
					}
				}

				if counter.incrementButton.Clicked(context) {
					err := counter.increment()
					if err != nil {
						fmt.Printf("Failed to increment %#v\n", err)
					}
				}
			}

			if reloadCountersButton.Clicked(context) {
				readFolderForCounters()
			}

			if folderPath == "" && !didReadFromCache {
				cachePath := getOrCreateConfigDir()
				didReadFromCache = true
				err := readFromCache(cachePath)
				if err != nil {
					fmt.Printf("Failed to read cache %s\n", err)
				}
			}

			layoutMargin.Layout(context,
				func(context layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Vertical,
						Alignment: layout.Middle,
					}.Layout(context,

						//NOTE: Select Folder 
						layout.Rigid(func(context layout.Context) layout.Dimensions {
							button := material.Button(theme, &selectFolderButton, "Select a folder")
							if folderPath == "" {
								return button.Layout(context)
							}

							label := material.Label(
								theme,
								unit.Sp(20),
								fmt.Sprintf("Current: %s", folderPath),
							)

							label.WrapPolicy = text.WrapHeuristically
							textMargin := layout.Inset{
								Top:  unit.Dp(unit.Sp(6)),
								Left: unit.Dp(5),
							}

							return layout.Flex{
								Axis: layout.Horizontal,
							}.Layout(context,
								layout.Rigid(button.Layout),
								layout.Rigid(
									func(context layout.Context) layout.Dimensions {
										return textMargin.Layout(context, label.Layout)
									},
								),
							)
						}),

						//NOTE: Reload Button
						layout.Rigid(func(context layout.Context) layout.Dimensions {
							if folderPath == "" {
								return layout.Spacer{}.Layout(context)
							}
							button := material.Button(theme, &reloadCountersButton, "Reload counters")
							return layout.Inset{
								Top: unit.Dp(20),
							}.Layout(context, button.Layout)
						}),

						layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout),

						//NOTE: Create Button
						layout.Rigid(func(context layout.Context) layout.Dimensions {
							if folderPath == "" {
								return layout.Spacer{}.Layout(context)
							}
							button := material.Button(theme, &createCounterButton, "Create counter")

							return button.Layout(context)
						}),

						layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout),

						//NOTE: Counter List
						layout.Rigid(func(context layout.Context) layout.Dimensions {
							return list.Layout(context, len(counters), func(context layout.Context, index int) layout.Dimensions {
								ymargin := unit.Dp(15)
								return layout.Inset{
									Top:    ymargin,
									Bottom: ymargin,
								}.Layout(
									context,
									func(context layout.Context) layout.Dimensions {
										return getCounter(context, index, theme)
									},
								)
							})
						}),
					)
				},
			)
			eventType.Frame(context.Ops)
		}
	}
}

func getCounter(context layout.Context, index int, theme *material.Theme) layout.Dimensions {
	counter := counters[index]

	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(context,
		//NOTE: Labels
		layout.Rigid(func(context layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(context,
				layout.Rigid(func(context layout.Context) layout.Dimensions {
					textLabel := material.Label(
						theme, unit.Sp(40),
						printFileName(*counter.fileName),
					)
					textLabel.Alignment = text.Middle
					return textLabel.Layout(context)
				}),
				layout.Rigid(func(context layout.Context) layout.Dimensions {
					labelText := strconv.FormatInt(int64(counter.value), 10)
					if counter.value < 10 && counter.value > -1 {
						labelText = "0" + labelText
					}
					textLabel := material.Label(
						theme, unit.Sp(35), labelText,
					)
					textLabel.Alignment = text.End
					return textLabel.Layout(context)
				}),
			)
		}),

		//NOTE: Buttons
		layout.Rigid(func(context layout.Context) layout.Dimensions {
			buttonsMargin := layout.Inset{
				Left: unit.Dp(10),
			}
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(context,
				layout.Rigid(func(context layout.Context) layout.Dimensions {
					button := material.Button(theme, counter.incrementButton, "+")
					return button.Layout(context)
				}),
				layout.Rigid(func(context layout.Context) layout.Dimensions {
					button := material.Button(theme, counter.decrementButton, "-")
					return buttonsMargin.Layout(context, button.Layout)
				}),
			)
		}),
	)
}
