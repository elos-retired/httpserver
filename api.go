package httpserver

import (
	"github.com/elos/data"
	h "github.com/elos/httpserver/handles"
	"github.com/elos/models"
	t "github.com/elos/transfer"
)

func list(v ...string) []string {
	return v
}

func setupAPI(s *HTTPServer) {
	s.POST("/v1/users/",
		h.Access(h.Post(models.UserKind, list("name")), data.NewAnonAccess(s.Store)))

	s.POST("/v1/events/",
		h.Auth(h.Post(models.EventKind, list("name")), t.Auth(t.HTTPCredentialer), s.Store))

	s.GET("/v1/authenticate",
		h.Auth(h.WebSocket(t.DefaultUpgrader, s), t.Auth(t.SocketCredentialer), s.Store))

	s.GET("/v1/repl",
		h.Auth(h.REPL(t.DefaultUpgrader, s), t.Auth(t.SocketCredentialer), s.Store))
}
