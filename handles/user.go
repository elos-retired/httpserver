package handles

import (
	"net/http"

	"github.com/elos/data"
	"github.com/elos/httpserver/templates"
	"github.com/elos/models"
	"github.com/julienschmidt/httprouter"
)

func UserCalendar(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	templates.RenderUserCalendar(w, r, a, a.Client().(models.User))
}

func UserEvents(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	templates.RenderUserEvents(w, r, a, a.Client().(models.User))
}

func UserTasks(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	templates.RenderUserTasks(w, r, a, a.Client().(models.User))
}

func UserRoutines(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	templates.RenderUserRoutines(w, r, a, a.Client().(models.User))
}

func UserSchedules(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	templates.RenderUserSchedules(w, r, a, a.Client().(models.User))
}

func UserSchedulesBase(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	templates.RenderUserSchedulesBase(w, r, a, a.Client().(models.User))
}

func UserSchedulesWeekly(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	templates.RenderUserSchedulesWeekly(w, r, a, a.Client().(models.User))
}

func UserSchedulesYearly(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	templates.RenderUserSchedulesYearly(w, r, a, a.Client().(models.User))
}

func UserSchedulesWeekday(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	templates.RenderUserSchedulesWeekday(w, r, a, a.Client().(models.User))
}

func UserSchedulesYearday(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	templates.RenderUserSchedulesYearday(w, r, a, a.Client().(models.User))
}
