package brain

import (
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"testing"
)

func Test_returnsForToxicResponse(t *testing.T) {
	data := map[float64]string{
		0.8:  "осторожнее",
		0.9:  "на грани щас",
		0.98: "токсик",
	}
	for score, expected := range data {
		brain := NewBrain(true, statistics.NewScriber(), &ToxicityDetectorNoop{score: score})
		respond, response := brain.Decision(0, "any")
		if !respond || response != expected {
			t.Error("Response: ", expected, " != ", response)
		}
	}
}
