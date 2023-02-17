package main

import (
	"mdgen/gui"
	"runtime/debug"
	"testing"
)

func TestGenLocal(t *testing.T) {
	generate(&gui.Data{
		Remote:           false,
		InputFolder:      "./posts",
		OutputFolder:     "./out",
		OutputPostFolder: "posts",
		GoogleAnalytics:  "AAAAA",
		BaseURL:          "https://localhost",
	}, make(chan struct{}))
}

func TestStack(t *testing.T) {
	stk := debug.Stack()
	t.Log(string(stk))
}
