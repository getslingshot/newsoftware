package main

import (
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestExceptionValidation(t *testing.T) {
	cases := []struct {
		In  Exception
		Out string
	}{
		{
			In: Exception{
				ID:         bson.NewObjectId(),
				StartAt:    time.Now(),
				EndAt:      time.Now(),
				Type:       "foo",
				Day:        0,
				Daytime:    true,
				AfterHours: false,
				Note:       "",
				Hours:      make([]Day, 0),
			},
			Out: "01",
		},
	}

	for i, c := range cases {
		valid, _ := c.In.Validate()

		if !valid {
			t.Errorf("%v -> Fail to validate object, expected: %v, got: %v", i, c.Out, valid)
		}
	}
}

func TestExceptionWeekday(t *testing.T) {
	cases := []struct {
		In  Exception
		Out string
	}{
		{Exception{Day: 0}, "Sunday"},
		{Exception{Day: 1}, "Monday"},
		{Exception{Day: 2}, "Tuesday"},
		{Exception{Day: 3}, "Wednesday"},
		{Exception{Day: 4}, "Thursday"},
		{Exception{Day: 5}, "Friday"},
		{Exception{Day: 6}, "Saturday"},
	}

	for i, c := range cases {
		r := c.In.Weekday()

		if r != c.Out {
			t.Errorf("%v -> Fail to convert weekday, expect: %v, got: %v", i, c.Out, r)
		}
	}
}
