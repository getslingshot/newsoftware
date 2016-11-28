package main

import "time"

func validateScheduleRange(collection Schedules) (bool, map[string]string) {
	var sundays Schedules
	var mondays Schedules
	var tuesdays Schedules
	var wednesdays Schedules
	var thursdays Schedules
	var fridays Schedules
	var saturdays Schedules

	e := sendError()

	for _, s := range collection {
		switch s.Day {

		case time.Sunday:
			sundays = append(sundays, s)

		case time.Monday:
			mondays = append(mondays, s)

		case time.Tuesday:
			tuesdays = append(tuesdays, s)

		case time.Wednesday:
			wednesdays = append(wednesdays, s)

		case time.Thursday:
			thursdays = append(thursdays, s)

		case time.Friday:
			fridays = append(fridays, s)

		case time.Saturday:
			saturdays = append(saturdays, s)

		}
	}

	matrix := []Schedules{
		sundays,
		mondays,
		tuesdays,
		wednesdays,
		thursdays,
		fridays,
		saturdays,
	}

	for _, coll := range matrix {
		if len(coll) == 0 {
			continue
		}

		for i, v := range coll {
			if i == 0 {
				// skip index 0, we will validate index 1 against 0, index 2 against 1...
				continue
			}

			lastRangeCloses := coll[i-1].Closes
			next := len(coll) > i

			if !next {
				continue
			}

			if v.Opens <= lastRangeCloses {
				e.Error = "Invalid range(s)"
				m := v.OpenHour() + ":" + v.OpenMinute() + v.OpenAmPm() +
					" to " + v.CloseHour() + ":" + v.CloseMinute() + v.CloseAmPm()
				e.Errors[v.Weekday()] = "Invalid range for: " + m
			}
		}
	}

	return e.Error == "", e.Errors
}
