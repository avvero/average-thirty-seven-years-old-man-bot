package main

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func decision(chatId int64, text string) (respond bool, response string) {
	if randomUpTo(100) == 0 {
		phrase := senselessPhrases[randomUpTo(len(senselessPhrases))]
		return true, phrase
	}
	//
	text = strings.ToLower(text)
	if !Contains([]string{"0", "-1001733786877", "245851441", "-578279468"}, strconv.FormatInt(chatId, 10)) {
		return true, "Mr Moony presents his compliments to Professor Snape, and begs him to keep his abnormally large nose out of other people’s business."
	}
	if text == "gg" {
		return true, "gg"
	}
	if normalizeRu(text) == "нет" {
		return true, "пидора ответ"
	}
	if strings.Contains(text, "morrowind") ||
		strings.Contains(text, "морровинд") ||
		strings.Contains(text, "моровинд") {
		return true, "Morrowind - одна из лучших игр эва"
	}
	if text == "er" ||
		text == "ер" ||
		text == "эр" ||
		strings.Contains(text, "elden ring") ||
		strings.Contains(text, " er ") ||
		strings.Contains(text, " ер ") ||
		strings.Contains(text, " эр ") {
		return true, "Elden Ring - это величие"
	}
	if strings.Contains(text, "купил") {
		return true, "А не пиздишь? Аренда это не покупка"
	}
	if strings.Contains(text, "spotify") || strings.Contains(text, "спотифай") {
		return true, "Эти пидоры Антону косарик должны за подписку"
	}
	if strings.Contains(normalizeEn(text), "devops") ||
		strings.Contains(normalizeRu(text), "девопс") {
		return true, "Девопсы не нужны"
	}
	return false, ""
}

func randomUpTo(max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max)
}

var charMap = map[string]string{
	"e": "е",
	"o": "о",
	"h": "н",
	"a": "а",
	"t": "т",
	"k": "к",
	"c": "с",
	"b": "б",
	"m": "м",
}

func normalizeRu(text string) string {
	result := text
	for k, v := range charMap {
		result = strings.Replace(result, k, v, -1)
	}
	return result
}

func normalizeEn(text string) string {
	result := text
	for k, v := range charMap {
		result = strings.Replace(result, v, k, -1)
	}
	return result
}

var senselessPhrases = []string{"хуйню не неси", "база", "мда", "вообще похую", "ничего нового"}
