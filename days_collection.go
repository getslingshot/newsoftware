package main

import "sort"

// Days collection
type Days []Day

// Len of an day collection
func (slice Days) Len() int {
	return len(slice)
}

// Less of day collection
func (slice Days) Less(i, j int) bool {
	return slice[i].Opens < slice[j].Opens
}

// Swap method for days based on day.Opens
func (slice Days) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// DaysBy is the type of a "less" function that defines the ordering of its Day arguments.
type DaysBy func(p1, p2 *Day) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by DaysBy) Sort(days []Day) {
	ps := &daySorter{
		days: days,
		by:   by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}

	sort.Sort(ps)
}

// daySorter joins a By function and a slice of Days to be sorted.
type daySorter struct {
	days []Day
	by   func(p1, p2 *Day) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *daySorter) Len() int {
	return len(s.days)
}

// Swap is part of sort.Interface.
func (s *daySorter) Swap(i, j int) {
	s.days[i], s.days[j] = s.days[j], s.days[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *daySorter) Less(i, j int) bool {
	return s.by(&s.days[i], &s.days[j])
}
