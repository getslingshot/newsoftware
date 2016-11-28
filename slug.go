package main

import "strings"

func slug(s string) string {
	trim := strings.Trim(s, " ")
	return strings.ToLower(strings.Join(strings.Split(trim, " "), "-"))
}
