package main

import "testing"

func Test_returnsOnGG(t *testing.T) {
	for _, text := range []string{"gg", "GG"} {
		respond, response := decision(text)
		if !respond || response != "gg" {
			t.Error("Expected: gg, got", response)
		}
	}
}

func Test_returnsOnNet(t *testing.T) {
	for _, text := range []string{"нет", "Нет", "НЕТ"} {
		respond, response := decision(text)
		if !respond || response != "пидора ответ" {
			t.Error("Expected: пидора ответ, got", response)
		}
	}
}
