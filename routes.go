package httpserver

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/elos/data"
	"github.com/elos/models"
	"github.com/julienschmidt/httprouter"
)

type Page struct {
}

func Template(name string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		path := filepath.Join(ResourcesDir, name)
		t, err := template.ParseFiles(path)
		if err != nil {
			log.Print("Template error: %s", err)
			return
		}
		t.Execute(w, Page{})
	}
}

func SignInHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	w.Write([]byte("hello" + a.Client().(models.User).Name()))
}
