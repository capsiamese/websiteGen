package main

import (
	"fmt"
	"mime"
	"net/http"
	"net/http/cgi"
)

func main() {
	if err := cgi.Serve((*Handler)(nil)); err != nil {
		fmt.Print("Content-type:", mime.TypeByExtension(".html"), "\n\n")
		fmt.Println("<h1>", err, "</h1>")
	}
}

type Handler struct{}

//w.Write(*(*[]byte)(unsafe.Pointer(&data)))

func (Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	r.URL.Query()
	_, _ = rw.Write([]byte("<h1>Hello World</h1>"))
}
