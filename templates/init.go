package templates

import (
	"html/template"
	"log"
	"path/filepath"
)

// not exposed because all templates should be refered to
// by some render function in the templates package, not externally
var templates = map[Name]*template.Template{}

func init() {
	// templateSets defined in conf.go
	if err := parseHTMLTemplates(templateSets); err != nil {
		log.Fatal(err)
	}
}

/*
	joinDir prepends a full path to a slice of relative paths

	i.e., base = "/root/here/", files  ["1.go", "2.go"]
	=> ["/root/here/1.go", "/root/here/2.go"]

	Useful for building templateSet paths
*/
func joinDir(base string, files []string) []string {
	r := make([]string, len(files))
	for i := range files {
		r[i] = filepath.Join(base, files[i])
	}
	return r
}

func parseHTMLTemplates(sets map[Name][]string) error {
	for name, set := range sets {
		t, err := template.ParseFiles(joinDir(TemplatesDir, set)...)
		if err != nil {
			return err
		}
		templates[name] = t
	}
	return nil
}
