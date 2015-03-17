package handles

import (
	"net/http"
	"strconv"

	"github.com/elos/data"
	"github.com/elos/httpserver/templates"
	"github.com/elos/transfer"
	"github.com/julienschmidt/httprouter"
)

type TemplateHandle func(*transfer.HTTPConnection) error

func Template(t TemplateHandle) AccessHandle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
		c := transfer.NewHTTPConnection(w, r, a)
		templates.CatchError(c, t(c))
	}
}

func UserSchedulesWeekday(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	c := transfer.NewHTTPConnection(w, r, a)
	CatchError(c, userSchedulesWeekday(c, p))
}

func userSchedulesWeekday(c *transfer.HTTPConnection, p httprouter.Params) error {
	weekday, err := strconv.Atoi(p.ByName("weekday"))
	if err != nil {
		return NewMissingParamError("weekday")
	}

	if weekday < 0 || weekday > 6 {
		return NewBadParamError("weekday", "must be in range 0-6 inclusive")
	}

	templates.CatchError(c, templates.RenderUserSchedulesWeekday(c, weekday))
	return nil
}

func UserSchedulesYearday(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	c := transfer.NewHTTPConnection(w, r, a)
	userSchedulesYearday(c, p)
}

func userSchedulesYearday(c *transfer.HTTPConnection, p httprouter.Params) error {
	yearday, err := strconv.Atoi(p.ByName("yearday"))
	if err != nil {
		return NewMissingParamError("yearday")
	}

	if yearday < 0 || yearday > 1231 {
		return NewBadParamError("yearday", "must at least be in range 0-1231 inclusive to be potentially valid")
	}

	templates.CatchError(c, templates.RenderUserSchedulesYearday(c, yearday))
	return nil
}
