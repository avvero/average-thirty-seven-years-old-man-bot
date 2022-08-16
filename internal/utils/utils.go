package utils

import (
	"math/rand"
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

func RandomUpTo(max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max)
}
