package main

import "sort"

// Exceptions collection
type Exceptions []Exception

// ExceptionsBy is the type of a "less" function that defines the ordering of its Exception arguments.
type ExceptionsBy func(p1, p2 *Exception) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by ExceptionsBy) Sort(exceptions []Exception) {
	ps := &exceptionSorter{
		exceptions: exceptions,
		by:         by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}

	sort.Sort(ps)
}

// exceptionSorter joins a By function and a slice of Exceptions to be sorted.
type exceptionSorter struct {
	exceptions []Exception
	by         func(p1, p2 *Exception) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *exceptionSorter) Len() int {
	return len(s.exceptions)
}

// Swap is part of sort.Interface.
func (s *exceptionSorter) Swap(i, j int) {
	s.exceptions[i], s.exceptions[j] = s.exceptions[j], s.exceptions[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *exceptionSorter) Less(i, j int) bool {
	return s.by(&s.exceptions[i], &s.exceptions[j])
}
