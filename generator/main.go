package main

import (
	"generator/config"
	"generator/process"
	"generator/rec"
	"github.com/caarlos0/env"
)

func main() {
	data := new(config.Data)
	err := env.Parse(data)
	if err != nil {
		panic(err)
	}
	ch := make(chan struct{})
	go process.Generate(data, ch)
	<-ch
	rec.Default().Writeln("generate ok")
}
