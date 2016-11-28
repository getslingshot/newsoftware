package validator

import "strconv"

// InRange Validate presence of a field
func (v *Validator) InRange(field string, start int, end int, value int) bool {
	if value < start || value > end {
		v.Valid = false
		v.Errors[field] = "Invalid value, " + strconv.Itoa(start) + "-" + strconv.Itoa(end) + "  range allowed"
	}

	return v.Valid
}
