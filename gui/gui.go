package main

import (
	"fmt"
	"generator/config"
	"generator/rec"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

type App struct {
	win    *ui.Window
	data   config.Data
	remote ui.Control

	setupFuncList []func()

	startBtn *ui.Button

	doneCh chan struct{}
	rc     rec.Recorder
}

func NewGUI(opt ...Option) *App {
	a := &App{}
	for _, v := range opt {
		v(a)
	}
	if a.rc == nil {
		a.rc = rec.Default()
	}
	return a
}

func (a *App) Update(f func()) {
	ui.QueueMain(f) // for other goroutine
}

func (a *App) Error(msg any) {
	ui.MsgBoxError(a.win, "Error", fmt.Sprintf("%s", msg))
}

func (a *App) Msg(msg any) {
	ui.MsgBox(a.win, "Message", fmt.Sprintf("%s", msg))
}

func (a *App) setup() {
	mw := ui.NewWindow("Hello", 800, 600, false)
	defer mw.Show()
	mw.OnClosing(func(window *ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mw.Destroy()
		return true
	})

	vBox := ui.NewVerticalBox()
	vBox.SetPadded(true)

	r := a.remoteComponent()

	vBox.Append(a.formComponent(), false)
	vBox.Append(ui.NewVerticalSeparator(), false)

	vBox.Append(a.pickerComponent("选择输入目录", "选择输入目录", a.data.InputFolder, false, func(folder string) {
		a.data.InputFolder = folder
	}), false)
	vBox.Append(ui.NewVerticalSeparator(), false)

	vBox.Append(a.pickerComponent("选择输出目录", "选择输出目录", a.data.OutputFolder, false, func(folder string) {
		a.data.OutputFolder = folder
	}), false)
	vBox.Append(ui.NewVerticalSeparator(), false)
	vBox.Append(r, false)

	goBtn := ui.NewButton("GO")
	a.startBtn = goBtn
	vBox.Append(goBtn, false)

	vBox.Append(a.draftComponent(), false)

	wl, ok := a.rc.(*GUIRec)
	if ok {
		wl.Attach(ui.NewMultilineEntry())
		vBox.Append(wl.Control(), false)
		wl.QueueMain(a.Update)
	}

	mw.SetChild(vBox)
	a.win = mw
	for _, v := range a.setupFuncList {
		v()
	}
	a.doneCh = make(chan struct{})
	a.startListen()
}

func (a *App) Data() *config.Data {
	return &a.data
}

func (a *App) OnStartBtnClicked(f func(button *ui.Button)) {
	a.SetupF(func() {
		a.startBtn.OnClicked(f)
	})
}

func (a *App) Run() error {
	return ui.Main(a.setup)
}

func (a *App) SetupF(f func()) {
	a.setupFuncList = append(a.setupFuncList, f)
}

func (a *App) DoneChan() chan<- struct{} {
	return a.doneCh
}

func (a *App) startListen() {
	go func(app *App) {
		for {
			select {
			case <-app.doneCh:
				app.Update(func() {
					app.WriteCache()
					app.startBtn.Enable()
					app.Msg("完成!")
				})
			}
		}
	}(a)
}
