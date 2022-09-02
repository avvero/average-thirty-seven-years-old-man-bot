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

func (detector *ToxicityDetector) ToxicityScore(text string) (float64, error) {
	return detector.apiClient.ToxicityScore(text)
}

type HuggingFaceToxicityDetectorNoop struct {
	toxic bool
}

func (detector *HuggingFaceToxicityDetectorNoop) ToxicityScore(text string) (float64, error) {
	return 0, nil
}
