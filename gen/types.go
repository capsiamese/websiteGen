package main

import (
	"fmt"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"go.abhg.dev/goldmark/mermaid"
	"html/template"
	"mdgen/gui"
	"os"
	"sort"
	"time"
)

type SiteInfo struct {
	SiteTitle       string
	GoogleAnalytics string
	BaseURL         template.URL
	Posts           []*Post
	BannerItems     []BannerItem
	BannerPost      []*Post
	Conf            *gui.Data

	complete       int
	total          int
	ProgressChange func(n int)
	OnErr          func(error)
}

func (si *SiteInfo) AddBannerPost(post *Post) {
	si.BannerPost = append(si.BannerPost, post)
}

func (si *SiteInfo) IsBannerPost(fileName string) bool {
	for i := range si.BannerItems {
		if si.BannerItems[i].FileName == fileName {
			return true
		}
	}
	return false
}

func (si *SiteInfo) AddTotal(n int) {
	si.total += n
	if si.total != 0 {
		si.ProgressChange(int(float32(si.complete/si.total) * 100))
	}
}

func (si *SiteInfo) AddComplete(n int) {
	si.complete += n
	if si.total != 0 {
		si.ProgressChange(int(float32(si.complete/si.total) * 100))
	}
}

func (si *SiteInfo) SortPost() {
	sort.Slice(si.Posts, func(i, j int) bool {
		return si.Posts[i].Meta.Date.Unix() > si.Posts[j].Meta.Date.Unix()
	})
}

type Post struct {
	FileName   string
	OutputName string
	URL        template.URL
	HTML       template.HTML
	Meta       PostMeta
	FileInfo   os.FileInfo
}

type PostMeta struct {
	UseStyle bool
	Title    string
	Tags     []string
	Date     time.Time
	Draft    bool

	Private bool
}

type BannerItem struct {
	URL      template.URL
	Target   string
	Name     string
	FileName string
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
