package misc

import "testing"

func TestCapitalize(t *testing.T) {
	// Kinda a meme test but ok lol
	x := "scaffold"
	y := "Scaffold"
	z := "123"
	if Capitalize(x) != y {
		t.Errorf("Wrong Capitalization. Expected: %s | Got: %s", y, Capitalize(x))
	}

	if Capitalize(y) == y {
		t.Errorf("Wrong Capitalization. Expected: %s | Got: %s", y, Capitalize(y))
	}
	if Capitalize(z) == z {
		t.Errorf("Wrong Capitalization. Expected: %s | Got: %s", z, Capitalize(z))
	}

}
