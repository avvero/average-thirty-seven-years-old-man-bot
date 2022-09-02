package brain

import (
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"testing"
)

func Test_returnsForToxicResponse(t *testing.T) {
	brain := NewBrain(NewMemory(), true, statistics.NewScriber(), &HuggingFaceToxicityDetectorNoop{toxic: true})
	respond, response := brain.Decision(0, "any")
	expected := "токсик ебаный"
	if !respond || response != expected {
		t.Error("Response: ", expected, " != ", response)
	}
}
