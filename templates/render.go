package templates

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Show(name TemplateName) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		renderTemplate(w, r, name, nil)
	}
}

func Template(name TemplateName, data interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		renderTemplate(w, r, name, data)
	}
}

func Render(w http.ResponseWriter, r *http.Request, name TemplateName, data interface{}) {
	renderTemplate(w, r, name, data)
}

func renderTemplate(w http.ResponseWriter, r *http.Request, name TemplateName, data interface{}) {
	t := templates[name]
	if t == nil {
		http.NotFound(w, r)
		log.Print("template not found")
		return
	}
	t.Execute(w, data)
}
