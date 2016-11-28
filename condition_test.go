package main

import "testing"

func TestCondition(t *testing.T) {
	cases := []struct {
		In  Condition
		Out bool // Valid ?
	}{
		{In: Condition{Day: 0, Then: "dosomething"}, Out: false},
	}

	for _, c := range cases {
		v := c.In.Validate()

		if v.Error != "" {
			t.Error("Error on validate Condition: ", v.Errors)
		}
	}
}

func TestLogicValdiation(t *testing.T) {
	cases := []struct {
		In  Logic
		Out bool // Valid ?
	}{
		{In: Logic{}, Out: false},
		{In: Logic{Field: "zip", Condition: "abr"}, Out: false},
	}

	for _, c := range cases {
		v := c.In.Validate()
		valid := v.Error == ""

		if valid != c.Out {
			t.Error("Error on validate Logic", v.Errors)
		}
	}
}
