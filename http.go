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

	*autonomous.Core
	*autonomous.AgentHub
	data.Store
	*httprouter.Router

	SocketRequests chan *agents.ClientDataAgent
}

func NewHTTPServer(host string, port int, s data.Store) *HTTPServer {
	return &HTTPServer{
		host:           host,
		port:           port,
		Core:           autonomous.NewCore(),
		AgentHub:       autonomous.NewAgentHub(),
		Store:          s,
		SocketRequests: make(chan *agents.ClientDataAgent, 10),
	}
}

func New(host string, port int, s data.Store) *HTTPServer {
	return NewHTTPServer(host, port, s)
}

func (s *HTTPServer) Run() {
	s.startup()
	stopChannel := s.Core.StopChannel()

	for {
		select {
		case a := <-s.SocketRequests:
			s.AgentHub.StartAgent(a)
		case _ = <-*stopChannel:
			s.shutdown()
			break
		}
	}
}

func (a *HTTPServer) startup() {
	a.Core.Startup()
	a.SetupRoutes()
	go a.Listen()
}

func (a *HTTPServer) shutdown() {
	a.Core.Shutdown()
}

func list(v ...string) []string {
	return v
}

func (s *HTTPServer) SetupRoutes() {
	router := httprouter.New()

	router.POST("/v1/users/", Auth(Post(models.UserKind, list("name")), s.Store))

	router.POST("/v1/events/", Auth(Post(models.EventKind, list("name")), s.Store))

	router.GET("/v1/authenticate", Auth(WebSocket(transfer.DefaultWebSocketUpgrader, s.SocketRequests), s.Store))

	s.Router = router
}

func (a *HTTPServer) Listen() {
	serving_url := fmt.Sprintf("%s:%d", a.host, a.port)

	log.Print("Serving at http://%s", serving_url)

	log.Fatal(http.ListenAndServe(serving_url, logging.LogRequest(a.Router)))
}
