package brain

import (
	"github.com/avvero/the_gamers_guild_bot/internal/utils"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"testing"
)

func Test_returnsForLuckySenselessPhrase(t *testing.T) {
	brain := NewBrain(NewMemory(), true, statistics.NewScriber())
	respond := false
	response := ""
	for i := 0; i < 500; i++ {
		thisRespond, thisResponse := brain.Decision(0, "any")
		if thisRespond {
			respond = thisRespond
			response = thisResponse
		}
	}
	if !respond || !utils.Contains(senselessPhrases, response) {
		t.Error("Expected and got: something from senselessPhrases != ", response)
	}
}
