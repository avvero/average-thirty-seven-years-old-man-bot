package brain

import (
	"fmt"
	"github.com/avvero/the_gamers_guild_bot/internal/huggingface"
)

type HuggingFaceToxicityDetector struct {
	apiClient *huggingface.HuggingFaceApiClient
	threshold float64
}

func NewHuggingFaceToxicityDetector(apiClient *huggingface.HuggingFaceApiClient, threshold float64) Opinion {
	return &HuggingFaceToxicityDetector{apiClient: apiClient, threshold: threshold}
}

func (detector HuggingFaceToxicityDetector) Express(text string) (has bool, response string) {
	score, error := detector.apiClient.ToxicityScore(text)
	if error != nil {
		fmt.Printf("Response code: %s", error)
	}
	if score > detector.threshold {
		return true, ""
	} else {
		return false, ""
	}
}

type HuggingFaceToxicityDetectorNoop struct {
	toxic bool
}

func (detector HuggingFaceToxicityDetectorNoop) Express(text string) (has bool, response string) {
	return detector.toxic, ""
}
