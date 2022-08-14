package main

import (
	"strconv"
	"testing"
)

func Test_responseOnlyToWhitelisted(t *testing.T) {
	data := map[string]string{
		"-1001733786877": "gg",
		"245851441":      "gg",
		"-578279468":     "gg",
		"1":              "Хули нада, пес?",
	}
	for k, v := range data {
		chatId, _ := strconv.ParseInt(k, 10, 64)
		respond, response := decision(chatId, "gg")
		if !respond || response != v {
			t.Error("Expected and got:", v, " != ", response)
		}
	}
}

func Test_returnsOnGG(t *testing.T) {
	data := map[string]string{
		"gg": "gg",
		"GG": "gg",
		"gG": "gg",
		"Gg": "gg",
	}
	for k, v := range data {
		respond, response := decision(0, k)
		if !respond || response != v {
			t.Error("Expected and got:", v, " != ", response)
		}
	}
}

func Test_returnsOnNet(t *testing.T) {
	data := map[string]string{
		"нет": "пидора ответ",
		"НЕТ": "пидора ответ",
		"Нет": "пидора ответ",
		"НеТ": "пидора ответ",
	}
	for k, v := range data {
		respond, response := decision(0, k)
		if !respond || response != v {
			t.Error("Expected and got:", v, " != ", response)
		}
	}
}
