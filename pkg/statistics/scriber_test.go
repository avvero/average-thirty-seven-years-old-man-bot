package statistics

import (
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"testing"
	"time"
)

func Test_messageCounter(t *testing.T) {
	scriber := NewScriber()
	// FIRST USERNAME
	{
		// Keep 100 messages
		for i := 0; i < 100; i++ {
			scriber.Keep(&telegram.WebhookRequestMessage{From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one"})
		}
		for len(scriber.messages) != 0 {
			time.Sleep(10 * time.Millisecond) // TODO none reliable
		}
		firstMessageCounter := scriber.statistics.UsersStatistics["first"].MessageCounter
		if firstMessageCounter != 100 {
			t.Errorf("Expected messages: \"%d\" but got: \"%d\"", 100, firstMessageCounter)
		}
	}
	// SECOND USERNAME
	{
		// Keep 200 messages
		for i := 0; i < 200; i++ {
			scriber.Keep(&telegram.WebhookRequestMessage{From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two"})
		}
		for len(scriber.messages) != 0 {
			time.Sleep(10 * time.Millisecond) // TODO none reliable
		}
		secondMessageCounter := scriber.statistics.UsersStatistics["second"].MessageCounter
		if secondMessageCounter != 200 {
			t.Errorf("Expected messages: \"%d\" but got: \"%d\"", 200, secondMessageCounter)
		}
	}
	// BOTH
	firstMessageCounter := scriber.statistics.UsersStatistics["first"].MessageCounter
	if firstMessageCounter != 100 {
		t.Errorf("Expected messages: \"%d\" but got: \"%d\"", 100, firstMessageCounter)
	}
	secondMessageCounter := scriber.statistics.UsersStatistics["second"].MessageCounter
	if secondMessageCounter != 200 {
		t.Errorf("Expected messages: \"%d\" but got: \"%d\"", 200, secondMessageCounter)
	}
}
