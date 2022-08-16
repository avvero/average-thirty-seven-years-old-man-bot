package main

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

func randomUpTo(max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max)
}
