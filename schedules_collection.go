// http://nerdyworm.com/blog/2013/05/15/sorting-a-slice-of-structs-in-go/ ftw!
package main

import "sort"

// Schedules collection
type Schedules []Schedule

// Len of an schedule collection
func (slice Schedules) Len() int {
	return len(slice)
}

// Less of schedule collection
func (slice Schedules) Less(i, j int) bool {
	return slice[i].Opens < slice[j].Opens
}

// Swap method for schedules based on schedule.Opens
func (slice Schedules) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// LogTable used to display schedule
func (slice Schedules) LogTable() {
	logger.Println()

	for i, d := range slice {
		logger.Printf("|%-2v|%-9s|%-6s|%-6s|\n", i, d.Weekday(), d.OpenHour(), d.CloseHour())
	}
}

// SchedulesBy is the type of a "less" function that defines the ordering of its Schedule arguments.
type SchedulesBy func(p1, p2 *Schedule) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by SchedulesBy) Sort(schedules []Schedule) {
	ps := &scheduleSorter{
		schedules: schedules,
		by:        by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}

	sort.Sort(ps)
}

// scheduleSorter joins a By function and a slice of Schedules to be sorted.
type scheduleSorter struct {
	schedules []Schedule
	by        func(p1, p2 *Schedule) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *scheduleSorter) Len() int {
	return len(s.schedules)
}

// Swap is part of sort.Interface.
func (s *scheduleSorter) Swap(i, j int) {
	s.schedules[i], s.schedules[j] = s.schedules[j], s.schedules[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *scheduleSorter) Less(i, j int) bool {
	return s.by(&s.schedules[i], &s.schedules[j])
}
