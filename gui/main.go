package main

import (
	"fmt"
	"generator/process"
	"generator/rec"
	"github.com/andlabs/ui"
)

func main() {
	recorder := rec.NewFileRec("run.log")

	app := NewGUI(WithLog(recorder))
	app.ReadCache()

	app.OnStartBtnClicked(func(button *ui.Button) {
		if !button.Enabled() {
			return
		}
		button.Disable()

		go process.Generate(app.Data(), app.DoneChan())
	})

	app.SetupF(func() {})

	err := app.Run()
	if err != nil {
		fmt.Println("----", err, "----")
	}
}
