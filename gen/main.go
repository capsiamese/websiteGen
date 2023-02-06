package main

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"fmt"
	"github.com/Masterminds/sprig"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"html/template"
	"log"
	"os"
	"path"
	"time"
)

/*
TODO:
	1. WriteFile static file compare create time
	2. template refactor (content start --- content end)
*/

//go:embed index.tmpl.html
var indexTmpl string

//go:embed content.tmpl.html
var contentTmpl string

//go:embed favicon.ico
var icon []byte

//go:embed bg.webm
var bg []byte

var (
	indexTemplate   = template.Must(template.New("").Funcs(sprig.FuncMap()).Parse(indexTmpl))
	contentTemplate = template.Must(template.New("").Funcs(sprig.FuncMap()).Parse(contentTmpl))
	markdown        = NewMarkDown()
	writer          FileWriter
)

const (
	about      = "about.md"
	references = "reference.md"
	aboutHtml  = "about.html"
	refHtml    = "reference.html"

	openMode = os.O_TRUNC | os.O_CREATE | os.O_RDWR
)

func newDraft(target string) {
	f, err := os.OpenFile(target, openMode, 0644)
	if err != nil {
		log.Fatalln("[fatal]", err)
	}
	defer f.Close()

	_, _ = fmt.Fprintln(f, "---")
	_, _ = fmt.Fprintf(f, "title: \"%s\"\n", path.Base(target))
	_, _ = fmt.Fprintf(f, "date: %s\n", time.Now().Format(time.RFC3339))
	_, _ = fmt.Fprintf(f, "draft: true\n")
	_, _ = fmt.Fprintf(f, "tags: []\n")
	_, _ = fmt.Fprintln(f, "---")
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "new" {
		newDraft(os.Args[2])
		return
	}

	ParseConfig()
	if config.Method == "local" {
		writer = LocalWriter{}
	} else if config.Method == "remote" {
		r := NewRemoteWriter()
		if err := r.Connect(); err != nil {
			log.Fatalln("[fatal] connect to remote", err)
		}
		writer = r
	} else {
		log.Fatalln("[fatal] config field 'method' must be local or remote")
	}

	s := &SiteInfo{
		SiteTitle:       "siamese",
		GoogleAnalytics: config.GoogleAnalytics,
		BaseURL:         template.URL(config.BaseURL),
		AboutURL:        "/" + aboutHtml,
		RefURL:          "/" + refHtml,
	}
	if err := scanPosts(s, config.InputDir, config.PostFolder); err != nil {
		log.Fatalln("[fatal] scan posts", err)
	}
	if err := Exec(s, config.OutputDir, config.PostFolder); err != nil {
		log.Fatalln("[fatal] gen html", err)
	}
	WriteFile(path.Join(config.OutputDir, "favicon.ico"), icon)
	WriteFile(path.Join(config.OutputDir, "bg.webm"), bg)

	if err := writer.PostRun(); err != nil {
		log.Fatalln("[fatal] run post command", err)
	}
}

func scanPosts(s *SiteInfo, inputDir, postDir string) error {
	entities, err := os.ReadDir(inputDir)
	if err != nil {
		return err
	}
	log.Println("[info] input directory", inputDir)
	log.Println("[info] start scanning post")
	for _, v := range entities {
		name := v.Name()
		isMd := OneOf(path.Ext(name), ".md", ".MD", ".mD", ".Md")
		if !v.IsDir() && isMd {
			info, err := parsePost(inputDir, name, "/"+postDir)
			if err != nil {
				log.Fatalln("[fatal] parse post", name, err)
			} else if info != nil {
				if name == about {
					info.UseStyle = true
					s.AboutPost = info
				} else if name == references {
					info.UseStyle = true
					s.RefPost = info
				} else {
					s.Posts = append(s.Posts, info)
				}
				log.Println("[info] scanning", name, "done")
			}
		}
	}
	SortPost(s.Posts)
	return nil
}

func parsePost(dir, fn, urlPrefix string) (*PostInfo, error) {
	srcFullName := path.Join(dir, fn)
	data, err := os.ReadFile(srcFullName)
	if err != nil {
		return nil, err
	}

	ctx := parser.NewContext()
	buf := bytes.NewBuffer(make([]byte, 0, 1024*1024*4))
	err = markdown.Convert(data, buf, parser.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	mdMeta := meta.Get(ctx)
	if isDraft(mdMeta) {
		return nil, nil
	}
	info := new(PostInfo)
	fillMeta(mdMeta, info)

	info.RawName = fn
	info.HexName = hex.EncodeToString([]byte(fn))
	info.UseStyle = false
	info.Head = mdMeta
	info.URLPath = template.URL(fmt.Sprintf("%s/%s.html", urlPrefix, info.HexName))
	info.HTML = template.HTML(buf.String())

	return info, nil
}

func isDraft(m map[string]any) bool {
	if d, ok := m["draft"]; ok {
		if dd, ok := d.(bool); ok {
			return dd
		}
	}
	return true
}

func fillMeta(m map[string]any, postInfo *PostInfo) {
	title, err := Get[string](m, "title")
	if err != nil {
		log.Println("[warning] fill meta key", "title", err)
	}
	postInfo.Title = title

	// todo: pretty time format
	date, err := Get[string](m, "date")
	if err != nil {
		log.Println("[warning] fill meta key", "date", err)
	}
	t, err := time.Parse(time.RFC3339, date)
	if err == nil {
		postInfo.Date = t.Format("Jan _2, 2006")
		postInfo.ts = t.Unix()
		//postInfo.Date = strconv.FormatInt(t.Unix(), 10)
	} else {
		postInfo.Date = date
		postInfo.ts = 0
	}
}

func Exec(s *SiteInfo, output, postDir string) error {
	err := writer.MkdirAll(path.Join(output, postDir))
	if err != nil {
		return err
	}
	log.Println("[info] start generate html")

	if err = ExecTemplate(s, nil, path.Join(output, "index.html")); err != nil {
		log.Println("[error] generate index.html", err)
	}
	if err = ExecTemplate(s, s.RefPost, path.Join(output, refHtml)); err != nil {
		log.Println("[error] generate reference.html", err)
	}
	if err = ExecTemplate(s, s.AboutPost, path.Join(output, aboutHtml)); err != nil {
		log.Println("[error] generate about.html", err)
	}
	log.Println("[info] generate index.html about.html reference.html done")

	for _, v := range s.Posts {
		if err = ExecTemplate(s, v, path.Join(output, postDir, v.HexName+".html")); err != nil {
			log.Println("[error] generate", v.RawName, err)
		} else {
			log.Println("[info] generate", v.RawName, "html done")
		}
	}
	return nil
}

func ExecTemplate(s *SiteInfo, post *PostInfo, fn string) error {
	file, err := writer.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()
	if post == nil {
		return indexTemplate.Execute(file, s)
	} else {
		return contentTemplate.Execute(file, map[string]any{
			"Site": s, "Post": post,
		})
	}
}

func WriteFile(fn string, data []byte) {
	f, err := writer.Open(fn)
	if err != nil {
		log.Println("[error] open", fn, err)
		return
	}
	defer f.Close()
	log.Println("[info] start write", fn)
	_, err = f.Write(data)
	if err != nil {
		log.Println("[error] write", fn, err)
		return
	}
	log.Println("[info] write", fn, "done")
}
