package validator

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/tekkismatt/ams_slingshot_app/parsetime"
)

// Validator struct based on:
type Validator struct {
	Valid  bool
	Errors map[string]string
}

// NewValidator constructor
func NewValidator() *Validator {
	return &Validator{Errors: make(map[string]string), Valid: true}
}

// Present Validate presence of a field
func (v *Validator) Present(field string, value string) bool {
	if value == "" {
		v.Valid = false
		v.Errors[field] = "Required"
	}

	return v.Valid
}

// Weekday Validate presence of a field
func (v *Validator) Weekday(field string, start time.Weekday, end time.Weekday, value time.Weekday) bool {
	if value < start || value > end {
		v.Valid = false
		v.Errors[field] = "Invalid value, 0-6 range allowed"
	}

	return v.Valid
}

// Positive Validate integer positive value
func (v *Validator) Positive(field string, value int) bool {
	if value < 0 {
		v.Valid = false
		v.Errors[field] = "must be greater than 0"
	}

	return v.Valid
}

// GreaterThanInt64 Validate integer GreaterThan value
func (v *Validator) GreaterThanInt64(field string, base int64, value int64, errmsg string) bool {
	if value > base {
		v.Valid = false
		v.Errors[field] = errmsg
	}

	return v.Valid
}

// AfterHour validates given value should be greater then a provided number
func (v *Validator) AfterHour(field string, after int, value int) bool {
	if value <= after {
		v.Valid = false
		// This error message works because we always want format like `07:05`
		v.Errors[field] = "must be after " + parsetime.Hour(after, "PM") + ":" + parsetime.Minute(after)
	}

	return v.Valid
}

// InRangeHours validates given value should be greater then a provided number
func (v *Validator) InRangeHours(field string, from int, to int, value int) bool {
	if from < value && value > to {
		v.Valid = false
		v.Errors[field] = "must be in range " + parsetime.Hour(from, "PM") + " to " + strconv.Itoa(to)
	}

	return v.Valid
}

// Email Validate presence of a field
// http://stackoverflow.com/questions/46155/Validate-email-address-in-javascript
func (v *Validator) Email(field string, value string) bool {
	Re := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

	if !Re.MatchString(value) {
		v.Valid = false
		v.Errors[field] = "InValid email format"
	}

	return v.Valid
}

// IsValid to check if we have any Errors arrociated to v
func (v *Validator) IsValid() bool {
	return !v.Valid
}

func stringInArray(arr []string, val string) bool {
	for _, s := range arr {
		if s == val {
			return true
		}
	}

	return false
}

// InArray validate presence of a field
func (v *Validator) InArray(field string, arr []string, value string) bool {
	if !stringInArray(arr, value) {
		v.Valid = false
		v.Errors[field] = "Type should be one of: " + strings.Join(arr, ", ")
	}

	return v.Valid
}
