package main

import (
	_ "embed"
	"fmt"
	"github.com/andlabs/ui"
	"mdgen/gui"
	"mdgen/rec"
)

// gui
func main() {
	recorder := rec.NewFileRec("run.log")

	app := gui.NewGUI(gui.WithLog(recorder))
	app.ReadCache()

	app.OnStartBtnClicked(func(button *ui.Button) {
		if !button.Enabled() {
			return
		}
		button.Disable()

		go generate(app.Data(), app.DoneChan())
	})

	app.SetupF(func() {})

	err := app.Run()
	if err != nil {
		fmt.Println("----", err, "----")
	}
}
