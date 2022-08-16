package main

import (
	"strconv"
	"strings"
)

type Brain struct {
	memory *Memory
}

func NewBrain(memory *Memory) *Brain {
	brain := &Brain{
		memory: memory,
	}
	return brain
}

func (brain Brain) decision(chatId int64, text string) (respond bool, response string) {
	if randomUpTo(100) == 0 {
		phrase := brain.memory.senselessPhrases[randomUpTo(len(brain.memory.senselessPhrases))]
		return true, phrase
	}
	if len(text) > 14 && randomUpTo(100) == 0 {
		phrase := brain.khaleesify(text)
		return true, phrase
	}
	if !Contains([]string{"0", "-1001733786877", "245851441", "-578279468"}, strconv.FormatInt(chatId, 10)) {
		return true, "Mr Moony presents his compliments to Professor Snape, and begs him to keep his abnormally large nose out of other people’s business."
	}
	//
	text = strings.ToLower(text)
	if text == "gg" {
		return true, "gg"
	}
	if brain.normalizeRu(text) == "нет" {
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
	if strings.Contains(brain.normalizeEn(text), "devops") ||
		strings.Contains(brain.normalizeRu(text), "девопс") {
		return true, "Девопсы не нужны"
	}
	if text == "трансформация" ||
		text == "трансформацию" ||
		text == "трансформации" ||
		strings.Contains(brain.normalizeRu(text), "трансформация ") ||
		strings.Contains(brain.normalizeRu(text), "трансформацию ") ||
		strings.Contains(brain.normalizeRu(text), "трансформации ") {
		tokens := map[string]string{
			"трансформация": "оргия гомогеев",
			"трансформацию": "оргию гомогеев",
			"трансформации": "оргии гомогеев",
		}
		result := text
		for k, v := range tokens {
			result = strings.Replace(result, k, v, -1)
		}
		return true, result
	}
	return false, ""
}

func (brain Brain) normalizeRu(text string) string {
	result := text
	for k, v := range brain.memory.normalisationMap {
		result = strings.Replace(result, k, v, -1)
	}
	return result
}

func (brain Brain) normalizeEn(text string) string {
	result := text
	for k, v := range brain.memory.normalisationMap {
		result = strings.Replace(result, v, k, -1)
	}
	return result
}

func (brain Brain) khaleesify(text string) string {
	result := strings.ToLower(text)
	for _, k := range brain.memory.mockingMapKeys {
		result = strings.Replace(result, k, brain.memory.mockingMap[k], -1)
	}
	return result
}
