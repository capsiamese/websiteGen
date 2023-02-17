package rec

import (
	"github.com/sirupsen/logrus"
	"os"
)

// todo: add bt

type FileRec struct {
	l *logrus.Logger
}

func NewFileRec(fn string) *FileRec {
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	l := logrus.New()
	l.SetOutput(f)
	l.SetLevel(logrus.DebugLevel)
	r := &FileRec{
		l: l,
	}
	defaultRec = r
	return r
}

func (fr *FileRec) Writeln(i ...any) {
	fr.l.Print(i...)
}

func (fr *FileRec) WritelnF(f string, i ...any) {
	fr.l.Printf(f, i...)
}
