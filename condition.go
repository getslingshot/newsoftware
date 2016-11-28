package main

import (
	"time"

	"bitbucket.org/tekkismatt/ams_slingshot_app/validator"

	"gopkg.in/mgo.v2/bson"
)

// Condition model
type Condition struct {
	ID           bson.ObjectId  `bson:"_id" json:"_id"`
	Day          time.Weekday   `bson:"day" json:"day"`
	Logic        []Logic        `bson:"logic" json:"logic"`
	Then         string         `bson:"then" json:"then"`
	ScheduleDays []time.Weekday `bson:"scheduleDays" json:"scheduleDays"`
	Hours        []Day          `bson:"hours" json:"hours"`
	Note         string         `bson:"note" json:"note"`
}

// Logic used to store conditional entries
// Field like 'zipcode'
// Condition like `equals`, `contains`...
// Start, End For range, greater_than, less_than
// Value String value for some conditions
type Logic struct {
	Field     string `bson:"field" json:"field"`
	Condition string `bson:"condition" json:"condition"`
	Value     string `bson:"value" json:"value"`
	Start     int    `bson:"start" json:"start"`
	End       int    `bson:"end" json:"end"`
}

var conditionTypes = []string{
	"contains",
	"range",
	"equal",
	"greater_than",
	"less_than",
}

// Public method for Logic struct
func (l Logic) Public() map[string]interface{} {
	m := map[string]interface{}{
		"field":     l.Field,
		"condition": l.Condition,
		"value":     l.Value,
		"start":     l.Start,
		"end":       l.End,
	}

	return m
}

// Validate Condition model
func (c Condition) Validate() E {
	var e E
	v := validator.NewValidator()
	v.Present("then", c.Then)
	v.Weekday("day", 0, 6, c.Day)

	// Validate hours
	for _, l := range c.Logic {
		// Each l is a model. We can call Validate and set `v` properties
		err := l.Validate()

		if len(err.Errors) > 0 {
			e.Error = "Invalid payload"

			v.Valid = false
			for key, value := range err.Errors {
				v.Errors[key] = value
			}
		}
	}

	if !v.Valid {
		e.Error = "Invalid"
		e.Errors = v.Errors
	}
	return e
}

// Validate Logic model
func (l Logic) Validate() E {
	var e E
	v := validator.NewValidator()

	v.Present("field", l.Field)
	v.Present("condition", l.Condition)
	v.InArray("condition", conditionTypes, l.Condition)

	switch l.Condition {
	case "range", "less_than", "greater_than":
		v.Positive("start", l.Start)
		v.Positive("end", l.Start)

	case "contains":
		v.Present("value", l.Value)
	}

	if !v.Valid {
		e.Error = "Invalid"
		e.Errors = v.Errors
	}

	return e
}

// Public method for Condition model
func (c Condition) Public() map[string]interface{} {
	logics := make([]map[string]interface{}, len(c.Logic))

	for i, l := range c.Logic {
		logics[i] = l.Public()
	}

	hours := make([]map[string]interface{}, len(c.Hours))
	for i, h := range c.Hours {
		hours[i] = h.Public()
	}

	m := map[string]interface{}{
		"_id":          c.ID.Hex(),
		"day":          c.Day,
		"scheduleDays": c.ScheduleDays,
		"hours":        hours,
		"logic":        logics,
		"then":         c.Then,
		"note":         c.Note,
	}

	return m
}
