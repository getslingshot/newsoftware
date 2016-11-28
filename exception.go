package main

import (
	"time"

	"bitbucket.org/tekkismatt/ams_slingshot_app/validator"

	"gopkg.in/mgo.v2/bson"
)

// Exception model
// Based on https://godoc.org/time#Weekday
// and http://schema.org/OpeningHoursSpecification
// and http://stackoverflow.com/questions/17460235/mongodb-opening-hours-schema-and-query-for-open-closed
// Day will use 0-6 -> Sunday - Saturday ( defined on golang and js and...)
// Opens | Closes uses minutes to define start | end
// Opens at 8:00am -> 8*60 minutes = 480
type Exception struct {
	ID         bson.ObjectId `bson:"_id" json:"_id"`
	StartAt    time.Time     `bson:"startAt" json:"startAt"`
	EndAt      time.Time     `bson:"endAt" json:"endAt"`
	Type       string        `bson:"type" json:"type"`
	Day        time.Weekday  `bson:"day" json:"day"`
	Daytime    bool          `bson:"daytime" json:"daytime"`
	AfterHours bool          `bson:"afterhours" json:"afterhours"`
	Note       string        `bson:"note" json:"note"`
	Hours      []Day         `bson:"hours" json:"hours"`
}

// DisplayStart helper method for rendering
func (d Exception) DisplayStart() string {
	month := d.StartAt.Month()
	day := d.StartAt.Day()
	year := d.StartAt.Year()

	t := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return t.Format("01/02")
}

// DisplayEnd helper method for rendering
func (d Exception) DisplayEnd() string {
	month := d.EndAt.Month()
	day := d.EndAt.Day()
	year := d.EndAt.Year()

	t := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return t.Format("01/02")
}

// Validate Exception model
func (d Exception) Validate() (bool, map[string]string) {
	v := validator.NewValidator()
	v.Weekday("day", 0, 6, d.Day)
	v.GreaterThanInt64("startAt", d.EndAt.Unix(), d.StartAt.Unix(), "Must ends after "+d.StartAt.String())

	// Validate hours
	for _, h := range d.Hours {
		// Each hour is a model. We can call Validate and set `v` properties
		invalid, errs := h.Validate()

		// If current (loop) hour is invalid, we set v.Valid to false and
		// attach any error from hour validation to main exception validations
		if invalid {
			v.Valid = false
			for key, value := range errs {
				v.Errors[key] = value
			}
		}
	}

	invalid, errs := validateRangeHours(d.Hours)

	if invalid {
		v.Valid = false

		for key, value := range errs {
			v.Errors[key] = value
		}
	}

	return v.Valid, v.Errors
}

// Public fields for Exception
func (d Exception) Public() map[string]interface{} {
	return map[string]interface{}{
		"day":        d.Day,
		"note":       d.Note,
		"dayDisplay": d.Day.String(),
		"hours":      d.Hours,
	}
}

// Weekday return the actual day, i.e: 'Sunday' of the week
func (d Exception) Weekday() string {
	return time.Weekday(d.Day).String()
}

// StartAtStr returns startAt string
func (d Exception) StartAtStr() string {
	return d.StartAt.Format("01/02/2006")
}

// EndAtStr returns startAt string
func (d Exception) EndAtStr() string {
	return d.EndAt.Format("01/02/2006")
}
