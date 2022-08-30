package utils

import (
	"testing"
)

func Test_PrintJson(t *testing.T) {
	data := map[string]any{
		"{\"key\":\"value\"}": map[string]string{"key": "value"},
		"{\"key\":\"\"}":      map[string]string{"key": ""},
		"":                    nil,
	}
	for expected, object := range data {
		result := PrintJson(object)
		if result != expected {
			t.Error("PrintJson for ", object, ": ", expected, " != ", result)
		}
	}
}
