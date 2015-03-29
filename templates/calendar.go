package templates

import (
	"fmt"
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
	Time      string
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

func hourString(t time.Time) string {
	var suffix string
	hour := t.Hour() + 1
	if hour <= 12 {
		suffix = "AM"
	} else {
		suffix = "PM"
	}

	var extraZero string
	if t.Minute()%60 < 10 {
		extraZero = "0"
	}

	return fmt.Sprintf("%d:%s%d %s", t.Hour()%12, extraZero, t.Minute()%60, suffix)
}

func timeString(f models.Fixture) string {
	return fmt.Sprintf("%s - %s", hourString(f.StartTime()), hourString(f.EndTime()))
}

func calendarFixture(f models.Fixture) *CalendarFixture {
	startM := f.StartTime().Hour()*60 + f.StartTime().Minute()
	endM := f.EndTime().Hour()*60 + f.EndTime().Minute()
	return &CalendarFixture{
		Name:      f.Name(),
		Time:      timeString(f),
		RelStart:  float32(startM) / 1440 * 100,
		RelHeight: float32(endM-startM) / 1440 * 100,
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
		sched, err := cal.YeardaySchedule(a, now)
		if err != nil {
			return cw
		}
		cw.Days[i] = calendarDay(a, sched, now)
		now.Add(24 * time.Hour)
	}
	return cw
}
