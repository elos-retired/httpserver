package templates

import (
	"go/build"
	"log"
	"path/filepath"
	"text/template"
)

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

var layoutTemplate string = "layout.tmpl"

func Prepend(s string, v ...string) []string {
	l := make([]string, len(v)+1)
	l[0] = s
	for i := range v {
		l[i+1] = v[i]
	}
	return l
}

func Layout(v ...string) []string {
	return Prepend(layoutTemplate, v...)
}

var sessionsLayoutTemplate string = "sessions/layout.tmpl"

func Sessions(v ...string) []string {
	return Layout(Prepend(sessionsLayoutTemplate, v...)...)
}

var templateSets = map[Name][]string{
	Index: Layout("index.html"),

	SignIn:         Sessions("sessions/sign-in.tmpl"),
	Register:       Sessions("sessions/register.tmpl"),
	AccountCreated: Sessions("sessions/account-created.tmpl"),

	UserCalendar:         Layout("user/calendar.tmpl"),
	UserEvents:           Layout("user/events.tmpl"),
	UserTasks:            Layout("user/tasks.tmpl"),
	UserRoutines:         Layout("user/routines.tmpl"),
	UserSchedules:        Layout("user/schedules.tmpl"),
	UserSchedulesBase:    Layout("user/schedules/layout.tmpl", "user/schedules/base.tmpl"),
	UserSchedulesWeekly:  Layout("user/schedules/weekly.tmpl"),
	UserSchedulesYearly:  Layout("user/schedules/yearly.tmpl"),
	UserSchedulesWeekday: Layout("user/schedules/weekday.tmpl"),
	UserSchedulesYearday: Layout("user/schedules/yearday.tmpl"),

	UserSchedulesBaseAddFixture: Layout("user/schedules/layout.tmpl", "user/schedules/base-add.tmpl"),
}

var templates = map[Name]*template.Template{}

func init() {
	if err := parseHTMLTemplates(templateSets); err != nil {
		log.Fatal(err)
	}
}

func joinDir(base string, files []string) []string {
	r := make([]string, len(files))
	for i := range files {
		r[i] = filepath.Join(TemplatesDir, files[i])
	}
	return r
}

func parseHTMLTemplates(sets map[Name][]string) error {
	for name, set := range sets {
		t, err := template.ParseFiles(joinDir(TemplatesDir, set)...)
		if err != nil {
			return err
		}
		/*
			t = t.Lookup("ROOT")
			if t == nil {
				return fmt.Errorf("ROOT template not found in %v", set)
			}
		*/
		templates[name] = t
	}
	return nil
}
