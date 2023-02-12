package main

import (
	"fmt"
	"github.com/caarlos0/env/v7"
	"mdgen/gui"
)

// cli
func main() {
	data := new(gui.Data)
	err := env.Parse(data)
	if err != nil {
		panic(err)
	}
	generate(data, nil, func(err error) {
		fmt.Println(err)
	})
}
