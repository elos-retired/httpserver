package templates

import (
	"log"
	"net/http"
	"time"

	"github.com/elos/data"
	"github.com/elos/models"
	"github.com/elos/models/calendar"
	"github.com/elos/transfer"
)

func RenderCalendar(c *transfer.HTTPConnection) {
	renderTemplate(c, UserCalendar, calendarWeek(c.Access, c.Client().(models.User)))
}

func RenderFakeCalendar(w http.ResponseWriter, r *http.Request) {
	renderTemplate(transfer.NewHTTPConnection(w, r, nil), UserCalendar, &CalendarWeek{
		Days: []*CalendarDay{
			&CalendarDay{
				Header: "Header 1",
				Fixtures: []*CalendarFixture{
					&CalendarFixture{
						Name:      "Fixture 1",
						RelStart:  50,
						RelHeight: 10,
					},
					&CalendarFixture{
						Name:      "Fixture 2",
						RelStart:  60,
						RelHeight: 20,
					},
				},
			},
			&CalendarDay{
				Header: "Header 2",
				Fixtures: []*CalendarFixture{
					&CalendarFixture{
						Name:      "Fixture 1",
						RelStart:  20,
						RelHeight: 5,
					},
					&CalendarFixture{
						Name:      "Fixture 2",
						RelStart:  80,
						RelHeight: 20,
					},
				},
			},
			&CalendarDay{
				Header: "Header 3",
			},
			&CalendarDay{
				Header: "Header 4",
			},
			&CalendarDay{
				Header: "Header 5",
			},
		},
	})
}

// Data structures for rendering calendar views
type CalendarWeek struct {
	Days []*CalendarDay
}

type CalendarDay struct {
	Header   string
	Fixtures []*CalendarFixture
}

type CalendarFixture struct {
	Name      string
	RelStart  float32
	RelHeight float32
}

func calHeader(t time.Time) string {
	n := time.Now()
	if t.YearDay() == n.YearDay() {
		return "Today"
	} else if t.YearDay()-n.YearDay() == 1 {
		return "Tomorrow"
	} else {
		return t.Weekday().String()
	}

}

func calendarFixture(f models.Fixture) *CalendarFixture {
	startM := f.StartTime().Hour()*60 + f.StartTime().Minute()
	endM := f.EndTime().Hour()*60 + f.EndTime().Minute()
	return &CalendarFixture{
		Name:      f.Name(),
		RelStart:  float32(startM) / 1440,
		RelHeight: float32(startM+endM) / 1440,
	}
}

func calendarDay(a data.Access, s models.Schedule, t time.Time) *CalendarDay {
	cd := new(CalendarDay)
	cd.Header = calHeader(t)
	fs, _ := s.OrderedFixtures(a)
	cd.Fixtures = make([]*CalendarFixture, len(fs))
	for i, f := range fs {
		cd.Fixtures[i] = calendarFixture(f)
	}
	return cd
}

func calendarWeek(a data.Access, u models.User) *CalendarWeek {
	cw := new(CalendarWeek)
	cw.Days = make([]*CalendarDay, 5)
	cal, err := u.Calendar(a)
	if err != nil {
		log.Print(err)
		if err == models.ErrEmptyRelationship {
			cal, _ = calendar.Create(a)
			u.SetCalendar(cal)
			a.Save(u)
			a.Save(cal)
		} else {
			return cw
		}
	}
	now := time.Now()
	for i := 0; i < 5; i++ {
		sched, err := cal.ScheduleForDay(a, now)
		if err != nil {
			return cw
		}
		cw.Days[i] = calendarDay(a, sched, now)
		now.Add(24 * time.Hour)
	}
	return cw
}
