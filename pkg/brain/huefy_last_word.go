package brain

import (
	"github.com/avvero/the_gamers_guild_bot/internal/knowledge"
	"github.com/avvero/the_gamers_guild_bot/internal/utils"
)

type HuefyLastWordIntention struct {
}

func (this HuefyLastWordIntention) Express(text string) (has bool, response string) {
	length := len(text)
	result := make([]rune, length*2)
	resultPosition := length*2 - 1
	runes := []rune(text)
	vowelsNumber := 0
	replaced := false
	for i := len(runes) - 1; i >= 0; i-- {
		if utils.ContainsRune(knowledge.Delimiters, runes[i]) {
			result[resultPosition] = runes[i]
			resultPosition--
			continue
		}
		// treat two vowels as one
		if utils.ContainsRune(knowledge.Vowels, runes[i]) && i > 0 && !utils.ContainsRune(knowledge.Vowels, runes[i-1]) {
			vowelsNumber++
		}
		// treat two vowels as one if text contains only two vowels and start from one of it
		if utils.ContainsRune(knowledge.Vowels, runes[i]) && i == 0 {
			vowelsNumber++
		}
		if !replaced && vowelsNumber == 2 {
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
			// skip runes that are replaced
			for i > 0 && !utils.ContainsRune(knowledge.Delimiters, runes[i-1]) {
				i--
			}
			replaced = true
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
		return true, string(trimmedResult)
	}
	return true, string(result)
}
