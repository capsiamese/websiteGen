package main

import (
	"github.com/caarlos0/env/v7"
	"mdgen/gui"
	"mdgen/rec"
)

// cli
func main() {
	data := new(gui.Data)
	err := env.Parse(data)
	if err != nil {
		panic(err)
	}
	ch := make(chan struct{})
	generate(data, ch)
	<-ch
	rec.Default().Writeln("generate ok")
}
