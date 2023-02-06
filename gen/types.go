package main

import (
	"fmt"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"go.abhg.dev/goldmark/mermaid"
	"html/template"
	"sort"
)

type SiteInfo struct {
	SiteTitle       string
	GoogleAnalytics string
	BaseURL         template.URL
	AboutURL        template.URL
	RefURL          template.URL
	Posts           []*PostInfo

	AboutPost *PostInfo
	RefPost   *PostInfo
}

type PostInfo struct {
	Title    string
	Date     string
	HexName  string
	RawName  string
	UseStyle bool
	Head     map[string]any
	URLPath  template.URL
	HTML     template.HTML

	ts int64
}

func OneOf[T comparable](src T, oneOf ...T) bool {
	for _, v := range oneOf {
		if v == src {
			return true
		}
	}
	return false
}

func Get[T any](m map[string]any, key string) (T, error) {
	v := m[key]
	r, ok := v.(T)
	if !ok {
		return r, fmt.Errorf("value of key %s is not %T", key, v)
	}
	return r, nil
}

func NewMarkDown() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			&mermaid.Extender{ /*```mermaid*/ },
			highlighting.NewHighlighting(
				highlighting.WithStyle("doom-one"),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(true),
				),
			),
		),
	)
}

func OrDefault[T comparable](v T, d T) T {
	var x T
	if v == x {
		return d
	}
	return v
}

func SortPost(list []*PostInfo) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].ts > list[j].ts
	})
}
