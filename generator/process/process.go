package process

import (
	"bytes"
	"fmt"
	"generator/rec"
	"github.com/mozillazg/go-pinyin"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v3"
	"html/template"
	"os"
	"path"
	"strings"
	"time"
)

func (si *SiteInfo) AddPost(p *Post) {
	si.Posts = append(si.Posts, p)
}

func (si *SiteInfo) ScanPostDir() error {
	rec.Writeln("start scan posts folder ", si.Conf.InputFolder)

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
				rec.WritelnF("get post %s %v", fileName, err)
				continue
			} else if post.Meta.Draft {
				rec.Writeln("post ", fileName, " is draft ignore")
				continue
			} else {
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

	rec.Writeln("scan posts folder ", si.Conf.InputFolder, " done!")

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
	fullName := path.Join(dir, name)

	rec.Writeln("start parse post ", fullName)

	post := new(Post)
	post.FileName = name
	if i, err := info.Info(); err != nil {
		return nil, err
	} else {
		post.fileInfo = i
	}

	data, err := os.ReadFile(fullName)
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

	rec.Writeln("parse post ", fullName, " done!")

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
		rec.Writeln(err)
		return
	}
	err = yaml.Unmarshal(data, postMeta)
	if err != nil {
		rec.Writeln(err)
		return
	}
}

func (si *SiteInfo) GenerateIndex() error {
	rec.Writeln("start generate index.html")

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
	rec.Writeln("generate index.html  done!")
	return nil
}

func (si *SiteInfo) GeneratePosts() error {
	prefix := path.Join(si.Conf.OutputFolder, si.Conf.OutputPostFolder)

	rec.Writeln("start generate posts in ", prefix)

	err := writer.MkdirAll(prefix)
	if err != nil {
		return err
	}

	gen := func(filePath string, s *SiteInfo, p *Post) error {
		rec.Writeln("start write file to ", filePath)

		outputFileStat, err := writer.Stat(filePath)
		if err != nil {
			rec.WritelnF("get %s stat %v ", filePath, err)
		} else if !si.Conf.BuildAllPosts && p.fileInfo.ModTime().Before(outputFileStat.ModTime()) {
			rec.WritelnF("markdown modTime:%s exists html modTime:%s ignore",
				p.fileInfo.ModTime().Format(time.RFC3339),
				outputFileStat.ModTime().Format(time.RFC3339))
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
			rec.WritelnF("gen post %s %v ", v.OutputName, err)
		}
	}
	for _, v := range si.BannerPost {
		err := gen(path.Join(si.Conf.OutputFolder, v.OutputName), si, v)
		if err != nil {
			rec.WritelnF("gen post %s %v ", v.OutputName, err)
		}
	}

	rec.Writeln("generate posts in ", prefix, " done!")
	return nil
}

func (si *SiteInfo) WriteStatics() error {
	entries, err := static.ReadDir("static")
	if err != nil {
		return err
	}
	for _, v := range entries {
		targetF := path.Join(si.Conf.OutputFolder, v.Name())

		rec.Writeln("write static file to ", targetF)

		data, _ := static.ReadFile(path.Join("static", v.Name()))
		file, err := writer.Open(targetF)
		if err != nil {
			rec.WritelnF("gen post %s %v", v.Name(), err)
			continue
		}
		_, err = file.Write(data)
		file.Close()
		if err != nil {
			rec.WritelnF("gen post %s %v", v.Name(), err)
			continue
		}
	}
	return nil
}
