package main

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/andlabs/ui"
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

/*
TODO:
*/
var (
	assetsDir   string = "./assets"
	templateDir string = "./templates"
	staticDir   string = "./statics"
)

var (
	markdown      = NewMarkDown()
	writer        FileWriter
	htmlTemplates *template.Template
)

func IniTemplate() {
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
	if path.Ext(name) == ".js" {
		data, _ := assetsFS.ReadFile("assets/" + name)
		return template.JS(data)
	} else {
		data, _ := assetsFS.ReadFile("assets/" + name)
		return template.CSS(data)
	}
}

const (
	openMode = os.O_TRUNC | os.O_CREATE | os.O_RDWR
)

func main() {
	IniTemplate()

	app := gui.NewGUI()

	app.OnStartBtnClicked(func(button *ui.Button) {
		button.Disable()
		app.SetProgress(0)

		generate(app.Data(), app.SetProgress, func(err error) {})

		app.Done()
	})

	app.SetupF(func() {

	})

	err := app.Run()
	if err != nil {
		fmt.Println("----", err, "----")
	}
}

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
		ProgressChange: change,
		OnErr:          onErr,
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
