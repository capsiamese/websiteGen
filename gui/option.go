package main

import (
	"generator/rec"
)

type Option func(a *App)

func WithLog(l rec.Recorder) Option {
	return func(a *App) {
		a.rc = l
	}
}
