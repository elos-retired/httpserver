package httpserver

import (
	"html/template"
	// "io"
	"log"
	"path/filepath"
)

var templateSets = map[string][]string{
	"index":    {"index.html"},
	"sign-in":  {"layout.tmpl", "sessions.tmpl", "sign-in.tmpl"},
	"register": {"layout.tmpl", "sessions.tmpl", "register.tmpl"},
}

func init() {
	if err := parseHTMLTemplates(templateSets); err != nil {
		log.Fatal(err)
	}
}

var templates = map[string]*template.Template{}

func joinDir(base string, files []string) []string {
	r := make([]string, len(files))
	for i := range files {
		r[i] = filepath.Join(templatesDir, files[i])
	}
	return r
}

func parseHTMLTemplates(sets map[string][]string) error {
	for name, set := range sets {
		t, err := template.ParseFiles(joinDir(templatesDir, set)...)
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
