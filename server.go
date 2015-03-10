package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elos/autonomous"
	"github.com/elos/data"
	"github.com/elos/models"
	"github.com/elos/stack/util/logging"
	t "github.com/elos/transfer"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

type HTTPServer struct {
	host string
	port int

	autonomous.Life
	autonomous.Stopper

	*autonomous.Hub
	data.Store
	*httprouter.Router
}

func NewHTTPServer(host string, port int, s data.Store) *HTTPServer {
	server := new(HTTPServer)

	server.host = host
	server.port = port
	server.Store = s
	server.Hub = autonomous.NewHub()
	server.Life = autonomous.NewLife()
	server.Stopper = make(autonomous.Stopper)

	return server
}

func New(host string, port int, s data.Store) *HTTPServer {
	return NewHTTPServer(host, port, s)
}

func (s *HTTPServer) Start() {
	s.SetupRoutes()
	go s.Hub.Start()
	go s.Listen()
	s.Life.Begin()
	<-s.Stopper
	s.Life.End()
}

func list(v ...string) []string {
	return v
}

var ResourcesDir = os.Getenv("ELOS_HTTPSERVER_RESOURCES_DIR")

func (s *HTTPServer) SetupRoutes() {
	router := httprouter.New()

	router.POST("/v1/users/",
		Access(Post(models.UserKind, list("name")), data.NewAnonAccess(s.Store)))

	router.POST("/v1/events/",
		Auth(Post(models.EventKind, list("name")), t.Auth(t.HTTPCredentialer), s.Store))

	router.GET("/v1/authenticate",
		Auth(WebSocket(t.DefaultUpgrader, s), t.Auth(t.SocketCredentialer), s.Store))

	router.GET("/v1/repl",
		Auth(REPL(t.DefaultUpgrader, s), t.Auth(t.SocketCredentialer), s.Store))

	router.GET("/", Template("index.html"))
	router.GET("/sign-in", Template("sign-in.html"))
	router.POST("/sign-in", Auth(SignInHandle, t.Auth(t.FormCredentialer), s.Store))
	router.ServeFiles("/resources/*filepath", http.Dir(ResourcesDir))

	s.Router = router
}

func (a *HTTPServer) Listen() {
	serving_url := fmt.Sprintf("%s:%d", a.host, a.port)

	log.Printf("Serving at http://%s", serving_url)

	log.Fatal(http.ListenAndServe(serving_url, context.ClearHandler(logging.LogRequest(a.Router))))
}
