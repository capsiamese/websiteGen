package gui

import (
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

type Data struct {
	Remote           bool   `env:"REMOTE"`
	InputFolder      string `env:"INPUT_FOLDER"`
	OutputFolder     string `env:"OUTPUT_FOLDER"`
	GoogleAnalytics  string `env:"GA"`
	BaseURL          string `env:"BASE_URL"`
	OutputPostFolder string `env:"POST_FOLDER"`

	RemoteAddr string `env:"ADDR"`
	RemotePort string `env:"PORT"`
	User       string `env:"USER"`
	Password   string `env:"PASSWORD"`
	KeyFile    string `env:"KEY_FILE"`
	KeyStr     string `env:"KEY_STR"`

	Clean bool // todo: impl clean dir
}

type App struct {
	win    *ui.Window
	data   Data
	remote ui.Control

	setupFuncList []func()

	startBtn *ui.Button
	progress *ui.ProgressBar
}

func NewGUI() *App {
	return &App{}
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

	pb := ui.NewProgressBar()
	a.progress = pb
	vBox.Append(pb, false)

	vBox.Append(a.draftComponent(), false)

	mw.SetChild(vBox)
	a.win = mw
	for _, v := range a.setupFuncList {
		v()
	}
}

func (a *App) Data() *Data {
	return &a.data
}

func (a *App) OnStartBtnClicked(f func(button *ui.Button)) {
	a.SetupF(func() {
		a.startBtn.OnClicked(f)
	})
}

func (a *App) SetProgress(n int) {
	a.Update(func() {
		a.progress.SetValue(n)
	})
}

func (a *App) Run() error {
	return ui.Main(a.setup)
}

func (a *App) Done() {
	a.progress.SetValue(0)
	a.startBtn.Enable()
	a.Msg("完成!")
}

func (a *App) SetupF(f func()) {
	a.setupFuncList = append(a.setupFuncList, f)
}
