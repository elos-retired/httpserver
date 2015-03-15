package httpserver

import (
	"net/http"

	"github.com/elos/httpserver/handles"
	"github.com/elos/httpserver/templates"
	t "github.com/elos/transfer"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	sessionsStore      = sessions.NewCookieStore([]byte("something-very-secret"), securecookie.GenerateRandomKey(32))
	CookieCredentialer = t.NewCookieCredentialer(sessionsStore)
	CookieAuth         = t.Auth(CookieCredentialer)
)

func setupRoutes(s *HTTPServer) {
	s.GET("/", templates.Show(templates.Index))

	s.GET("/sign-in", templates.Show(templates.SignIn))
	s.POST("/sign-in", handles.Auth(handles.SignIn(sessionsStore), t.Auth(t.FormCredentialer), s.Store))
	s.GET("/register", templates.Show(templates.Register))
	s.POST("/register", handles.RegisterHandle(s.Store))

	s.GET("/calendar", handles.Auth(handles.CalendarHandle, CookieAuth, s.Store))

	s.ServeFiles("/img/*filepath", http.Dir(templates.ImgDir))
	s.ServeFiles("/css/*filepath", http.Dir(templates.CSSDir))
}
