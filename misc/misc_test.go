package misc

import "testing"

func TestCapitalize(t *testing.T) {
	x := "name"
	y := "Name"
	z := "123"

	if capitalizedX := Capitalize(x); capitalizedX != y {
		t.Errorf("Wrong Capitalization. Expected: %s | Got: %s", y, capitalizedX)
	}

	if capitalizedY := Capitalize(y); capitalizedY != y {
		t.Errorf("Wrong Capitalization. Expected: %s | Got: %s", y, capitalizedY)
	}

	if capitalizedZ := Capitalize(z); capitalizedZ != z {
		t.Errorf("Wrong Capitalization. Expected: %s | Got: %s", z, capitalizedZ)
	}
}
