package utils

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"
)

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func ContainsRune(s []rune, e rune) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func RandomUpTo(max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max)
}

func PrintJson(value any) string {
	js, _ := json.Marshal(value)
	return string(js)
}

func PrettyPrint(value any) string {
	result := PrintJson(value)
	result = strings.ReplaceAll(result, "\"", "")
	//result = strings.ReplaceAll(result, "{", ",")
	//result = strings.ReplaceAll(result, "}", ",")
	return result
}
