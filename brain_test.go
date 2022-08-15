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
		"1":              "Mr Moony presents his compliments to Professor Snape, and begs him to keep his abnormally large nose out of other people’s business.",
	}
	for k, expected := range data {
		chatId, _ := strconv.ParseInt(k, 10, 64)
		respond, response := decision(chatId, "gg")
		if !respond || response != expected {
			t.Error("Expected and got:", expected, " != ", response)
		}
	}
}

func Test_returnsOnSomeText(t *testing.T) {
	data := map[string]string{
		"gg": "gg",
		"GG": "gg",
		"gG": "gg",
		"Gg": "gg",

		"нет": "пидора ответ",
		"НЕТ": "пидора ответ",
		"Нет": "пидора ответ",
		"НеТ": "пидора ответ",
		"Heт": "пидора ответ",
		"Нeт": "пидора ответ",

		"morrowind":             "Morrowind - одна из лучших игр эва",
		"моровинд":              "Morrowind - одна из лучших игр эва",
		"морровинд":             "Morrowind - одна из лучших игр эва",
		"бла бла бла morrowind": "Morrowind - одна из лучших игр эва",
		"бла бла бла morrowind бла бла бла": "Morrowind - одна из лучших игр эва",

		"Elden Ring": "Elden Ring - это величие",
		"бла бла бла Elden Ring бла бла бла": "Elden Ring - это величие",
		"elden ring": "Elden Ring - это величие",
		"бла бла бла elden ring бла бла бла": "Elden Ring - это величие",
		"ER": "Elden Ring - это величие",
		"бла бла бла ER бла бла бла": "Elden Ring - это величие",
		"ЕР": "Elden Ring - это величие",
		"бла бла бла ЕР бла бла бла": "Elden Ring - это величие",
		"ЭР": "Elden Ring - это величие",
		"бла бла бла ЭР бла бла бла": "Elden Ring - это величие",

		"купил":             "А не пиздишь? Аренда это не покупка",
		"бла бла бла купил": "А не пиздишь? Аренда это не покупка",
		"бла бла бла купил бла бла бла": "А не пиздишь? Аренда это не покупка",

		"spotify":  "Эти пидоры Антону косарик должны за подписку",
		"Spotify":  "Эти пидоры Антону косарик должны за подписку",
		"спотифай": "Эти пидоры Антону косарик должны за подписку",
		"бла бла бла spotify бла бла бла":  "Эти пидоры Антону косарик должны за подписку",
		"бла бла бла спотифай бла бла бла": "Эти пидоры Антону косарик должны за подписку",

		"devops": "Девопсы не нужны",
		"девопс": "Девопсы не нужны",
	}
	for k, expected := range data {
		respond, response := decision(0, k)
		if !respond || response != expected {
			t.Error("Expected and got:", expected, " != ", response)
		}
	}
}

func Test_returnsOnNotElderRing(t *testing.T) {
	data := []string{
		"pERt",
		"sdfERdfd",
		"аааЕРваа",
	}
	for _, text := range data {
		respond, response := decision(0, text)
		if respond {
			t.Error("Not expected: ", response)
		}
	}
}

func Test_returnsForLucky(t *testing.T) {
	respond := false
	response := ""
	for i := 0; i < 500; i++ {
		thisRespond, thisResponse := decision(0, "any")
		if thisRespond {
			respond = thisRespond
			response = thisResponse
		}
	}
	if !respond || !Contains(senselessPhrases, response) {
		t.Error("Expected and got: something from senselessPhrases != ", response)
	}
}
