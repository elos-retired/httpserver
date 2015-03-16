package templates

import (
	"net/http"

	"github.com/elos/data"
	"github.com/elos/models"
)

func RenderUserCalendar(w http.ResponseWriter, r *http.Request, a data.Access, u models.User) {
	RenderFakeCalendar(w, r)
}

func RenderUserEvents(w http.ResponseWriter, r *http.Request, a data.Access, u models.User) {
	Render(w, r, UserEvents, u)
}

func RenderUserTasks(w http.ResponseWriter, r *http.Request, a data.Access, u models.User) {
	Render(w, r, UserTasks, u)
}

func RenderUserRoutines(w http.ResponseWriter, r *http.Request, a data.Access, u models.User) {
	Render(w, r, UserRoutines, u)
}

func RenderUserSchedules(w http.ResponseWriter, r *http.Request, a data.Access, u models.User) {
	Render(w, r, UserSchedules, u)
}

func RenderUserSchedulesBase(w http.ResponseWriter, r *http.Request, a data.Access, u models.User) {
	Render(w, r, UserSchedulesBase, u)
}

func RenderUserSchedulesWeekly(w http.ResponseWriter, r *http.Request, a data.Access, u models.User) {
	Render(w, r, UserSchedulesWeekly, u)
}

func RenderUserSchedulesYearly(w http.ResponseWriter, r *http.Request, a data.Access, u models.User) {
	Render(w, r, UserSchedulesYearly, u)
}

func RenderUserSchedulesWeekday(w http.ResponseWriter, r *http.Request, a data.Access, u models.User) {
	Render(w, r, UserSchedulesWeekday, u)
}

func RenderUserSchedulesYearday(w http.ResponseWriter, r *http.Request, a data.Access, u models.User) {
	Render(w, r, UserSchedulesYearday, u)
}
