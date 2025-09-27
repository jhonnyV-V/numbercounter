package main

import (
	"fmt"
	"log"
	"mime"
	"os"
	"strconv"
	"strings"

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

type Counter struct {
	value           int
	fileName        *string
	incrementButton *widget.Clickable
	decrementButton *widget.Clickable
}

func (counter *Counter) increment() error {
	data := []byte(strconv.FormatInt(int64(counter.value+1), 10))
	name := fmt.Sprintf("%s/%s", folderPath, *counter.fileName)

	w, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("(counter.increment) Failed open file %s: %w", name, err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("(counter.increment) Failed write file %s: %w", name, err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("(counter.increment) Failed close file %s: %w", name, err)
	}

	counter.value += 1

	return nil
}

func (counter *Counter) decrement() error {
	data := []byte(strconv.FormatInt(int64(counter.value-1), 10))
	name := fmt.Sprintf("%s/%s", folderPath, *counter.fileName)

	w, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("(counter.decrement) Failed open file %s: %w", name, err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("(counter.decrement) Failed write file %s: %w", name, err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("(counter.decrement) Failed close file %s: %w", name, err)
	}

	counter.value -= 1

	return nil
}

var (
	//TODO: find a way to "cache" the last folder
	//TODO: remove this default value
	folderPath string = ""
	files      []string
	counters   []*Counter
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Size(unit.Dp(800), unit.Dp(700)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func readFolderForCounters() {
	listOfFiles, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Printf("listOfFiles ERR %#v\n", err.Error())
	}
	fmt.Printf("listOfFiles %#v\n", listOfFiles)
	counters = []*Counter{}
	for _, v := range listOfFiles {
		if v.IsDir() {
			continue
		}
		name := v.Name()
		fmt.Printf("name %s\n", name)

		splitName := strings.Split(name, ".")
		fmt.Printf("splitName %#v\n", splitName)

		extension := mime.TypeByExtension("." + splitName[len(splitName)-1])
		fmt.Printf("extension %s\n", extension)

		if !strings.Contains(extension, "text/") {
			fmt.Printf("does not contain text/\n")
			continue
		}

		data, err := os.ReadFile(fmt.Sprintf("%s/%s", folderPath, name))

		if err != nil {
			fmt.Printf("failed to read file %s %#v\n", name, err)
		}
		value, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err != nil {
			fmt.Printf("failed to read value %s %#v\n", string(data), err)
		}

		counters = append(counters, &Counter{
			value:           value,
			fileName:        &name,
			incrementButton: new(widget.Clickable),
			decrementButton: new(widget.Clickable),
		})
	}
}

func loop(window *app.Window) error {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	var list layout.List
	var createCounterButton widget.Clickable
	var selectFolderButton widget.Clickable

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
				}
			}

			if createCounterButton.Clicked(context) && folderPath != "" {
				ex := explorer.NewExplorer(window)
				//modify to take the folderPath as a default directory
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

			layoutMargin.Layout(context,
				func(context layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Vertical,
						Alignment: layout.Middle,
					}.Layout(context,
						layout.Rigid(func(context layout.Context) layout.Dimensions {
							button := material.Button(theme, &selectFolderButton, "Select a folder")
							if folderPath == "" {
								return button.Layout(context)
							}

							label := material.Label(
								theme,
								unit.Sp(20),
								fmt.Sprintf("Current Directory: %s", folderPath),
							)

							label.WrapPolicy = text.WrapHeuristically
							textMargin := layout.Inset{
								Top: unit.Dp(unit.Sp(6)),
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
						layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout),
						layout.Rigid(func(context layout.Context) layout.Dimensions {
							if folderPath == "" {
								return layout.Spacer{}.Layout(context)
							}
							button := material.Button(theme, &createCounterButton, "Create counter")

							return button.Layout(context)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout),
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
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(context,
		layout.Rigid(func(context layout.Context) layout.Dimensions {
			textLabel := material.Label(
				theme, unit.Sp(24),
				strconv.FormatInt(int64(counter.value), 10),
			)
			textLabel.Alignment = text.Middle
			return textLabel.Layout(context)
		}),
		layout.Rigid(func(context layout.Context) layout.Dimensions {
			buttonsMargin := layout.Inset{
				Bottom: unit.Dp(10),
				Left:   unit.Dp(15),
			}
			return layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(context,
				layout.Rigid(func(context layout.Context) layout.Dimensions {
					button := material.Button(theme, counter.incrementButton, "+")
					return buttonsMargin.Layout(context, button.Layout)
				}),
				layout.Rigid(func(context layout.Context) layout.Dimensions {
					button := material.Button(theme, counter.decrementButton, "-")
					return buttonsMargin.Layout(context, button.Layout)
				}),
			)
		}),
	)
}
