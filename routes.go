package httpserver

import (
	"fmt"
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
	s.GET("/", Template("index.html"))
	s.GET("/sign-in", Template("sign-in.html"))
	s.POST("/sign-in", Auth(SignInHandle, t.Auth(t.FormCredentialer), s.Store))
	s.GET("/register", Template("register.html"))
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

func Template(name string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		path := filepath.Join(templatesDir, name)
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

func RegisterHandle(s data.Store) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		u, err := user.NewWithName(s, r.FormValue("name"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		WelcomeHandle(w, u)
	}
}

func WelcomeHandle(w http.ResponseWriter, u models.User) {
	s := fmt.Sprintf("Welcome %s, we have created an account for you. These are your credentials: ID: %s | Key: %s ",
		u.Name(), u.ID().String(), u.Key())

	w.Write([]byte(s))
}
