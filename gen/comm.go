package main

import (
	"embed"
	"fmt"
	"github.com/Masterminds/sprig"
	"html/template"
	"mdgen/gui"
	"os"
	"path"
	"strings"
)

//go:embed assets
var assetsFS embed.FS

//go:embed templates
var templates embed.FS

//go:embed static
var static embed.FS

var (
	markdown      = NewMarkDown()
	writer        FileWriter
	htmlTemplates *template.Template
)

func init() {
	htmlTemplates = template.Must(
		template.New("").
			Funcs(sprig.FuncMap()).
			Funcs(map[string]any{
				"asset": GetAsset,
			}).
			ParseFS(templates, "templates/*.html"),
	)
}

func GetAsset(name string) any {
	fn := path.Join("assets", name)
	switch path.Ext(name) {
	case ".js":
		data, _ := assetsFS.ReadFile(fn)
		return template.JS(data)
	case ".css":
		data, _ := assetsFS.ReadFile(fn)
		return template.CSS(data)
	}
	return nil
}

const (
	openMode = os.O_TRUNC | os.O_CREATE | os.O_RDWR
)

func generate(data *gui.Data, change func(int), onErr func(error)) {
	if data.Remote {
		r := NewRemoteWriter()
		if err := r.Connect(data); err != nil {
			onErr(fmt.Errorf("[fatal] connect to remote %v", err))
			return
		}
		writer = r
	} else {
		writer = LocalWriter{}
	}

	s := &SiteInfo{
		SiteTitle:       "siamese",
		GoogleAnalytics: data.GoogleAnalytics,
		BaseURL:         template.URL(strings.TrimRight(data.BaseURL, "/")),
		Conf:            data,
		BannerItems: []BannerItem{
			{URL: "/index.html", Target: "_self", Name: "Home"},
			{URL: "/about.html", Target: "_self", Name: "About", FileName: "about.md"},
			{URL: "/ref.html", Target: "_blank", Name: "Ref", FileName: "ref.md"},
		},
		ProgressChange:    change,
		OnErr:             onErr,
		IndexTemplateName: "blueprint",
		PostTemplateName:  "blueprint",
	}

	step := []func() error{
		s.ScanPostDir,
		s.GenerateIndex,
		s.GeneratePosts,
		s.WriteStatics,
	}
	for i, f := range step {
		if err := f(); err != nil {
			onErr(fmt.Errorf("run step:%d %v", i, err))
		}
	}

	if err := writer.PostRun(); err != nil {
		onErr(fmt.Errorf("post run %v", err))
	}
}
