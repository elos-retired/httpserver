package httpserver

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"go/build"

	"github.com/elos/data"
	"github.com/elos/models"
	"github.com/elos/models/user"
	t "github.com/elos/transfer"
	"github.com/julienschmidt/httprouter"
)

var (
	assetsDir    = filepath.Join(defaultBase("github.com/elos/httpserver"), "assets")
	templatesDir = filepath.Join(assetsDir, "templates")
	imgDir       = filepath.Join(assetsDir, "img")
	cssDir       = filepath.Join(assetsDir, "css")
)

func setupRoutes(s *HTTPServer) {
	s.GET("/", Template("index", nil))
	s.GET("/sign-in", Template("sign-in", nil))
	s.POST("/sign-in", Auth(SignInHandle, t.Auth(t.FormCredentialer), s.Store))
	s.GET("/register", Template("register", nil))
	s.POST("/register", RegisterHandle(s.Store))

	s.ServeFiles("/css/*filepath", http.Dir(cssDir))
	s.ServeFiles("/img/*filepath", http.Dir(imgDir))
}

func defaultBase(path string) string {
	p, err := build.Default.Import(path, "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

type Page struct {
}

func Template(name string, data interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		renderTemplate(w, r, name, data)
	}
}

func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	t := templates[name]
	if t == nil {
		http.NotFound(w, r)
		log.Print("template not found")
		return
	}
	t.Execute(w, data)
}

func RegisterTemplate(n string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		path1 := filepath.Join(templatesDir, n)
		path2 := filepath.Join(templatesDir, "layout.html")
		t, err := template.ParseFiles(path1, path2)
		if err != nil {
			log.Print("Template error: %s", err)
			return
		}
		log.Print(t.Execute(w, Page{}))
	}
}

func SignInHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	w.Write([]byte("hello" + a.Client().(models.User).Name()))
}

func RegisterHandle(s data.Store) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		u, err := user.NewWithName(s, r.FormValue("name"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		renderTemplate(w, r, "account-created", u)
	}
}
