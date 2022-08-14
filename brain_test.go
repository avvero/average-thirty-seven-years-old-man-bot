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

func Test_returnsOnGG(t *testing.T) {
	data := map[string]string{
		"gg": "gg",
		"GG": "gg",
		"gG": "gg",
		"Gg": "gg",
	}
	for k, expected := range data {
		respond, response := decision(0, k)
		if !respond || response != expected {
			t.Error("Expected and got:", expected, " != ", response)
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
	for k, expected := range data {
		respond, response := decision(0, k)
		if !respond || response != expected {
			t.Error("Expected and got:", expected, " != ", response)
		}
	}
}

func Test_returnsOnMorrowind(t *testing.T) {
	data := []string{
		"morrowind",
		"моровинд",
		"морровинд",
		"морровинд",
		"бла бла бла morrowind",
		"бла бла бла morrowind бла бла бла ",
	}
	expected := "Morrowind - одна из лучших игр эва"
	for _, text := range data {
		respond, response := decision(0, text)
		if !respond || response != expected {
			t.Error("Expected and got:", expected, " != ", response)
		}
	}
}

func Test_returnsOnElderRing(t *testing.T) {
	data := []string{
		"Elden Ring",
		"elden ring",
		"бла бла бла Elden Ring бла бла бла",
		"elden ring",
		"бла бла бла elden ring бла бла бла",
		"ER",
		"бла бла бла ER бла бла бла",
		"ЕР",
		"бла бла бла ЕР бла бла бла",
		"ЭР",
		"бла бла бла ЭР бла бла бла",
	}
	expected := "Elden Ring - это величие"
	for _, text := range data {
		respond, response := decision(0, text)
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

func Test_returnsOnBuy(t *testing.T) {
	data := []string{
		"купил",
		"бла бла бла купил",
		"бла бла бла купил бла бла бла",
	}
	expected := "А не пиздишь? Аренда это не покупка"
	for _, text := range data {
		respond, response := decision(0, text)
		if !respond || response != expected {
			t.Error("Expected and got:", expected, " != ", response)
		}
	}
}

func Test_returnsOnSpotify(t *testing.T) {
	data := []string{
		"spotify",
		"Spotify",
		"спотифай",
		"бла бла бла spotify бла бла бла",
		"бла бла бла спотифай бла бла бла",
	}
	expected := "Эти пидоры Антону косарик должны за подписку"
	for _, text := range data {
		respond, response := decision(0, text)
		if !respond || response != expected {
			t.Error("Expected and got:", expected, " != ", response)
		}
	}
}

func Test_returnsForLucky(t *testing.T) {
	expected := "хуйню не неси"
	respond := false
	response := ""
	for i := 0; i < 500; i++ {
		thisRespond, thisResponse := decision(0, "any")
		if thisRespond {
			respond = thisRespond
			response = thisResponse
		}
	}
	if !respond || response != expected {
		t.Error("Expected and got:", expected, " != ", response)
	}
}

func Test_returnsOnDevops(t *testing.T) {
	data := []string{
		"devops",
		"девопс",
	}
	expected := "Девопсы не нужны"
	for _, text := range data {
		respond, response := decision(0, text)
		if !respond || response != expected {
			t.Error("Expected and got:", expected, " != ", response)
		}
	}
}
