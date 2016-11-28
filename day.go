package main

import (
	"bitbucket.org/tekkismatt/ams_slingshot_app/parsetime"
	"bitbucket.org/tekkismatt/ams_slingshot_app/validator"
)

// Day model
type Day struct {
	Opens  int `bson:"opens" json:"opens"`
	Closes int `bson:"closes" json:"closes"`
}

// Public interface for Day objects
func (d Day) Public() map[string]interface{} {
	return map[string]interface{}{
		"opens":       d.Opens,
		"closes":      d.Closes,
		"openHour":    d.OpenHour(),
		"openMinute":  d.OpenMinute(),
		"openAmPm":    d.OpenAmPm(),
		"closeHour":   d.CloseHour(),
		"closeMinute": d.CloseMinute(),
		"closeAmPm":   d.CloseAmPm(),
	}
}

// Validate Day model
func (d Day) Validate() (bool, map[string]string) {
	v := validator.NewValidator()
	v.InRangeHours("opens", -1, 1439, d.Opens)
	v.Positive("opens", d.Opens)

	// Can start from 1 (since 0 will be used as opens) up to 1440, the end of the day
	v.InRangeHours("closes", 1, 1441, d.Closes)
	v.AfterHour("closes", d.Opens, d.Closes)
	v.Positive("closes", d.Closes)

	return v.IsValid(), v.Errors
}

// OpenHour returns format like: '05'
func (d Day) OpenHour() string {
	return parsetime.Hour(d.Opens, d.OpenAmPm())
}

// OpenMinute returns format like: '05'
func (d Day) OpenMinute() string {
	return parsetime.Minute(d.Opens)
}

// CloseHour returns format like: '05'
func (d Day) CloseHour() string {
	return parsetime.Hour(d.Closes, d.CloseAmPm())
}

// CloseMinute returns format like: '05'
func (d Day) CloseMinute() string {
	return parsetime.Minute(d.Closes)
}

// CloseAmPm returns AM or PM
func (d Day) CloseAmPm() string {
	return parsetime.AmPm(d.Closes)
}

// OpenAmPm returns AM or PM
func (d Day) OpenAmPm() string {
	return parsetime.AmPm(d.Opens)
}
