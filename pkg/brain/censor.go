package brain

import (
	"strings"
)

type Censor struct {
}

func (this Censor) Check(chatId int64, text string) (forbidden bool, response string) {
	//text = normalize(text)
	// List of word-prefixes
	for _, word := range []string{"Росси", "Путин", "Украин", "Аллах", "Мухамед", "бог", "Иисус"} {
		if strings.Contains(strings.ToLower(text), strings.ToLower(word)) {
			return true, ""
		}
	}
	return false, ""
}

func normalize(text string) string {
	text = strings.ToLower(text)
	textRunes := []rune(text)
	result := make([]rune, len(textRunes)*2)
	resultPosition := len(result) - 1
	for i := len(textRunes) - 1; i >= 0; i-- {
		if textRunes[i] == textRunes[i-1] {
			continue
		}
	}
	resultPosition--
	return text
}
