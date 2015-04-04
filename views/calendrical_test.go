package views

import (
	"testing"
	"time"

	"github.com/elos/data"
	"github.com/elos/models"
	"github.com/elos/models/fixture"
	"github.com/elos/models/persistence"
	"github.com/elos/models/schedule"
	"github.com/elos/testing/expect"
)

var (
	now      = time.Now()
	testTime = time.Date(2015, 1, 1, 9, 0, 0, 0, time.UTC)
)

func TestCalendarHeader(t *testing.T) {
	t.Parallel()

	testTime := time.Now()

	if CalendarHeader(testTime) != TodayHeader {
		t.Errorf("Today should have a header of: %s", TodayHeader)
	}

	testTime = testTime.Add(24 * time.Hour)

	if CalendarHeader(testTime) != TomorrowHeader {
		t.Errorf("Tomorrow should have a header of %s, got: %s", TomorrowHeader, CalendarHeader(testTime))
	}

	testTime = testTime.Add(24 * time.Hour)

	if CalendarHeader(testTime) != testTime.Weekday().String() {
		t.Errorf("Random day should have weekday as its header, got: %s", CalendarHeader(testTime))
	}
}

func TestFormattedHour(t *testing.T) {
	t.Parallel()

	testTime := time.Date(2015, 1, 1, 9, 0, 0, 0, time.UTC)

	formattedTimes := map[time.Time]string{
		testTime:                                  "9:00 AM",
		testTime.Add(1 * time.Minute):             "9:01 AM",
		testTime.Add(3*time.Hour + 1*time.Minute): "12:01 PM",
	}

	for ti, f := range formattedTimes {
		if FormattedHour(ti) != f {
			t.Errorf("Expected: %s, but got: %s", f, FormattedHour(ti))
		}
	}
}

type timeable struct {
	s time.Time
	e time.Time
}

func (t timeable) StartTime() time.Time {
	return t.s
}

func (t timeable) SetStartTime(s time.Time) {
	t.s = s
}

func (t timeable) EndTime() time.Time {
	return t.e
}

func (t timeable) SetEndTime(e time.Time) {
	t.e = e
}

func TestFormattedTimeable(t *testing.T) {
	t.Parallel()

	testTime := time.Date(2015, 1, 1, 9, 0, 0, 0, time.UTC)
	timeable := timeable{testTime, testTime.Add(1 * time.Hour)}

	expectedString := "9:00 AM - 10:00 AM"

	if FormattedTimeable(timeable) != expectedString {
		t.Errorf("Expected %s, but got: %s", expectedString, FormattedTimeable(timeable))
	}
}

func TestAbsMinute(t *testing.T) {
	t.Parallel()

	timeTable := map[time.Time]int{
		testTime:                                   9 * 60,
		testTime.Add(2 * time.Minute):              9*60 + 2,
		testTime.Add(15*time.Hour - 1*time.Minute): 1439,
		testTime.Add(15 * time.Hour):               0,
	}

	for ti, min := range timeTable {
		if AbsMinute(ti) != min {
			t.Errorf("Expected %d, but got %d", min, AbsMinute(ti))
		}
	}
}

func TestRelativeStartPosition(t *testing.T) {
	t.Parallel()

	s := persistence.Store(persistence.MongoMemoryDB())
	f, err := fixture.New(s)
	expect.NoError("creating new fixture", err, t)

	f.SetStartTime(testTime)
	f.SetEndTime(testTime.Add(1 * time.Hour))

	if RelativeStartPosition(f) != .375 {
		t.Errorf("Expected .375, but got %d", RelativeStartPosition(f))
	}
}

func TestRelativeHeight(t *testing.T) {
	t.Parallel()

	s := persistence.Store(persistence.MongoMemoryDB())
	f, err := fixture.New(s)
	expect.NoError("creating new fixture", err, t)

	f.SetStartTime(testTime)
	f.SetEndTime(testTime.Add(1 * time.Hour))

	expectedHeight := 1.0 / 24.0

	if RelativeHeight(f) != float32(expectedHeight) {
		t.Errorf("Expected %d, but got %d", expectedHeight, RelativeHeight(f))
	}
}

func TestMakeCalendarFixture(t *testing.T) {
	t.Parallel()

	s := persistence.Store(persistence.MongoMemoryDB())
	f, err := fixture.New(s)
	expect.NoError("creating new fixture", err, t)

	testName := "testName"
	f.SetName(testName)

	f.SetStartTime(testTime)
	f.SetEndTime(testTime.Add(1 * time.Hour))

	cf := MakeCalendarFixture(f)

	if f.Name() != cf.Name {
		t.Errorf("Fixture and CalendarFixture names should match")
	}

	if cf.Time != FormattedTimeable(f) {
		t.Errorf("CalendarFixture's Time should be the formatted timeable of the fixture")
	}

	if cf.RelStart != RelativeStartPosition(f)*100 {
		t.Errorf("CalendarFixture's RelStart should be 100 times its RelativeStartPosition decimal")
	}

	if cf.RelHeight != RelativeHeight(f)*100 {
		t.Errorf("CalendarFixture's RelHeight should be 100 times its RelativeHeight decimal")
	}
}

func TestMakeCalendarDay(t *testing.T) {
	t.Parallel()

	s := persistence.Store(persistence.MongoMemoryDB())

	sched, err := schedule.New(s)
	expect.NoError("creating new schedule", err, t)

	fs := make([]models.Fixture, 3)
	for i := 0; i < 3; i++ {
		f, err := fixture.New(s)
		expect.NoError("creating new fixture", err, t)

		fs[i] = f

		err = sched.IncludeFixture(f)
		expect.NoError("including fixture in schedule", err, t)

		err = s.Save(f)
		expect.NoError("saving fixture", err, t)
	}

	err = s.Save(sched)
	expect.NoError("saving schedule", err, t)

	access := data.NewAnonAccess(s)

	cd, err := MakeCalendarDay(access, sched)
	expect.NoError("making calendar day", err, t)

	if len(cd.Fixtures) != 3 {
		t.Errorf("Expected there to be three fixtures on this schedule")
	}

	if cd.Header != CalendarHeader(sched.StartTime()) {
		t.Errorf("Expected calendar header to be %s", CalendarHeader(sched.StartTime()))
	}
}

func TestMakeCalendarWeek(t *testing.T) {
	t.Parallel()

	//	s := persistence.Store(persistence.MongoMemoryDB())

}

func TestMakeSchedule(t *testing.T) {
	t.Parallel()

	s := persistence.Store(persistence.MongoMemoryDB())

	sched, err := schedule.Create(s)
	expect.NoError("creating schedule", err, t)
	f1, err := fixture.Create(s)
	expect.NoError("creating fixture", err, t)
	f2, err := fixture.Create(s)
	expect.NoError("creating fixture", err, t)
	err = sched.IncludeFixture(f1)
	expect.NoError("including fixture", err, t)
	err = sched.IncludeFixture(f2)
	expect.NoError("including fixture", err, t)

	schedView, err := MakeSchedule(data.NewAnonAccess(s), sched)
	expect.NoError("making schedule", err, t)

	if len(schedView.Fixtures) != 2 {
		t.Errorf("Schedule should generate correct number of fixtures, in this case 2")
	}
}
