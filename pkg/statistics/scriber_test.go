package statistics

import (
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"github.com/avvero/the_gamers_guild_bot/internal/utils"
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
			t.Errorf("Expected: \"%d\" but got: \"%d\"", 100, firstMessageCounter)
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
			t.Errorf("Expected: \"%d\" but got: \"%d\"", 200, secondMessageCounter)
		}
	}
	// BOTH
	firstMessageCounter := scriber.statistics.UsersStatistics["first"].MessageCounter
	if firstMessageCounter != 100 {
		t.Errorf("Expected: \"%d\" but got: \"%d\"", 100, firstMessageCounter)
	}
	secondMessageCounter := scriber.statistics.UsersStatistics["second"].MessageCounter
	if secondMessageCounter != 200 {
		t.Errorf("Expected: \"%d\" but got: \"%d\"", 200, secondMessageCounter)
	}
}

func Test_statisticsSerialization(t *testing.T) {
	scriber := NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one"})
	scriber.Keep(&telegram.WebhookRequestMessage{From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two"})
	jsonString := utils.PrintJson(scriber.GetStatistics())
	expected := "{\"userStatistics\":{\"first\":{\"username\":\"first\",\"messageCounter\":1},\"second\":{\"username\":\"second\",\"messageCounter\":1}}}"
	if jsonString != expected {
		t.Errorf("Expected: \"%s\" but got: \"%s\"", expected, jsonString)
	}
}

func Test_statisticsPrettyPrint(t *testing.T) {
	scriber := NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one"})
	scriber.Keep(&telegram.WebhookRequestMessage{From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two"})
	jsonString := utils.PrettyPrint(scriber.GetStatistics())
	expected := "{\"userStatistics\":{\"first\":{\"username\":\"first\",\"messageCounter\":1},\"second\":{\"username\":\"second\",\"messageCounter\":1}}}"
	if jsonString != expected {
		t.Errorf("Expected: \"%s\" but got: \"%s\"", expected, jsonString)
	}
}
