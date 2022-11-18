package statistics

import (
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"github.com/avvero/the_gamers_guild_bot/internal/utils"
	"testing"
	"time"
)

func Test_messageCounter(t *testing.T) {
	scriber := NewScriber()
	date := time.Now().Format("2006-01-02")
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
		dailyMessageCounter := scriber.data.ChatStatistics[1].DailyStatistics[date].MessageCounter
		if dailyMessageCounter != 100 {
			t.Errorf("Expected: \"%d\" but got: \"%d\"", 100, dailyMessageCounter)
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
	dailyMessageCounter := scriber.data.ChatStatistics[1].DailyStatistics[date].MessageCounter
	if dailyMessageCounter != 300 {
		t.Errorf("Expected: \"%d\" but got: \"%d\"", 100, dailyMessageCounter)
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
	date := time.Now().Format("2006-01-02")
	expected := `{"userStatistics":{"first":{"messageCounter":1},"second":{"messageCounter":1}},"dailyStatistics":{"` + date + `":{"messageCounter":2}},"dailyWordStatistics":{"` + date + `":{"one":1,"two":1}}}`
	if jsonString != expected {
		t.Errorf("Expected: \"%s\" but got: \"%s\"", expected, jsonString)
	}
}

func Test_statisticsPrettyPrint(t *testing.T) {
	scriber := NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	})
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	})
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	})
	for len(scriber.messages) != 0 {
		time.Sleep(10 * time.Millisecond) // TODO none reliable
	}
	text := scriber.GetStatisticsPrettyPrint(1)
	date := time.Now().Format("2006-01-02")
	expected := "Statistics by user:\n - second: 2\n - first: 1\nStatistics by day:\n - " + date + ": 3\n"
	if text != expected {
		t.Errorf("Expected: \"%s\" but got: \"%s\"", expected, text)
	}
}

func Test_statisticsForSingleWords(t *testing.T) {
	scriber := NewScriber()
	for i := 0; i < 10; i++ {
		scriber.Keep(&telegram.WebhookRequestMessage{
			From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
			Chat: &telegram.WebhookRequestMessageChat{Id: 1},
		})
	}
	for i := 0; i < 20; i++ {
		scriber.Keep(&telegram.WebhookRequestMessage{
			From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two",
			Chat: &telegram.WebhookRequestMessageChat{Id: 1},
		})
	}
	for len(scriber.messages) != 0 {
		time.Sleep(10 * time.Millisecond) // TODO none reliable
	}
	date := time.Now().Format("2006-01-02")
	oneWordCounter := scriber.GetStatistics(1).DailyWordStatistics[date]["one"]
	twoWordCounter := scriber.GetStatistics(1).DailyWordStatistics[date]["two"]
	if oneWordCounter != 10 {
		t.Errorf("Expected: \"%d\" but got: \"%d\"", 10, oneWordCounter)
	}
	if twoWordCounter != 20 {
		t.Errorf("Expected: \"%d\" but got: \"%d\"", 20, twoWordCounter)
	}
	println("Result:", utils.PrintJson(scriber.GetStatistics(1)))
}

func Test_statisticsText(t *testing.T) {
	text := `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore 
		et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita 
		kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur 
		sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. 
		At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem 
		ipsum dolor sit amet.`

	scriber := NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: text,
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	})
	for len(scriber.messages) != 0 {
		time.Sleep(10 * time.Millisecond) // TODO none reliable
	}
	expectedWordStatistics := map[string]int{
		"lorem":  4,
		"ipsum":  4,
		"dolor":  4,
		"dolore": 2,
		"no":     2,
	}
	date := time.Now().Format("2006-01-02")
	for word, number := range expectedWordStatistics {
		counter := scriber.GetStatistics(1).DailyWordStatistics[date][word]
		if counter != number {
			t.Errorf("Expected for %v: \"%d\" but got: \"%d\"", word, number, counter)
		}
	}
	println("Result:", utils.PrintJson(scriber.GetStatistics(1)))
}

func Test_statisticsRussianText(t *testing.T) {
	text := `Значимость этих проблем настолько очевидна, что реализация намеченных плановых заданий способствует 
		подготовки и реализации дальнейших направлений развития. Равным образом постоянный количественный рост и сфера нашей 
		активности позволяет оценить значение соответствующий условий активизации. Разнообразный и богатый опыт реализация 
		намеченных плановых заданий требуют определения и уточнения направлений прогрессивного развития. Не следует, однако 
		забывать, что реализация намеченных плановых заданий представляет собой интересный эксперимент проверки дальнейших 
		направлений развития.
		проблема проблемы в на с и`

	scriber := NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: text,
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	})
	for len(scriber.messages) != 0 {
		time.Sleep(10 * time.Millisecond) // TODO none reliable
	}
	expectedWordStatistics := map[string]int{
		"значимость": 1,
		"проблем":    1,
		"проблема":   1,
		"проблемы":   1,
		"в":          0,
		"на":         0,
		"с":          0,
		"и":          0,
	}
	date := time.Now().Format("2006-01-02")
	for word, number := range expectedWordStatistics {
		counter := scriber.GetStatistics(1).DailyWordStatistics[date][word]
		if counter != number {
			t.Errorf("Expected for %v: \"%d\" but got: \"%d\"", word, number, counter)
		}
	}
	println("Result:", utils.PrintJson(scriber.GetStatistics(1)))
}
