package templates

import (
	"go/build"
	"path/filepath"
)

/*
	PackagePath finds the full path for the specified
	golang import path

	i.e. PackagePath("github.com/elos/httpserver/templates")
	     -> "~/Nick/workspace/go/src/github.com/elos/httpserver/templates"
			or some equivlent root path
*/
func PackagePath(importPath string) string {
	p, err := build.Default.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

var (
	httpserverPath = "github.com/elos/httpserver"
	AssetsDir      = filepath.Join(PackagePath(httpserverPath), "assets")
	TemplatesDir   = filepath.Join(AssetsDir, "templates")
	ImgDir         = filepath.Join(AssetsDir, "img")
	CSSDir         = filepath.Join(AssetsDir, "css")
)

type Name int

const (
	Index Name = iota
	SignIn
	Register
	AccountCreated

	UserCalendar
	UserEvents
	UserTasks
	UserRoutines
	UserSchedules
	UserSchedulesBase
	UserSchedulesWeekly
	UserSchedulesYearly
	UserSchedulesWeekday
	UserSchedulesYearday
	UserSchedulesBaseAddFixture
)

var (
	layoutTemplate          string = "layout.tmpl"
	sessionsLayoutTemplate  string = "sessions/layout.tmpl"
	schedulesLayoutTemplate string = "user/schedules/layout.tmpl"
)

/*
	Prepend creates a slice of strings from variadic
	arguments with the guarantee that the slice will be
	of size >= 1, with index 0 equal to s

	Prepend is useful for constructing templateSets
	i.e.,
		func Root(v ...string) []string {
			return Prepend("root.tmpl", v...)
		}
*/
func Prepend(s string, v ...string) []string {
	l := make([]string, len(v)+1)
	l[0] = s
	for i := range v {
		l[i+1] = v[i]
	}
	return l
}

// Layout prepends variadic arguments with the layoutTemplate
func Layout(v ...string) []string {
	return Prepend(layoutTemplate, v...)
}

// Sessions prepends variadic arguments with the layout and sessions templates
func Sessions(v ...string) []string {
	return Layout(Prepend(sessionsLayoutTemplate, v...)...)
}

// Schedules prepends variadic arguments with the layout and schedules templates
func Schedules(v ...string) []string {
	return Layout(Prepend(schedulesLayoutTemplate, Prepend("user/schedules/common.tmpl", v...)...)...)
}

// Definition of the available templateSets for elos
// used in initialization of the templates, see: init.go
var templateSets = map[Name][]string{
	Index: Layout("index.html"),

	SignIn:         Sessions("sessions/sign-in.tmpl"),
	Register:       Sessions("sessions/register.tmpl"),
	AccountCreated: Sessions("sessions/account-created.tmpl"),

	UserCalendar:  Layout("user/schedules/common.tmpl", "user/calendar.tmpl"),
	UserEvents:    Layout("user/events.tmpl"),
	UserTasks:     Layout("user/tasks.tmpl"),
	UserRoutines:  Layout("user/routines.tmpl"),
	UserSchedules: Layout("user/schedules.tmpl"),

	UserSchedulesBase:           Schedules("user/schedules/base.tmpl"),
	UserSchedulesBaseAddFixture: Schedules("user/schedules/base-add.tmpl"),

	UserSchedulesWeekly:  Layout("user/schedules/weekly.tmpl"),
	UserSchedulesYearly:  Layout("user/schedules/yearly.tmpl"),
	UserSchedulesWeekday: Layout("user/schedules/weekday.tmpl"),
	UserSchedulesYearday: Layout("user/schedules/yearday.tmpl"),
}
