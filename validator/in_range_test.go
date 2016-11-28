package validator_test

import "testing"
import "bitbucket.org/tekkismatt/ams_slingshot_app/validator"

func TestInRange(t *testing.T) {
	cases := []struct {
		In  int
		Out bool
	}{
		{0, false},
		{-1, false},
		{1, true},
		{52, true},
		{53, false},
	}

	for i, c := range cases {
		v := validator.NewValidator()
		actual := v.InRange("foo", 1, 52, c.In)

		if actual != c.Out {
			t.Errorf("%v, Expected validator to be %v, but got %v", i, c.Out, actual)
		}
	}
}
