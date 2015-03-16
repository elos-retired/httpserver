package httpserver

import (
	"net/http"

	"github.com/elos/data"
	"github.com/elos/httpserver/handles"
	"github.com/elos/httpserver/templates"
	t "github.com/elos/transfer"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

var (
	sessionsStore      = sessions.NewCookieStore([]byte("something-very-secret"), securecookie.GenerateRandomKey(32))
	CookieCredentialer = t.NewCookieCredentialer(sessionsStore)
	CookieAuth         = t.Auth(CookieCredentialer)
)

func UserAuth(f handles.AccessHandle, s data.Store) httprouter.Handle {
	return handles.Auth(f, CookieAuth, s)
}

const (
	UserBase             = "/user"
	UserCalendar         = UserBase + "/calendar"
	UserEvents           = UserBase + "/events"
	UserTasks            = UserBase + "/tasks"
	UserRoutines         = UserBase + "/routines"
	UserSchedules        = UserBase + "/schedules"
	UserSchedulesBase    = UserSchedules + "/base"
	UserSchedulesWeekly  = UserSchedules + "/weekly"
	UserSchedulesYearly  = UserSchedules + "/yearly"
	UserSchedulesWeekday = UserSchedulesWeekly + "/:weekday"
	UserSchedulesYearday = UserSchedulesYearly + "/:yearday"
)

func setupRoutes(s *HTTPServer) {
	store := s.Store

	s.GET("/", templates.Show(templates.Index))

	s.GET("/sign-in", templates.Show(templates.SignIn))
	s.POST("/sign-in", handles.Auth(handles.SignIn(sessionsStore), t.Auth(t.FormCredentialer), s.Store))
	s.GET("/register", templates.Show(templates.Register))
	s.POST("/register", handles.RegisterHandle(s.Store))

	s.GET(UserCalendar, UserAuth(handles.UserCalendar, store))
	s.GET(UserEvents, UserAuth(handles.UserEvents, store))
	s.GET(UserTasks, UserAuth(handles.UserTasks, store))
	s.GET(UserRoutines, UserAuth(handles.UserRoutines, store))
	s.GET(UserSchedules, UserAuth(handles.UserSchedules, store))
	s.GET(UserSchedulesBase, UserAuth(handles.UserSchedulesBase, store))
	s.GET(UserSchedulesWeekly, UserAuth(handles.UserSchedulesWeekly, store))
	s.GET(UserSchedulesYearly, UserAuth(handles.UserSchedulesYearly, store))
	s.GET(UserSchedulesWeekday, UserAuth(handles.UserSchedulesWeekday, store))
	s.GET(UserSchedulesYearday, UserAuth(handles.UserSchedulesYearday, store))

	s.ServeFiles("/img/*filepath", http.Dir(templates.ImgDir))
	s.ServeFiles("/css/*filepath", http.Dir(templates.CSSDir))
}
