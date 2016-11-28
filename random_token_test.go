package main

import "testing"

func TestRandomToken(t *testing.T) {
	for _, key := range cases {
		result := len(randomToken(key.in))

		if result != key.out {
			t.Errorf("Expected %v to be %v but got %v", key.in, key.out, result)
		}
	}
}

var cases = []struct {
	in  int
	out int
}{
	{1, 1},
	{2, 2},
	{10, 10},
	{15, 15},
	{200, 200},
}
