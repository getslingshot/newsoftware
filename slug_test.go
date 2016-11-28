package main

import "testing"

var tableTest = []struct {
	in  string
	out string
}{
	{"my awesome company", "my-awesome-company"},
	{"MY awesome company", "my-awesome-company"},
	{"mY awesome company", "my-awesome-company"},
	{"my company", "my-company"},
	{"123 MY COMPANY ", "123-my-company"},
}

func TestSlug(t *testing.T) {
	for _, key := range tableTest {
		result := slug(key.in)

		if result != key.out {
			t.Error("Expected " + key.in + " to be: " + key.out + "but got: " + result)
		}
	}
}
