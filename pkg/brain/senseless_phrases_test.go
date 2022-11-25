package brain

import (
	"github.com/avvero/the_gamers_guild_bot/internal/utils"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"testing"
)

func Test_returnsForLuckySenselessPhrase(t *testing.T) {
	brain := NewBrain(true, statistics.NewScriber(), &ToxicityDetectorNoop{})
	respond := false
	for i := 0; i < 1000; i++ {
		thisRespond, thisResponse, _ := brain.Decision(0, "any")
		if thisRespond && utils.Contains(senselessPhrases, thisResponse) {
			respond = thisRespond
			break
		}
	}
	if !respond {
		t.Error("Expected something from senselessPhrases but got nothing")
	}
}
