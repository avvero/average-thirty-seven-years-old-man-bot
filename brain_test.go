package main

import (
	"strconv"
	"testing"
)

func Test_responseOnlyToWhitelisted(t *testing.T) {
	for _, chatIdString := range []string{"-1001733786877", "245851441", "-578279468"} {
		chatId, _ := strconv.ParseInt(chatIdString, 10, 64)
		respond, response := decision(chatId, "gg")
		if !respond || response != "gg" {
			t.Error("Expected: gg, got", response)
		}
	}
	for _, chatIdString := range []string{"123"} {
		chatId, _ := strconv.ParseInt(chatIdString, 10, 64)
		respond, response := decision(chatId, "gg")
		if !respond || response != "Хули нада, пес?" {
			t.Error("Expected: Хули нада, пес?, got", chatId, response)
		}
	}
}

func Test_returnsOnGG(t *testing.T) {
	for _, text := range []string{"gg", "GG"} {
		respond, response := decision(0, text)
		if !respond || response != "gg" {
			t.Error("Expected: gg, got", response)
		}
	}
}

func Test_returnsOnNet(t *testing.T) {
	for _, text := range []string{"нет", "Нет", "НЕТ"} {
		respond, response := decision(0, text)
		if !respond || response != "пидора ответ" {
			t.Error("Expected: пидора ответ, got", response)
		}
	}
}
