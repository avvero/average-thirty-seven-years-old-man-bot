package main

import (
	"strconv"
	"strings"
)

func decision(chatId int64, text string) (respond bool, response string) {
	if !Contains([]string{"0", "-1001733786877", "245851441", "-578279468"}, strconv.FormatInt(chatId, 10)) {
		return true, "Хули нада, пес?"
	}
	if strings.EqualFold(text, "gg") {
		return true, "gg"
	}
	if strings.EqualFold(text, "нет") {
		return true, "пидора ответ"
	}
	return false, ""
}
