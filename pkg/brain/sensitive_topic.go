package brain

import (
	"strings"
)

type SensitiveTopic struct {
}

func (this SensitiveTopic) Check(chatId int64, text string) (forbidden bool, response string) {
	text = strings.ToLower(text)
	for _, word := range []string{"росси", "путин", "украин", "аллах", "мухаммед", "бог", "иисус"} {
		if strings.Contains(text, word) {
			return true, ""
		}
	}
	return false, ""
}
