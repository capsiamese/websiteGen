package main

import (
	"bytes"
	"fmt"
	"github.com/mozillazg/go-pinyin"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v3"
	"html/template"
	"os"
	"path"
	"strings"
)

func (si *SiteInfo) Error(format string, args ...any) {
	si.OnErr(fmt.Errorf(format, args...))
}

func (si *SiteInfo) AddPost(p *Post) {
	si.Posts = append(si.Posts, p)
}

func (si *SiteInfo) ScanPostDir() error {
	entities, err := os.ReadDir(si.Conf.InputFolder)
	if err != nil {
		return err
	}
	for _, v := range entities {
		fileName := v.Name()
		isMarkDown := OneOf(path.Ext(fileName), ".md", ".MD", ".mD", ".Md")
		if !v.IsDir() && isMarkDown {
			post, err := GetPost(si.Conf.InputFolder, fileName, v)
			if err != nil {
				si.OnErr(fmt.Errorf("get post %s %v", fileName, err))
				continue
			} else if post.Meta.Draft {
				continue
			} else {
				si.AddTotal(1)
				if si.IsBannerPost(fileName) {
					post.URL = template.URL(fmt.Sprintf("/%s", post.OutputName))
					si.AddBannerPost(post)
				} else {
					post.URL = template.URL(fmt.Sprintf("/%s/%s", si.Conf.OutputPostFolder, post.OutputName))
					si.AddPost(post)
				}
			}
		}
	}
	si.SortPost()
	return nil
}

var pinyinArgs = &pinyin.Args{
	Style:     pinyin.NORMAL,
	Heteronym: false,
	Separator: pinyin.Separator,
	Fallback: func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	},
}

func GetPost(dir, name string, info os.DirEntry) (*Post, error) {
	post := new(Post)
	post.FileName = name
	if i, err := info.Info(); err != nil {
		return nil, err
	} else {
		post.fileInfo = i
	}

	data, err := os.ReadFile(path.Join(dir, name))
	if err != nil {
		return nil, err
	}

	ctx := parser.NewContext()
	buf := bytes.NewBuffer(make([]byte, 0, 4*1024*1024))
	err = markdown.Convert(data, buf, parser.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	post.ctx = meta.Get(ctx)
	SetMeta(ctx, &post.Meta)

	post.HTML = template.HTML(buf.String())

	py := strings.Builder{}
	for _, v := range pinyin.LazyConvert(TrimExt(name), pinyinArgs) {
		py.WriteString(v)
	}
	py.WriteString(".html")
	post.OutputName = py.String()

	return post, nil
}

func TrimExt(name string) string {
	s := path.Ext(name)
	if len(s) < len(name) {
		return name[:len(name)-len(s)]
	}
	return name
}

func SetMeta(ctx parser.Context, postMeta *PostMeta) {
	// todo: handle error
	data, err := yaml.Marshal(meta.Get(ctx))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = yaml.Unmarshal(data, postMeta)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (si *SiteInfo) GenerateIndex() error {
	err := writer.MkdirAll(si.Conf.OutputFolder)
	if err != nil {
		return err
	}

	file, err := writer.Open(path.Join(si.Conf.OutputFolder, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()
	err = htmlTemplates.ExecuteTemplate(file, si.IndexTemplateName, map[string]any{
		"Site":  si,
		"Index": true,
	})
	if err != nil {
		return err
	}
	return nil
}

func (si *SiteInfo) GeneratePosts() error {
	prefix := path.Join(si.Conf.OutputFolder, si.Conf.OutputPostFolder)
	err := writer.MkdirAll(prefix)
	if err != nil {
		return err
	}

	gen := func(filePath string, s *SiteInfo, p *Post) error {
		outputFileStat, err := writer.Stat(filePath)
		if err != nil {
			// todo: handle error
		} else if p.fileInfo.ModTime().Before(outputFileStat.ModTime()) {
			return nil
		}

		file, err := writer.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = htmlTemplates.ExecuteTemplate(file, si.PostTemplateName, map[string]any{
			"Site": s,
			"Post": p,
		})
		if err != nil {
			return err
		}
		return nil
	}

	for _, v := range si.Posts {
		err := gen(path.Join(prefix, v.OutputName), si, v)
		if err != nil {
			si.Error("gen post %s %v", v.OutputName, err)
		}
		si.AddComplete(1)
	}
	for _, v := range si.BannerPost {
		err := gen(path.Join(si.Conf.OutputFolder, v.OutputName), si, v)
		if err != nil {
			si.Error("gen post %s %v", v.OutputName, err)
		}
		si.AddComplete(1)
	}
	return nil
}

func (si *SiteInfo) WriteStatics() error {
	entries, err := static.ReadDir("static")
	if err != nil {
		return err
	}
	for _, v := range entries {
		data, _ := static.ReadFile(path.Join("static", v.Name()))
		file, err := writer.Open(path.Join(si.Conf.OutputFolder, v.Name()))
		if err != nil {
			si.Error("gen post %s %v", v.Name(), err)
			continue
		}
		_, err = file.Write(data)
		file.Close()
		if err != nil {
			si.Error("gen post %s %v", v.Name(), err)
			continue
		}
	}
	return nil
}
