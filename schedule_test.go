package main

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestOpenScheduleHour(t *testing.T) {
	cases := []struct {
		In  Schedule
		Out string
	}{
		{Schedule{Opens: 60}, "01"},
		{Schedule{Opens: 120}, "02"},
		{Schedule{Opens: 60 * 18}, "06"},
	}

	for i, c := range cases {
		r := c.In.OpenHour()

		if r != c.Out {
			t.Errorf("%v -> Fail to convert open hour, expect: %v, got: %v", i, c.Out, r)
		}
	}
}

func TestOpenScheduleMinute(t *testing.T) {
	cases := []struct {
		In  Schedule
		Out string
	}{
		{Schedule{Opens: 60}, "00"},
		{Schedule{Opens: 120}, "00"},
		{Schedule{Opens: 30}, "30"},
		{Schedule{Opens: 61}, "01"},
		{Schedule{Opens: 70}, "10"},
		{Schedule{Opens: 18 * 60}, "00"},
	}

	for i, c := range cases {
		r := c.In.OpenMinute()

		if r != c.Out {
			t.Errorf("%v -> Fail to parse minute, expect: %v, got: %v", i, c.Out, r)
		}
	}
}

func TestCloseScheduleHour(t *testing.T) {
	cases := []struct {
		In  Schedule
		Out string
	}{
		{Schedule{Closes: 60}, "01"},
		{Schedule{Closes: 120}, "02"},
		{Schedule{Closes: 18 * 60}, "06"},
	}

	for i, c := range cases {
		r := c.In.CloseHour()

		if r != c.Out {
			t.Errorf("%v -> Fail to convert close hour, expect: %v, got: %v", i, c.Out, r)
		}
	}
}

func TestCloseScheduleMinute(t *testing.T) {
	cases := []struct {
		In  Schedule
		Out string
	}{
		{Schedule{Closes: 60}, "00"},
		{Schedule{Closes: 120}, "00"},
		{Schedule{Closes: 30}, "30"},
		{Schedule{Closes: 61}, "01"},
		{Schedule{Closes: 70}, "10"},
		{Schedule{Closes: 18 * 60}, "00"},
	}

	for i, c := range cases {
		r := c.In.CloseMinute()

		if r != c.Out {
			t.Errorf("%v -> Fail to parse minute, expect: %v, got: %v", i, c.Out, r)
		}
	}
}

func TestCloseScheduleAmPm(t *testing.T) {
	cases := []struct {
		In  Schedule
		Out string
	}{
		{Schedule{Closes: 60}, "AM"},
		{Schedule{Closes: 120}, "AM"},
		{Schedule{Closes: 18 * 60}, "PM"},
	}

	for i, c := range cases {
		r := c.In.CloseAmPm()

		if r != c.Out {
			t.Errorf("%v -> Fail to convert close am / pm, expect: %v, got: %v", i, c.Out, r)
		}
	}
}

func TestOpenScheduleAmPm(t *testing.T) {
	cases := []struct {
		In  Schedule
		Out string
	}{
		{Schedule{Opens: 60}, "AM"},
		{Schedule{Opens: 2 * 60}, "AM"},
		{Schedule{Opens: 18 * 60}, "PM"},
	}

	for i, c := range cases {
		r := c.In.OpenAmPm()

		if r != c.Out {
			t.Errorf("%v -> Fail to convert Open am / pm, expect: %v, got: %v", i, c.Out, r)
		}
	}
}

func TestScheduleWeekday(t *testing.T) {
	cases := []struct {
		In  Schedule
		Out string
	}{
		{Schedule{Day: 0}, "Sunday"},
		{Schedule{Day: 1}, "Monday"},
		{Schedule{Day: 2}, "Tuesday"},
		{Schedule{Day: 3}, "Wednesday"},
		{Schedule{Day: 4}, "Thursday"},
		{Schedule{Day: 5}, "Friday"},
		{Schedule{Day: 6}, "Saturday"},
	}

	for i, c := range cases {
		r := c.In.Weekday()

		if r != c.Out {
			t.Errorf("%v -> Fail to convert weekday, expect: %v, got: %v", i, c.Out, r)
		}
	}
}

func TestSchedulePublic(t *testing.T) {
	c := Schedule{}
	r := c.Public()

	v, _ := r["dayDisplay"]

	if v != "Sunday" {
		t.Errorf("Fail to get public day info, expect: %v, got: %v", "Sunday", v)
	}
}

func TestScheduleValidate(t *testing.T) {

	cases := []struct {
		In  Schedule
		Out bool // is invalid ?
	}{
		{Schedule{ID: bson.NewObjectId(), Day: 0, Opens: 0, Closes: 0}, true},
		{Schedule{ID: bson.NewObjectId(), Day: 1, Opens: 0, Closes: 2}, false},
		{Schedule{ID: bson.NewObjectId(), Day: 7, Opens: 0, Closes: 0}, true},
		{Schedule{ID: bson.NewObjectId(), Day: -1, Opens: 0, Closes: 0}, true},
		{Schedule{ID: bson.NewObjectId(), Day: 1, Opens: -1, Closes: 0}, true},
		{Schedule{ID: bson.NewObjectId(), Day: 1, Opens: 1, Closes: -1}, true},
		{Schedule{ID: bson.NewObjectId(), Day: 1, Opens: 1, Closes: 2}, false},
		{Schedule{ID: bson.NewObjectId(), Day: 1, Opens: 1, Closes: 1}, true},
	}

	for i, c := range cases {
		invalid, errs := c.In.Validate()

		if invalid != c.Out {
			t.Errorf("%v -> Fail to validate datetime, expect: %v, got: %v", i, c.Out, invalid)
			t.Error(errs)
		}
	}
}
