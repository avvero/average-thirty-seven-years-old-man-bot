package main

import "strings"

func decision(text string) (respond bool, response string) {
	if strings.EqualFold(text, "gg") {
		return true, "gg"
	}
	if strings.EqualFold(text, "нет") {
		return true, "пидора ответ"
	}
	return false, ""
}
