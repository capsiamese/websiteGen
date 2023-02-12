package main

import (
	"mdgen/gui"
	"testing"
)

func TestGenLocal(t *testing.T) {
	IniTemplate()
	generate(&gui.Data{
		Remote:           false,
		InputFolder:      "./posts",
		OutputFolder:     "./out",
		OutputPostFolder: "posts",
		GoogleAnalytics:  "AAAAA",
		BaseURL:          "https://localhost",
	}, func(i int) {
		t.Log("progress", i)
	}, func(err error) {
		t.Log(err)
	})
}
