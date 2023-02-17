package main

import (
	"fmt"
	"generator/rec"
	"github.com/andlabs/ui"
	"log"
	"time"
)

func NewGUIRec() *GUIRec { // can not resize widget, so not use
	r := &GUIRec{}
	rec.SetDefault(r)
	return r
}

type GUIRec struct {
	win   *ui.MultilineEntry
	queue func(func())
}

func (l *GUIRec) Attach(w *ui.MultilineEntry) {
	l.win = w
}

func (l *GUIRec) Control() ui.Control {
	return l.win
}

func (l *GUIRec) QueueMain(f func(func())) {
	l.queue = f
}

func (l *GUIRec) write(str string) {
	if l.queue == nil {
		log.Print(str)
		return
	}
	l.queue(func() {
		l.win.Append(time.Now().Format(time.RFC3339))
		l.win.Append(" ")
		l.win.Append(str)
	})
}

func (l *GUIRec) Writeln(i ...any) {
	l.win.Append(fmt.Sprintln(i...))
}
func (l *GUIRec) WritelnF(f string, i ...any) {
	l.win.Append(fmt.Sprintf(f+"\n", i...))
}
