package rec

import "log"

var defaultRec Recorder = def{}

type Recorder interface {
	Writeln(i ...any)
	WritelnF(f string, i ...any)
}

func Writeln(i ...any) {
	defaultRec.Writeln(i...)
}

func WritelnF(f string, i ...any) {
	defaultRec.WritelnF(f, i...)
}

func Default() Recorder {
	return defaultRec
}

type def struct{}

func (def) Writeln(i ...any) {
	log.Println(i...)
}
func (def) WritelnF(f string, i ...any) {
	log.Printf(f, i...)
	log.Println()
}
