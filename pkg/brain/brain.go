package brain

import (
	"strconv"
	"strings"

	"github.com/avvero/the_gamers_guild_bot/internal/knowledge"
	"github.com/avvero/the_gamers_guild_bot/internal/utils"
)

type Brain struct {
	memory       *Memory
	randomFactor bool
}

func NewBrain(memory *Memory, randomFactor bool) *Brain {
	return &Brain{memory: memory, randomFactor: randomFactor}
}

func (brain *Brain) Decision(chatId int64, text string) (respond bool, response string) {
	if !utils.Contains([]string{"0", "-1001733786877", "245851441", "-578279468"}, strconv.FormatInt(chatId, 10)) {
		return true, "Mr Moony presents his compliments to Professor Snape, and begs him to keep his abnormally large nose out of other people’s business."
	}
	if brain.randomFactor {
		if utils.RandomUpTo(100) == 0 {
			phrase := brain.GetSenselessPhrases()[utils.RandomUpTo(len(brain.GetSenselessPhrases()))]
			return true, phrase
		}
		if len(text) > 5 && utils.RandomUpTo(50) == 0 {
			phrase := brain.huefy(text)
			return true, phrase
		}
		if len(text) > 14 && utils.RandomUpTo(100) == 0 {
			phrase := brain.khaleesify(text)
			return true, phrase
		}
	}
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
	if strings.Contains(brain.normalizeRu(text), "трансформация") ||
		strings.Contains(brain.normalizeRu(text), "трансформацию") ||
		strings.Contains(brain.normalizeRu(text), "трансформации") {
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
	if strings.Contains(text, "java") ||
		strings.Contains(text, "джаба") ||
		strings.Contains(text, "джава") {
		return true, "джава-хуява, а я работаю на го"
	}
	if strings.Contains(text, "блокир") {
		return true, "пусть себе анус заблокируют"
	}
	if strings.Contains(text, "опять") {
		return true, "не опять, а снова"
	}
	return false, ""
}

func (brain *Brain) normalizeRu(text string) string {
	result := text
	for k, v := range brain.GetNormalizationMap() {
		result = strings.Replace(result, k, v, -1)
	}
	return result
}

func (brain *Brain) normalizeEn(text string) string {
	result := text
	for k, v := range brain.GetNormalizationMap() {
		result = strings.Replace(result, v, k, -1)
	}
	return result
}

func (brain *Brain) khaleesify(text string) string {
	result := strings.ToLower(text)
	for _, k := range brain.GetMockingMapKeys() {
		result = strings.Replace(result, k, brain.GetMockingMap()[k], -1)
	}
	return result
}

func (brain *Brain) GetSenselessPhrases() []string {
	return brain.memory.senselessPhrases
}

func (brain *Brain) GetMockingMap() map[string]string {
	return brain.memory.mockingMap
}

func (brain *Brain) GetMockingMapKeys() []string {
	return brain.memory.mockingMapKeys
}

func (brain *Brain) GetNormalizationMap() map[string]string {
	return brain.memory.normalisationMap
}

func (brain *Brain) huefy(text string) string {
	length := len(text)
	result := make([]rune, length*2)
	resultPosition := length*2 - 1
	runes := []rune(text)
	vowelsNumber := 0
	wordLength := 0
	for i := len(runes) - 1; i >= 0; i-- {
		if utils.ContainsRune(knowledge.Delimiters, runes[i]) {
			vowelsNumber = 0
			wordLength = 0
			result[resultPosition] = runes[i]
			resultPosition--
			continue
		} else {
			wordLength++
		}
		// treat two vowels as one
		if utils.ContainsRune(knowledge.Vowels, runes[i]) && i > 0 && !utils.ContainsRune(knowledge.Vowels, runes[i-1]) {
			vowelsNumber++
		}
		// look forward and take word length
		if vowelsNumber == 2 {
			//wordLength--
			for f := i; f >= 0 && !utils.ContainsRune(knowledge.Delimiters, runes[f]); f-- {
				wordLength++
			}
		}
		if vowelsNumber == 2 && wordLength < 5 {
			// skip
			vowelsNumber = 0
			for i >= 0 && !utils.ContainsRune(knowledge.Delimiters, runes[i]) {
				result[resultPosition] = runes[i]
				i--
				resultPosition--
			}
		} else if vowelsNumber == 2 {
			softRune := knowledge.VowelsSoftenMap[runes[i]]
			if softRune != 0 {
				result[resultPosition] = softRune
			} else {
				result[resultPosition] = runes[i]
			}
			resultPosition--

			result[resultPosition] = 'у'
			resultPosition--

			result[resultPosition] = 'х'
			resultPosition--
			// skip
			vowelsNumber = 0
			for i > 0 && !utils.ContainsRune(knowledge.Delimiters, runes[i-1]) {
				i--
			}
		} else {
			result[resultPosition] = runes[i]
			resultPosition--
		}
	}
	//trim
	payloadPosition := 0
	for ; payloadPosition < len(result); payloadPosition++ {
		if result[payloadPosition] != 0 {
			break
		}
	}
	if payloadPosition > 0 {
		trimmedResult := make([]rune, len(result)-payloadPosition)

		for i := 0; i < len(trimmedResult); i++ {
			trimmedResult[i] = result[payloadPosition]
			payloadPosition++
		}
		return string(trimmedResult)
	}
	return string(result)
}
