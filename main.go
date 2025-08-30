package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var (
	//TODO: find a way to "cache" the last folder
	folderPath string
	files      []string
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

func loop(window *app.Window) error {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	var list layout.List

	for {
		event := window.Event()
		switch eventType := event.(type) {
		case app.DestroyEvent:
			return eventType.Err
		case app.FrameEvent:
			context := app.NewContext(&ops, eventType)
			xmargin := unit.Dp(context.Constraints.Max.X) / 5
			ymargin := unit.Dp(context.Constraints.Max.Y) / 10
			list = layout.List{Axis: layout.Vertical}

			margins := layout.Inset{
				Left:   xmargin,
				Right:  xmargin,
				Top:    ymargin,
				Bottom: ymargin,
			}

			if folderPath == "" {
				fmt.Printf("no folder\n")
			}

			margins.Layout(context,
				func(context layout.Context) layout.Dimensions {
					return list.Layout(context, len(files)+10, func(context layout.Context, index int) layout.Dimensions {
						textLabel := material.Label(theme, unit.Sp(16), "Placeholder")
						textLabel.Alignment = text.Middle
						return textLabel.Layout(context)
					})
				},
			)
			eventType.Frame(context.Ops)
		}
	}
}
