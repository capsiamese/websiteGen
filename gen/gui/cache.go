package gui

import (
	"encoding/json"
	"os"
)

func (a *App) ReadCache() {
	data, err := os.ReadFile("./.cache")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &a.data)
	if err != nil {
		return
	}
}

func (a *App) WriteCache() {
	tmp := a.data.Password
	a.data.Password = ""
	defer func() {
		a.data.Password = tmp
	}()
	data, err := json.Marshal(a.data)
	if err != nil {
		return
	}
	err = os.WriteFile("./.cache", data, os.ModePerm)
	if err != nil {
		return
	}
}
