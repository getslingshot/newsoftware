package main

import (
	"time"

	"bitbucket.org/tekkismatt/ams_slingshot_app/validator"

	"gopkg.in/mgo.v2/bson"
)

// Call model
type Call struct {
	ID     bson.ObjectId `bson:"_id" json:"_id"`
	OnDay  time.Weekday  `bson:"onday" json:"onday"`
	ForDay time.Weekday  `bson:"forday" json:"forday"`
	Week   int           `bson:"week" json:"week"`
	Hours  []Day         `bson:"hours" json:"hours"`
	Note   string        `bson:"note" json:"note"`
}

// Validate Call model
func (d Call) Validate() (bool, map[string]string) {
	v := validator.NewValidator()

	// Can start with 0 up to 1439 minutes
	v.Weekday("onday", 0, 6, d.OnDay)
	v.Weekday("forday", 0, 6, d.ForDay)

	v.InRange("week", 1, 52, d.Week)

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

	return v.IsValid(), v.Errors
}

// Public fields for Call
func (d Call) Public() map[string]interface{} {
	hours := make([]map[string]interface{}, len(d.Hours))
	for i, h := range d.Hours {
		hours[i] = h.Public()
	}
	return map[string]interface{}{
		"_id":    d.ID,
		"forDay": d.ForDay,
		"hours":  hours,
		"note":   d.Note,
		"onDay":  d.OnDay,
		"week":   d.Week,
	}
}

// OnWeekday Weekday return the actual day, i.e: 'Sunday' of the week
func (d Call) OnWeekday() string {
	return time.Weekday(d.OnDay).String()
}

func (d Call) forWeekday() string {
	return time.Weekday(d.ForDay).String()
}

// RangeHour virtual struct
type RangeHour struct {
	Opens  int
	Closes int
}
