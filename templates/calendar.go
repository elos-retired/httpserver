package templates

import (
	"net/http"

	"github.com/elos/httpserver/views"
	"github.com/elos/models"
	"github.com/elos/transfer"
)

func RenderCalendar(c *transfer.HTTPConnection) {
	renderTemplate(c, UserCalendar, views.MakeCalendarWeek(c.Access, c.Client().(models.User)))
}

func RenderFakeCalendar(w http.ResponseWriter, r *http.Request) {
	renderTemplate(transfer.NewHTTPConnection(w, r, nil), UserCalendar, &views.CalendarWeek{
		Days: []*views.CalendarDay{
			&views.CalendarDay{
				Header: "Header 1",
				Fixtures: []*views.CalendarFixture{
					&views.CalendarFixture{
						Name:      "Fixture 1",
						RelStart:  50,
						RelHeight: 10,
					},
					&views.CalendarFixture{
						Name:      "Fixture 2",
						RelStart:  60,
						RelHeight: 20,
					},
				},
			},
			&views.CalendarDay{
				Header: "Header 2",
				Fixtures: []*views.CalendarFixture{
					&views.CalendarFixture{
						Name:      "Fixture 1",
						RelStart:  20,
						RelHeight: 5,
					},
					&views.CalendarFixture{
						Name:      "Fixture 2",
						RelStart:  80,
						RelHeight: 20,
					},
				},
			},
			&views.CalendarDay{
				Header: "Header 3",
			},
			&views.CalendarDay{
				Header: "Header 4",
			},
			&views.CalendarDay{
				Header: "Header 5",
			},
		},
	})
}
