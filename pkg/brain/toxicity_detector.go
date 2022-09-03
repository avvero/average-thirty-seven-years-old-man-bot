package brain

import (
	"github.com/avvero/the_gamers_guild_bot/internal/huggingface"
)

type ToxicityDetector interface {
	ToxicityScore(text string) (float64, error)
}
type HuggingFaceToxicityDetector struct {
	apiClient huggingface.HuggingFaceApiClient
}

func NewToxicityDetector(apiClient huggingface.HuggingFaceApiClient) *HuggingFaceToxicityDetector {
	return &HuggingFaceToxicityDetector{apiClient: apiClient}
}

func (detector HuggingFaceToxicityDetector) ToxicityScore(text string) (float64, error) {
	return detector.apiClient.ToxicityScore(text)
}

type ToxicityDetectorNoop struct {
	result float64
}

func (detector ToxicityDetectorNoop) ToxicityScore(text string) (float64, error) {
	return detector.result, nil
}
