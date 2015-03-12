package httpserver

import (
	"github.com/elos/data"
	"github.com/elos/models"
	t "github.com/elos/transfer"
)

func list(v ...string) []string {
	return v
}

func setupAPI(s *HTTPServer) {
	s.POST("/v1/users/",
		Access(Post(models.UserKind, list("name")), data.NewAnonAccess(s.Store)))

	s.POST("/v1/events/",
		Auth(Post(models.EventKind, list("name")), t.Auth(t.HTTPCredentialer), s.Store))

	s.GET("/v1/authenticate",
		Auth(WebSocket(t.DefaultUpgrader, s), t.Auth(t.SocketCredentialer), s.Store))

	s.GET("/v1/repl",
		Auth(REPL(t.DefaultUpgrader, s), t.Auth(t.SocketCredentialer), s.Store))
}
