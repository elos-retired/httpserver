package handles

import (
	"net/http"
	"strconv"
	"time"

	"github.com/elos/data"
	"github.com/elos/httpserver/templates"
	"github.com/elos/models"
	"github.com/elos/models/fixture"
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
	CatchError(c, userSchedulesYearday(c, p))
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

func UserSchedulesBaseAddFixture(w http.ResponseWriter, r *http.Request, p httprouter.Params, a data.Access) {
	c := transfer.NewHTTPConnection(w, r, a)
	CatchError(c, userSchedulesBaseAddFixture(c, a))
}

var formTimeLayout = "15:04"

func userSchedulesBaseAddFixture(c *transfer.HTTPConnection, a data.Access) error {
	r := c.Request()

	params, err := getVals(r, "name", "start_time", "end_time")
	if err != nil {
		return err
	}

	/*
		label, err := strconv.ParseBool(params["label"])
		if err != nil {
			return NewBadParamError("label", err.Error())
		}
	*/

	start_time, err := time.Parse(formTimeLayout, params["start_time"])
	if err != nil {
		return NewBadParamError("start_time", err.Error())
	}
	end_time, err := time.Parse(formTimeLayout, params["end_time"])
	if err != nil {
		return NewBadParamError("end_time", err.Error())
	}

	cal, err := c.Client().(models.User).Calendar(a)
	if err != nil {
		return err
	}

	s, err := cal.BaseSchedule(a)
	if err != nil {
		return err
	}

	f, err := fixture.New(a)
	if err != nil {
		return err
	}

	f.SetName(params["name"])
	f.SetStartTime(start_time)
	f.SetEndTime(end_time)
	// f.SetLabel(label)

	if err = a.Save(f); err != nil {
		return err
	}

	s.IncludeFixture(f)

	a.Save(f)
	a.Save(s)

	http.Redirect(c.ResponseWriter(), c.Request(), "/user/schedules/base", http.StatusFound)
	return nil
}

func getVals(r *http.Request, v ...string) (map[string]string, error) {
	params := make(map[string]string)

	for _, v := range v {
		s := r.FormValue(v)
		if s == "" {
			return nil, NewMissingParamError(v)
		}
		params[v] = s
	}

	return params, nil
}
