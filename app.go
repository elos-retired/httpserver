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

func UserTemplate(f handles.TemplateHandle, s data.Store) httprouter.Handle {
	return UserAuth(handles.Template(f), s)
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

	UserSchedulesBaseAddFixture = UserSchedulesBase + "/add_fixture"
)

func setupRoutes(s *HTTPServer) {
	store := s.Store

	s.GET("/", templates.Show(templates.Index))

	s.GET("/sign-in", templates.Show(templates.SignIn))
	s.POST("/sign-in", handles.Auth(handles.SignIn(sessionsStore, UserCalendar), t.Auth(t.FormCredentialer), s.Store))
	s.GET("/register", templates.Show(templates.Register))
	s.POST("/register", handles.RegisterHandle(s.Store))

	s.GET(UserCalendar, UserTemplate(templates.RenderUserCalendar, store))
	s.GET(UserEvents, UserTemplate(templates.RenderUserEvents, store))
	s.GET(UserTasks, UserTemplate(templates.RenderUserTasks, store))
	s.GET(UserRoutines, UserTemplate(templates.RenderUserRoutines, store))
	s.GET(UserSchedules, UserTemplate(templates.RenderUserSchedules, store))
	s.GET(UserSchedulesBase, UserTemplate(templates.RenderUserSchedulesBase, store))
	s.GET(UserSchedulesWeekly, UserTemplate(templates.RenderUserSchedulesWeekly, store))
	s.GET(UserSchedulesYearly, UserTemplate(templates.RenderUserSchedulesYearly, store))

	s.GET(UserSchedulesWeekday, UserAuth(handles.UserSchedulesWeekday, store))
	s.GET(UserSchedulesYearday, UserAuth(handles.UserSchedulesYearday, store))

	s.ServeFiles("/img/*filepath", http.Dir(templates.ImgDir))
	s.ServeFiles("/css/*filepath", http.Dir(templates.CSSDir))
}
