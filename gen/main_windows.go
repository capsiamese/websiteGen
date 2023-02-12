package main

import (
	_ "embed"
	"fmt"
	"github.com/andlabs/ui"
	"mdgen/gui"
)

// gui
func main() {
	app := gui.NewGUI()
	app.ReadCache()

	app.OnStartBtnClicked(func(button *ui.Button) {
		button.Disable()
		app.SetProgress(0)

		generate(app.Data(), app.SetProgress, func(err error) {})
		app.WriteCache()

		app.Done()
	})

	app.SetupF(func() {

	})

	err := app.Run()
	if err != nil {
		fmt.Println("----", err, "----")
	}
}
