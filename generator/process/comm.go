package process

import (
	"embed"
	"generator/config"
	"generator/rec"
	"github.com/Masterminds/sprig"
	"html/template"
	"os"
	"path"
	"strings"
	"time"
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

func Generate(data *config.Data, done chan<- struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	if data.Remote {
		r := NewRemoteWriter()
		if err := r.Connect(data); err != nil {
			rec.WritelnF("[fatal] connect to remote %v", err)
			return
		}
		writer = r
	} else {
		writer = LocalWriter{}
	}
	rec.Writeln("*****************************************")
	rec.WritelnF("start generate use remote %v", data.Remote)

	s := &SiteInfo{
		SiteTitle:       "siamese",
		GoogleAnalytics: data.GoogleAnalytics,
		Conf:            data,
		BaseURL:         template.URL(strings.TrimRight(data.BaseURL, "/")),
		BannerItems: []BannerItem{
			{URL: "/index.html", Target: "_self", Name: "Home"},
			{URL: "/about.html", Target: "_self", Name: "About", FileName: "about.md"},
			{URL: "/ref.html", Target: "_blank", Name: "Ref", FileName: "ref.md"},
		},
		IndexTemplateName: "blueprint",
		PostTemplateName:  "blueprint",
		BuildTime:         time.Now(),
	}

	step := []func() error{
		s.ScanPostDir,
		s.GenerateIndex,
		s.GeneratePosts,
		s.WriteStatics,
	}
	for i, f := range step {
		if err := f(); err != nil {
			rec.WritelnF("run step:%d %v", i, err)
			return
		}
	}

	if err := writer.PostRun(); err != nil {
		rec.WritelnF("post run %v", err)
		return
	}
}
