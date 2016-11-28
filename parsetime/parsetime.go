package parsetime

import "time"

// ref: https://bl.ocks.org/joyrexus/a56717634a672dcdfd48
func parsedTime(minutes int) time.Time {
	h := int(minutes / 60)
	m := minutes % 60
	t := time.Date(2000, 1, 1, h, m, 0, 0, time.UTC)

	return t
}

// Hour based on integer value for an hour
// Table to work with proper format values:
// https://andrey.nering.com.br/2015/how-to-format-date-and-time-with-go-lang/
func Hour(hour int, ampm string) string {

	t := parsedTime(hour)
	format := "15"
	pm := ampm == "PM"

	// We need this `hack` to convert d.Opens = 30 to `00:30` instead of `12:30`
	if pm {
		format = "03"
	}

	return t.Format(format)
}

// Minute based on hour int
func Minute(minute int) string {
	t := parsedTime(minute)
	return t.Format("04")
}

// AmPm based on our int value
func AmPm(hour int) string {
	t := parsedTime(hour)
	return t.Format("PM")
}
