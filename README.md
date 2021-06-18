# Yet Another Graphical toolKit

A very opinionated GUI toolkit with the most basic style definition possible.

Useful to build a quick multiplatform (linux, win & macos) desktop app.  

## Basic Example
```go
package main

import (
	"github.com/enimatek-nl/yagk"
	"os"
)

func main() {
	app := yagk.New("My App", 512, 512)

	style := yagk.NewStyle(&yagk.StyleDefinition{})
	
	pane := yagk.NewPane(0, 0, app.Width, style)

	pane.Button(
		"Quit",
		func() {
			os.Exit(1)
		},
	)
	
	app.PushPane(pane)
	
	app.Run()
}
```