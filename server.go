package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elos/autonomous"
	"github.com/elos/data"
	"github.com/elos/stack/util/logging"
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

	server.Life = autonomous.NewLife()
	server.Stopper = make(autonomous.Stopper)

	server.Hub = autonomous.NewHub()
	server.Store = s
	server.Router = httprouter.New()

	return server
}

func New(host string, port int, s data.Store) *HTTPServer {
	return NewHTTPServer(host, port, s)
}

func (s *HTTPServer) Start() {
	setupAPI(s)
	setupRoutes(s)

	go s.Hub.Start()
	go s.Listen()
	s.Life.Begin()
	<-s.Stopper
	s.Life.End()
}

func (a *HTTPServer) Listen() {
	serving_url := fmt.Sprintf("%s:%d", a.host, a.port)

	log.Printf("Serving at http://%s", serving_url)

	log.Fatal(http.ListenAndServe(serving_url, context.ClearHandler(logging.LogRequest(a.Router))))
}
