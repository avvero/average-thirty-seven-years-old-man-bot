package main

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func decision(chatId int64, text string) (respond bool, response string) {
	// Lucky if random returns 0
	if rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100) == 0 {
		return true, "хуйню не неси"
	}
	//
	text = strings.ToLower(text)
	if !Contains([]string{"0", "-1001733786877", "245851441", "-578279468"}, strconv.FormatInt(chatId, 10)) {
		return true, "Mr Moony presents his compliments to Professor Snape, and begs him to keep his abnormally large nose out of other people’s business."
	}
	if text == "gg" {
		return true, "gg"
	}
	if text == "нет" {
		return true, "пидора ответ"
	}
	if strings.Contains(text, "morrowind") ||
		strings.Contains(text, "морровинд") ||
		strings.Contains(text, "моровинд") {
		return true, "Morrowind - одна из лучших игр эва"
	}
	if text == "er" ||
		text == "ер" ||
		strings.Contains(text, "elden ring") ||
		strings.Contains(text, " er ") ||
		strings.Contains(text, " ер ") {
		return true, "Elden Ring - это величие"
	}
	if strings.Contains(text, "купил") {
		return true, "А не пиздишь? Аренда это не покупка"
	}
	if strings.Contains(text, "spotify") || strings.Contains(text, "спотифай") {
		return true, "Эти пидоры Антону косарик должны за подписку"
	}
	return false, ""
}

func senselessPhrases() []string {
	return []string{"хуйню не неси", "база"}
}
