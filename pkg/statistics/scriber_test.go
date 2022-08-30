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
			scriber.Keep(&telegram.WebhookRequestMessage{
				Chat: &telegram.WebhookRequestMessageChat{Id: 1},
				From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
			})
		}
		for len(scriber.messages) != 0 {
			time.Sleep(10 * time.Millisecond) // TODO none reliable
		}
		firstMessageCounter := scriber.data.ChatStatistics[1].UsersStatistics["first"].MessageCounter
		if firstMessageCounter != 100 {
			t.Errorf("Expected: \"%d\" but got: \"%d\"", 100, firstMessageCounter)
		}
	}
	// SECOND USERNAME
	{
		// Keep 200 messages
		for i := 0; i < 200; i++ {
			scriber.Keep(&telegram.WebhookRequestMessage{
				From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two",
				Chat: &telegram.WebhookRequestMessageChat{Id: 1},
			})
		}
		for len(scriber.messages) != 0 {
			time.Sleep(10 * time.Millisecond) // TODO none reliable
		}
		secondMessageCounter := scriber.data.ChatStatistics[1].UsersStatistics["second"].MessageCounter
		if secondMessageCounter != 200 {
			t.Errorf("Expected: \"%d\" but got: \"%d\"", 200, secondMessageCounter)
		}
	}
	// BOTH
	firstMessageCounter := scriber.data.ChatStatistics[1].UsersStatistics["first"].MessageCounter
	if firstMessageCounter != 100 {
		t.Errorf("Expected: \"%d\" but got: \"%d\"", 100, firstMessageCounter)
	}
	secondMessageCounter := scriber.data.ChatStatistics[1].UsersStatistics["second"].MessageCounter
	if secondMessageCounter != 200 {
		t.Errorf("Expected: \"%d\" but got: \"%d\"", 200, secondMessageCounter)
	}
}

func Test_statisticsSerialization(t *testing.T) {
	scriber := NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	})
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	})
	for len(scriber.messages) != 0 {
		time.Sleep(10 * time.Millisecond) // TODO none reliable
	}
	jsonString := utils.PrintJson(scriber.GetStatistics(1))
	expected := "{\"userStatistics\":{\"first\":{\"username\":\"first\",\"messageCounter\":1},\"second\":{\"username\":\"second\",\"messageCounter\":1}}}"
	if jsonString != expected {
		t.Errorf("Expected: \"%s\" but got: \"%s\"", expected, jsonString)
	}
}

func _Test_statisticsPrettyPrint(t *testing.T) {
	scriber := NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	})
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	})
	for len(scriber.messages) != 0 {
		time.Sleep(10 * time.Millisecond) // TODO none reliable
	}
	jsonString := utils.PrettyPrint(scriber.GetStatistics(1))
	expected := "{userStatistics:{first:{username:first,messageCounter:1},second:{username:second,messageCounter:1}}}"
	if jsonString != expected {
		t.Errorf("Expected: \"%s\" but got: \"%s\"", expected, jsonString)
	}
}
