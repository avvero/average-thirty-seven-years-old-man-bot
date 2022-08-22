package brain

import (
	"testing"
)

func _Test_normalization(t *testing.T) {
	censor := &Censor{}
	data := map[string]string{
		"gg":        "gg",
		"GG":        "gg",
		"корова":    "корова",
		"коорооваа": "корова",
		"коoрoоваа": "корова", // seconds o's is eng
	}
	for value, expected := range data {
		respond, response := censor.Check(0, value)
		if !respond || response != expected {
			t.Errorf("Expected: \"%s\" but got: \"%s\"", expected, response)
		}
	}
}
