package main

import (
	"net/http"

	"time"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var hours = []string{
	"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12",
}

var minutes = []string{
	"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10",
	"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
	"21", "22", "23", "24", "25", "26", "27", "28", "29", "30",
	"31", "32", "33", "34", "35", "36", "37", "38", "39", "40",
	"41", "42", "43", "44", "45", "46", "47", "48", "49", "50",
	"51", "52", "53", "54", "55", "56", "57", "58", "59",
}

var days = []time.Weekday{0, 1, 2, 3, 4, 5, 6}

// SchedulingLogicHandler http handler
func SchedulingLogicHandler(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("schedules")
	query := bson.M{}
	var schedules Schedules

	if err := c.Find(query).Sort("day").All(&schedules); err != nil {
		respond(w, r, http.StatusInternalServerError, "tmpl/content/schedulingLogic.tmpl", nil, nil)
		return
	}

	var exceptions Exceptions
	if err := db.C("exceptions").Find(bson.M{}).All(&exceptions); err != nil {
		respond(w, r, http.StatusInternalServerError, "tmpl/content/schedulingLogic.tmpl", nil, nil)
		return
	}

	data := Schedules{}
	for _, d := range schedules {
		d.Selected = true
		data = append(data, d)
	}

	// Get days added to schedules
	existingDays := make(map[int]int, len(days))
	for _, s := range data {
		existingDays[int(s.Day)] = int(s.Day)
	}

	// Let's check for non existing days on schedule.
	// Add an 'empty' record so user can add range of hours to that weekday
	for _, d := range days {
		_, ok := existingDays[int(d)]

		if !ok {
			day := Schedule{Day: d, Selected: false}
			data = append(data, day)
		}
	}

	daySort := func(a, b *Schedule) bool {
		return a.Day < b.Day
	}

	// Let's sort so we display content from Sunday -> Saturday
	SchedulesBy(daySort).Sort(data)

	e := Exception{StartAt: time.Now()}
	exceptions = append(exceptions, e)

	exceptionSort := func(a, b *Exception) bool {
		return a.StartAt.Before(b.StartAt)
	}

	// Let's sort so we display exceptions by date
	ExceptionsBy(exceptionSort).Sort(exceptions)

	var collection Calls
	c = db.C("calls")

	if err := c.Find(bson.M{}).All(&collection); err != nil {
		respondJSON(w, r, http.StatusInternalServerError, e)
		return
	}

	calls := make([]map[string]interface{}, len(collection))

	for i, c := range collection {
		calls[i] = c.Public()
	}

	var conditionCollection Conditions
	c = db.C("conditions")

	if err := c.Find(bson.M{}).All(&conditionCollection); err != nil {
		e := sendError()
		e.Error = "Cant process request"
		respondJSON(w, r, http.StatusInternalServerError, e)
		return
	}

	conditions := make([]map[string]interface{}, len(conditionCollection))

	for i, c := range conditionCollection {
		conditions[i] = c.Public()
	}

	d := M{
		"schedules":  data,
		"hours":      hours,
		"minutes":    minutes,
		"exceptions": exceptions,
		"calls":      calls,
		"conditions": conditions,
	}

	respond(w, r, http.StatusOK, "tmpl/content/schedulingLogic.tmpl", d, nil)
}
