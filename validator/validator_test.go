package validator_test

import "testing"
import "bitbucket.org/tekkismatt/ams_slingshot_app/validator"

func TestValidator(t *testing.T) {
	// arrange & act
	v := validator.NewValidator()

	if !v.Valid {
		t.Errorf("expected instance to init with true, got: %v", v.Valid)
	}

	if v.Errors == nil {
		t.Errorf("expected instance to init with empty Errors map, got: %v", v.Errors)
	}
}

func TestPresent(t *testing.T) {
	v := validator.NewValidator()

	b := v.Present("foo", "")

	if b {
		t.Errorf("expect Present to validate precence of field")
	}
}

func TestEmail(t *testing.T) {
	// arrange & act
	var cases = []struct {
		in  string
		out bool
	}{
		{"foo", false},
		{"f@f.com", true},
		{"f@", false},
		{"@f.com", false},
	}

	for _, option := range cases {
		v := validator.NewValidator()
		valid := v.Email("email", option.in)

		if valid != option.out {
			t.Errorf("expect %v to be invalid email, %v, %v", option.in, valid, option.out)
		}
	}
}

func TestIsValid(t *testing.T) {
	// arrage & act
	v := validator.NewValidator()
	b := v.IsValid()

	// assert
	if b {
		t.Errorf("expect Present to validate precence of field")
	}
}
