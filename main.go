package main

import (
	"fmt"
	"log"
	"mime"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Counter struct {
	value           int
	fileName        *string
	incrementButton *widget.Clickable
	decrementButton *widget.Clickable
}

var (
	//TODO: find a way to "cache" the last folder
	//TODO: remove this default value
	folderPath string = "/home/jhonny/Documents/counters-DELETEME"
	files      []string
	counters   []*Counter
)

func main() {

	listOfFiles, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Printf("listOfFiles ERR %#v\n", err.Error())
	}
	fmt.Printf("listOfFiles %#v\n", listOfFiles)
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

		counters = append(counters, &Counter{
			value:           0,
			fileName:        &name,
			incrementButton: new(widget.Clickable),
			decrementButton: new(widget.Clickable),
		})
	}

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

func loop(window *app.Window) error {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	var list layout.List
	var createCounterButton widget.Clickable

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

			if folderPath == "" {
				fmt.Printf("no folder\n")
			}

			layoutMargin.Layout(context,
				func(context layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Vertical,
						Alignment: layout.Middle,
					}.Layout(context,
						layout.Rigid(func(context layout.Context) layout.Dimensions {
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
									// textLabel.Layout,
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
			textLabel := material.Label(theme, unit.Sp(24), "0")
			textLabel.Alignment = text.Middle
			return textLabel.Layout(context)
		}),
		layout.Rigid(func(context layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(context,
				layout.Rigid(func(context layout.Context) layout.Dimensions {
					button := material.Button(theme, counter.incrementButton, "+")
					return button.Layout(context)
				}),
				layout.Rigid(func(context layout.Context) layout.Dimensions {
					button := material.Button(theme, counter.incrementButton, "-")
					return button.Layout(context)
				}),
			)
		}),
	)
}
