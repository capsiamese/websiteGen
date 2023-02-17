package main

import (
	"errors"
	"fmt"
	"github.com/andlabs/ui"
	"github.com/sqweek/dialog"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

func (a *App) pickerComponent(btnName, pickerTitle, pickerTxt string, file bool, recv func(folder string)) ui.Control {
	h := ui.NewHorizontalBox()
	label := ui.NewLabel(pickerTxt)
	btn := ui.NewButton(btnName)
	btn.OnClicked(func(button *ui.Button) {
		var pth string
		var err error
		if file {
			pth, err = dialog.File().Title(pickerTitle).Load()
		} else {
			pth, err = dialog.Directory().Title(pickerTitle).Browse()
		}
		if err != nil {
			if !errors.Is(err, dialog.ErrCancelled) {
				a.Error(err)
				return
			}
		} else {
			recv(pth)
			label.SetText(pth)
		}
	})

	h.SetPadded(true)
	h.Append(btn, false)
	h.Append(ui.NewVerticalSeparator(), false)
	h.Append(label, true)

	return h
}

func (a *App) targetSelector(remote bool) ui.Control {
	rb := ui.NewRadioButtons()
	rb.Append("local")
	rb.Append("remote")
	if remote {
		rb.SetSelected(1)
		a.hideRemote(false)
	} else {
		rb.SetSelected(0)
		a.hideRemote(true)
	}
	rb.OnSelected(func(buttons *ui.RadioButtons) {
		switch buttons.Selected() {
		case 0:
			a.hideRemote(true)
		case 1:
			a.hideRemote(false)
		}
	})
	return rb
}

func EntryStr(v *string) ui.Control {
	entry := ui.NewEntry()
	entry.SetText(*v)
	entry.OnChanged(func(entry *ui.Entry) {
		*v = entry.Text()
	})
	return entry
}

func (a *App) formComponent() ui.Control {
	form := ui.NewForm()
	form.SetPadded(true)
	form.Append("选择输出模式", a.targetSelector(a.data.Remote), false)

	form.Append("GoogleAnalytics", EntryStr(&a.data.GoogleAnalytics), false)
	form.Append("BaseURL", EntryStr(&a.data.BaseURL), false)
	form.Append("文章存储目录名", EntryStr(&a.data.OutputPostFolder), false)

	return form
}

func (a *App) draftComponent() ui.Control {
	btn := ui.NewButton("生成草稿")
	btn.OnClicked(func(button *ui.Button) {
		fn := ui.SaveFile(a.win)
		newDraft(fn)
	})
	return btn
}

func newDraft(target string) {
	if err := os.MkdirAll(path.Dir(target), 0644); err != nil {
		log.Fatalln("[fatal]", err)
	}
	f, err := os.OpenFile(target, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalln("[fatal]", err)
	}
	defer f.Close()

	_, fn := path.Split(filepath.ToSlash(target))

	_, _ = fmt.Fprintln(f, "---")
	_, _ = fmt.Fprintf(f, "title: \"%s\"\n", fn[:len(fn)-len(path.Ext(fn))])
	_, _ = fmt.Fprintf(f, "date: %s\n", time.Now().Format(time.RFC3339))
	_, _ = fmt.Fprintf(f, "draft: true\n")
	_, _ = fmt.Fprintf(f, "tags: []\n")
	_, _ = fmt.Fprintln(f, "---")
}
