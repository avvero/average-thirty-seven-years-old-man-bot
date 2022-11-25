package brain

import (
	"github.com/avvero/the_gamers_guild_bot/internal/utils"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"testing"
)

func Test_returnsForNoThenPhrase(t *testing.T) {
	brain := NewBrain(true, statistics.NewScriber(), &ToxicityDetectorNoop{})
	thisRespond, thisResponse, _ := brain.Decision(0, "нет")
	if !thisRespond || !utils.Contains(noThenResponses, thisResponse) {
		t.Error("Expected something from noThenResponses but got:", thisResponse)
	}
}
