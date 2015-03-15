package templates

import (
	"go/build"
	"log"
	"path/filepath"
	"text/template"
)

func PackagePath(importPath string) string {
	p, err := build.Default.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

var (
	httpserverPath = "github.com/elos/httpserver"
	AssetsDir      = filepath.Join(PackagePath(httpserverPath), "assets")
	TemplatesDir   = filepath.Join(AssetsDir, "templates")
	ImgDir         = filepath.Join(AssetsDir, "img")
	CSSDir         = filepath.Join(AssetsDir, "css")
)

type TemplateName int

const (
	Index TemplateName = iota
	SignIn
	Register
	AccountCreated
	Calendar
)

var templateSets = map[TemplateName][]string{
	Index:          {"layout.tmpl", "index.html"},
	SignIn:         {"layout.tmpl", "sessions.tmpl", "sign-in.tmpl"},
	Register:       {"layout.tmpl", "sessions.tmpl", "register.tmpl"},
	AccountCreated: {"layout.tmpl", "sessions.tmpl", "account-created.tmpl"},
	Calendar:       {"layout.tmpl", "calendar.tmpl"},
}

var templates = map[TemplateName]*template.Template{}

func init() {
	if err := parseHTMLTemplates(templateSets); err != nil {
		log.Fatal(err)
	}
}

func joinDir(base string, files []string) []string {
	r := make([]string, len(files))
	for i := range files {
		r[i] = filepath.Join(TemplatesDir, files[i])
	}
	return r
}

func parseHTMLTemplates(sets map[TemplateName][]string) error {
	for name, set := range sets {
		t, err := template.ParseFiles(joinDir(TemplatesDir, set)...)
		if err != nil {
			return err
		}
		/*
			t = t.Lookup("ROOT")
			if t == nil {
				return fmt.Errorf("ROOT template not found in %v", set)
			}
		*/
		templates[name] = t
	}
	return nil
}
