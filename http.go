package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elos/agents"
	"github.com/elos/autonomous"
	"github.com/elos/data"
	"github.com/elos/models"
	"github.com/elos/stack/util/logging"
	"github.com/elos/transfer"
	"github.com/julienschmidt/httprouter"
)

type HTTPServer struct {
	host string
	port int

	autonomous.Life
	autonomous.Stopper

	*autonomous.AgentHub
	data.Store
	*httprouter.Router

	SocketRequests chan *agents.ClientDataAgent
}

func NewHTTPServer(host string, port int, s data.Store) *HTTPServer {
	server := new(HTTPServer)
	server.host = host
	server.port = port
	server.AgentHub = autonomous.NewAgentHub()
	server.Store = s
	server.Life = autonomous.NewLife()

	return server
}

func New(host string, port int, s data.Store) *HTTPServer {
	return NewHTTPServer(host, port, s)
}

func (s *HTTPServer) Run() {
	s.startup()
	s.Life.Begin()

Run:
	for {
		select {
		case _ = <-s.Stopper:
			break Run
		}
	}

	s.shutdown()
	s.Life.End()
}

func (a *HTTPServer) startup() {
	a.SetupRoutes()
	go a.Listen()
}

func (a *HTTPServer) shutdown() {
}

func list(v ...string) []string {
	return v
}

func (s *HTTPServer) SetupRoutes() {
	router := httprouter.New()

	router.POST("/v1/users/", Auth(Post(models.UserKind, list("name")), s.Store))

	router.POST("/v1/events/", Auth(Post(models.EventKind, list("name")), s.Store))

	router.GET("/v1/authenticate", Auth(WebSocket(transfer.DefaultWebSocketUpgrader, s), s.Store))

	s.Router = router
}

func (a *HTTPServer) Listen() {
	serving_url := fmt.Sprintf("%s:%d", a.host, a.port)

	log.Print("Serving at http://%s", serving_url)

	log.Fatal(http.ListenAndServe(serving_url, logging.LogRequest(a.Router)))
}
