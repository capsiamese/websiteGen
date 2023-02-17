package process

import (
	"generator/config"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"go.abhg.dev/goldmark/mermaid"
	"html/template"
	"io/fs"
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
	Conf            *config.Data
	BuildTime       time.Time

	IndexTemplateName string
	PostTemplateName  string
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
	ctx        map[string]any
	fileInfo   fs.FileInfo
}

func (p *Post) Get(k string) any {
	return p.ctx[k]
}

type PostMeta struct {
	UseStyle bool      `yaml:"styled"`
	Title    string    `yaml:"title"`
	Tags     []string  `yaml:"tags"`
	Date     time.Time `yaml:"date"`
	Draft    bool      `yaml:"draft"`

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

func (si *SiteInfo) CleanOutput() {
}
