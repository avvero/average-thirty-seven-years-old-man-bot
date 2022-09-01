package brain

import (
	"github.com/avvero/the_gamers_guild_bot/internal/huggingface"
)

type ToxicityDetector struct {
	apiClient *huggingface.HuggingFaceApiClient
}

func NewToxicityDetector(apiClient *huggingface.HuggingFaceApiClient) *ToxicityDetector {
	return &ToxicityDetector{apiClient: apiClient}
}

func (this ToxicityDetector) Express(text string) (has bool, response string) {
	score, _ := this.apiClient.ToxicityScore(text)
	if score > 0.97 {
		return true, ""
	} else {
		return false, ""
	}
}
