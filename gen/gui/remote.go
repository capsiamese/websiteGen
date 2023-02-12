package gui

import "github.com/andlabs/ui"

func (a *App) remoteComponent() ui.Control {
	v := ui.NewVerticalBox()
	a.remote = v

	form := ui.NewForm()

	addrC := ui.NewHorizontalBox()
	addrC.Append(EntryStr(&a.data.RemoteAddr), false)
	port := ui.NewEntry()
	port.SetText("22")
	a.data.RemotePort = "22"
	port.OnChanged(func(entry *ui.Entry) {
		a.data.RemotePort = entry.Text()
	})
	addrC.Append(port, false)

	form.Append("SSH地址", addrC, false)

	form.Append("User", EntryStr(&a.data.User), false)
	pe := ui.NewPasswordEntry()
	pe.OnChanged(func(entry *ui.Entry) {
		a.data.Password = entry.Text()
	})
	form.Append("Password", pe, false)
	form.Append("远程输出目录", EntryStr(&a.data.OutputFolder), false)

	v.Append(form, false)
	v.Append(a.pickerComponent("选择SSH key文件", "选择SSH key文件", true, func(folder string) {
		a.data.KeyFile = folder
	}), false)

	return v
}

func (a *App) hideRemote(h bool) {
	a.data.Remote = !h
	if a.remote != nil {
		if h {
			a.remote.Hide()
			a.remote.Disable()
		} else {
			a.remote.Show()
			a.remote.Enable()
		}
	}
}

func FolderPicker() {
	win := ui.NewWindow("", 100, 200, false)
	win.OnClosing(func(window *ui.Window) bool {
		win.Destroy()
		return true
	})
	frame := ui.NewVerticalBox()

	frame.Append(ui.NewHorizontalSeparator(), false)
	h := ui.NewHorizontalBox()
	okBtn := ui.NewButton("确定")
	okBtn.OnClicked(func(button *ui.Button) {
		win.Destroy()
	})
	cancelBtn := ui.NewButton("取消")
	cancelBtn.OnClicked(func(button *ui.Button) {
		win.Destroy()
	})
	h.Append(okBtn, false)
	h.Append(cancelBtn, false)
	frame.Append(h, false)
	win.SetChild(frame)
	win.Show()
}
