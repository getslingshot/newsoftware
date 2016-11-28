package main

import (
	"sort"
	"strconv"
)

func validateRangeHours(coll Days) (bool, map[string]string) {
	sort.Sort(coll)
	e := sendError()

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
			e.Errors["Range: "+strconv.Itoa(i+1)] = "Invalid range for: " + m
		}
	}

	return len(e.Errors) > 0, e.Errors
}
