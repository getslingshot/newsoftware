package main

import (
	"time"

	"bitbucket.org/tekkismatt/ams_slingshot_app/parsetime"
	"bitbucket.org/tekkismatt/ams_slingshot_app/validator"

	"gopkg.in/mgo.v2/bson"
)

// Schedule model
// Based on https://godoc.org/time#Weekday
// and http://schema.org/OpeningHoursSpecification
// and http://stackoverflow.com/questions/17460235/mongodb-opening-hours-schema-and-query-for-open-closed
// Day will use 0-6 -> Sunday - Saturday ( defined on golang and js and...)
// Opens | Closes uses minutes to define start | end
// Opens at 8:00am -> 8*60 minutes = 480
type Schedule struct {
	ID         bson.ObjectId
	Day        time.Weekday `bson:"day" json:"day"`
	Opens      int          `bson:"opens" json:"opens"`
	Closes     int          `bson:"closes" json:"closes"`
	AfterHours bool         `bson:"afterhours" json:"afterhours"`
	Daytime    bool         `bson:"daytime" json:"daytime"`
	Note       string       `bson:"note" json:"note"`

	Selected bool
}

// Validate Schedule model
func (d Schedule) Validate() (bool, map[string]string) {
	v := validator.NewValidator()

	// Can start with 0 up to 1439 minutes
	v.InRangeHours("opens", 0, 1439, d.Opens)
	v.Positive("opens", d.Opens)
	v.Weekday("day", 0, 6, d.Day)

	// Can start from 1 (since 0 will be used as opens) up to 1440, the end of the day
	v.InRangeHours("closes", 0, 1440, d.Closes)
	v.AfterHour("closes", d.Opens, d.Closes)
	v.Positive("closes", d.Closes)

	return v.IsValid(), v.Errors
}

// Public fields for Schedule
func (d Schedule) Public() map[string]interface{} {
	return map[string]interface{}{
		"day":              d.Day,
		"opens":            d.Opens,
		"closes":           d.Closes,
		"note":             d.Note,
		"dayDisplay":       d.Day.String(),
		"opensDisplay":     d.OpenHour(),
		"closesDisplay":    d.CloseHour(),
		"closeampmDisplay": d.CloseAmPm(),
		"openampmDisplay":  d.OpenAmPm(),
	}
}

// Weekday return the actual day, i.e: 'Sunday' of the week
func (d Schedule) Weekday() string {
	return time.Weekday(d.Day).String()
}

// OpenHour returns format like: '05:00'
func (d Schedule) OpenHour() string {
	return parsetime.Hour(d.Opens, d.OpenAmPm())
}

// OpenMinute returns format like: '05:00'
func (d Schedule) OpenMinute() string {
	return parsetime.Minute(d.Opens)
}

// CloseHour returns format like: '05:00'
func (d Schedule) CloseHour() string {
	return parsetime.Hour(d.Closes, d.CloseAmPm())
}

// CloseMinute returns format like: '05:00'
func (d Schedule) CloseMinute() string {
	return parsetime.Minute(d.Closes)
}

// CloseAmPm returns AM or PM
func (d Schedule) CloseAmPm() string {
	return parsetime.AmPm(d.Closes)
}

// OpenAmPm returns AM or PM
func (d Schedule) OpenAmPm() string {
	return parsetime.AmPm(d.Opens)
}
