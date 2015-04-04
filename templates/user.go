package templates

import (
	"github.com/elos/httpserver/views"
	"github.com/elos/models"
	"github.com/elos/models/calendar"
	"github.com/elos/models/user"
	"github.com/elos/transfer"
)

func RenderUserCalendar(c *transfer.HTTPConnection) error {
	RenderFakeCalendar(c.ResponseWriter(), c.Request())
	return nil
}

func RenderUserEvents(c *transfer.HTTPConnection) error {
	return Render(c, UserEvents, c.Client().(models.User))
}

func RenderUserTasks(c *transfer.HTTPConnection) error {
	return Render(c, UserTasks, c.Client().(models.User))
}

func RenderUserRoutines(c *transfer.HTTPConnection) error {
	return Render(c, UserRoutines, c.Client().(models.User))
}

func RenderUserSchedules(c *transfer.HTTPConnection) error {
	return Render(c, UserSchedules, c.Client().(models.User))
}

func RenderUserSchedulesBase(c *transfer.HTTPConnection) error {
	sv, err := userSchedulesBaseView(c)
	if err != nil {
		return err
	}

	return Render(c, UserSchedulesBase, sv)
}

func userSchedulesBaseView(c *transfer.HTTPConnection) (*ScheduleView, error) {
	u := c.Client().(models.User)
	a := c.Access

	cal, err := u.Calendar(a)
	if err != nil {
		if err == models.ErrEmptyRelationship {
			if err = user.NewCalendar(a, u); err != nil {
				return nil, NewServerError(err)
			}
		} else {
			return nil, NewServerError(err)
		}
	}

	sch, err := cal.BaseSchedule(a)
	if err != nil {
		if err == models.ErrEmptyRelationship {
			if err = calendar.NewBaseSchedule(a, cal); err != nil {
				return nil, NewServerError(err)
			}
		} else {
			return nil, NewServerError(err)
		}
	}

	fixtures, err := sch.Fixtures(a)
	if err != nil {
		return nil, NewServerError(err)
	}

	return &ScheduleView{
		Fixtures: viewFixtures(fixtures),
	}, nil
}

func RenderUserSchedulesBaseAddFixture(c *transfer.HTTPConnection) error {
	sv, err := userSchedulesBaseView(c)
	if err != nil {
		return err
	}

	return Render(c, UserSchedulesBaseAddFixture, sv)
}

func RenderUserSchedulesWeekly(c *transfer.HTTPConnection) error {
	return Render(c, UserSchedulesWeekly, c.Client().(models.User))
}

func RenderUserSchedulesYearly(c *transfer.HTTPConnection) error {
	return Render(c, UserSchedulesYearly, c.Client().(models.User))
}

func RenderUserSchedulesWeekday(c *transfer.HTTPConnection, weekday int) error {
	return Render(c, UserSchedulesWeekday, c.Client().(models.User))
}

func RenderUserSchedulesYearday(c *transfer.HTTPConnection, yearday int) error {
	return Render(c, UserSchedulesYearday, c.Client().(models.User))
}

func viewFixtures(fs []models.Fixture) []*views.CalendarFixture {
	calfs := make([]*views.CalendarFixture, len(fs))

	for i := range fs {
		calfs[i] = views.MakeCalendarFixture(fs[i])
	}

	return calfs
}

type ScheduleView struct {
	SelectedFixture models.Fixture
	Fixtures        []*views.CalendarFixture
}
