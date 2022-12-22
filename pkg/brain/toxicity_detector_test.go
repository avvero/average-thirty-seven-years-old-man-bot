package brain

import (
	"fmt"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"testing"
)

func Test_returnsForToxicResponse(t *testing.T) {
	data := map[float64]string{
		0.92:  "осторожнее",
		0.93:  "осторожнее",
		0.98:  "на грани щас",
		0.981: "на грани щас",
		0.99:  "токсик",
		1.00:  "токсик",
	}
	for score, expected := range data {
		brain := NewBrain(true, statistics.NewScriber(), &ToxicityDetectorNoop{score: score}, nil)
		respond := false
		response := ""
		for i := 0; i < 20; i++ {
			thisRespond, thisResponse, _ := brain.Decision(0, "any")
			if thisRespond && thisResponse == expected {
				respond = thisRespond
				response = thisResponse
			}
		}
		if !respond || response != expected {
			t.Error("Response for score "+fmt.Sprintf("%v", score)+" : ", expected, " != ", response)
		}
	}
}
