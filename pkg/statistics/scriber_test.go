package statistics

import (
	"fmt"
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"github.com/avvero/the_gamers_guild_bot/internal/utils"
	"strconv"
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
			}, 0.2)
		}
		for len(scriber.packs) != 0 {
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
			}, 0.4)
		}
		for len(scriber.packs) != 0 {
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
	firstToxicityScore := scriber.data.ChatStatistics[1].UsersStatistics["first"].ToxicityScore
	firstToxicityExpected := 0.00
	if fmt.Sprintf("%.2f", firstToxicityScore) != fmt.Sprintf("%.2f", firstToxicityExpected) {
		t.Errorf("Expected: \"%.2f\" but got: \"%.2f\"", firstToxicityExpected, firstToxicityScore)
	}
	secondMessageCounter := scriber.data.ChatStatistics[1].UsersStatistics["second"].MessageCounter
	if secondMessageCounter != 200 {
		t.Errorf("Expected: \"%d\" but got: \"%d\"", 200, secondMessageCounter)
	}
	secondToxicityScore := scriber.data.ChatStatistics[1].UsersStatistics["second"].ToxicityScore
	secondToxicityScoreExpected := 0.00
	if fmt.Sprintf("%.2f", secondToxicityScore) != fmt.Sprintf("%.2f", secondToxicityScoreExpected) {
		t.Errorf("Expected: \"%.2f\" but got: \"%.2f\"", secondToxicityScoreExpected, secondToxicityScore)
	}
	dailyMessageCounter := scriber.data.ChatStatistics[1].DailyStatistics[date].MessageCounter
	if dailyMessageCounter != 300 {
		t.Errorf("Expected: \"%d\" but got: \"%d\"", 100, dailyMessageCounter)
	}
	dailyToxicityScore := scriber.data.ChatStatistics[1].DailyStatistics[date].ToxicityScore
	dailyToxicityScoreExpected := 0.00
	if fmt.Sprintf("%.2f", dailyToxicityScore) != fmt.Sprintf("%.2f", dailyToxicityScoreExpected) {
		t.Errorf("Expected: \"%.2f\" but got: \"%.2f\"", dailyToxicityScoreExpected, dailyToxicityScore)
	}
}

func Test_statisticsSerialization(t *testing.T) {
	scriber := NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	}, 0)
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	}, 0)
	for len(scriber.packs) != 0 {
		time.Sleep(10 * time.Millisecond) // TODO none reliable
	}
	jsonString := utils.PrintJson(scriber.GetStatistics(1))
	date := time.Now().Format("2006-01-02")
	dateTime := time.Now().Format("2006-01-02 15:04:05")
	expected := `{"userStatistics":{"first":{"messageCounter":1,"toxicityScore":0,"tension":0,"lastMessageDate":"` +
		date + `","lastMessageDateTime":"` + dateTime + `"},"second":{"messageCounter":1,"toxicityScore":0,"tension":0,"lastMessageDate":"` +
		date + `","lastMessageDateTime":"` + dateTime + `"}},"dailyStatistics":{"` + date + `":{"messageCounter":2,"toxicityScore":0}}}`
	if jsonString != expected {
		t.Errorf("Test failed.\nExpected: \"%s\" \nbut got : \"%s\"", expected, jsonString)
	}
}

func Test_increaseMessageCounter(t *testing.T) {
	scriber := NewScriber()
	message := &telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	}
	scriber.Keep(message, 0)
	time.Sleep(100 * time.Millisecond) // TODO none reliable
	counter := scriber.GetUserStatistics(message)
	if counter != 1 {
		t.Errorf("Test failed.\nExpected: \"%d\" \nbut got : \"%d\"", 1, counter)
	}
	// increase
	scriber.IncreaseUserMessageStatistics(1, "first", 10)
	counter = scriber.GetUserStatistics(message)
	if counter != 11 {
		t.Errorf("Test failed.\nExpected: \"%d\" \nbut got : \"%d\"", 11, counter)
	}
	// decrease
	scriber.IncreaseUserMessageStatistics(1, "first", -10)
	counter = scriber.GetUserStatistics(message)
	if counter != 1 {
		t.Errorf("Test failed.\nExpected: \"%d\" \nbut got : \"%d\"", 1, counter)
	}
}

func Test_statisticsPrettyPrint(t *testing.T) {
	scriber := NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	}, 0.1)
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	}, 0.2)
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two",
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	}, 0.3)
	for len(scriber.packs) != 0 {
		time.Sleep(10 * time.Millisecond) // TODO none reliable
	}
	text := scriber.GetStatisticsPrettyPrint(1)
	date := time.Now().Format("2006-01-02")
	expected := `Top 10 users:
 - second: 2 (t: 0.00)
 - first: 1 (t: 0.00)

Last 10 days:
 - ` + date + `: 3 (t: 0.00)

To get more information visit: ?id=1`
	if text != expected {
		t.Errorf("Expected: \"%s\" but got: \"%s\"", expected, text)
	}
}

func Test_statisticsPrettyPrintReturnTopUsers(t *testing.T) {
	scriber := NewScriber()
	for i := 0; i < 10; i++ {
		for j := 0; j < i; j++ {
			scriber.Keep(&telegram.WebhookRequestMessage{
				From: &telegram.WebhookRequestMessageSender{Username: "user" + strconv.Itoa(i)}, Text: "message" + strconv.Itoa(j),
				Chat: &telegram.WebhookRequestMessageChat{Id: 1},
			}, 0)
		}
	}
	time.Sleep(100 * time.Millisecond) // TODO none reliable
	for len(scriber.packs) != 0 {
		time.Sleep(10 * time.Millisecond) // TODO none reliable
	}
	text := scriber.GetStatisticsPrettyPrint(1)
	date := time.Now().Format("2006-01-02")
	expected := `Top 10 users:
 - user9: 9 (t: 0.00)
 - user8: 8 (t: 0.00)
 - user7: 7 (t: 0.00)
 - user6: 6 (t: 0.00)
 - user5: 5 (t: 0.00)
 - user4: 4 (t: 0.00)
 - user3: 3 (t: 0.00)
 - user2: 2 (t: 0.00)
 - user1: 1 (t: 0.00)

Last 10 days:
 - ` + date + `: 45 (t: 0.00)

To get more information visit: ?id=1`
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
		}, 0)
	}
	for i := 0; i < 20; i++ {
		scriber.Keep(&telegram.WebhookRequestMessage{
			From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two",
			Chat: &telegram.WebhookRequestMessageChat{Id: 1},
		}, 0)
	}
	for len(scriber.packs) != 0 {
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
		ipsum dolor sit amet. the of and to in is it et by from or but has that are a o so for on as an no not t's t s
		http https www com ru`

	scriber := NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: text,
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	}, 0)
	time.Sleep(100 * time.Millisecond) // TODO none reliable
	for len(scriber.packs) != 0 {
		time.Sleep(10 * time.Millisecond) // TODO none reliable
	}
	expectedWordStatistics := map[string]int{
		"lorem":  4,
		"ipsum":  4,
		"dolor":  4,
		"dolore": 2,
		"no":     0,
		"the":    0,
		"of":     0,
		"and":    0,
		"to":     0,
		"in":     0,
		"is":     0,
		"it":     0,
		"et":     0,
		"by":     0,
		"from":   0,
		"or":     0,
		"but":    0,
		"has":    0,
		"that":   0,
		"are":    0,
		"a":      0,
		"o":      0,
		"so":     0,
		"for":    0,
		"on":     0,
		"as":     0,
		"an":     0,
		"not":    0,
		"t":      0,
		"s":      0,
		"http":   0,
		"https":  0,
		"www":    0,
		"ru":     0,
		"com":    0,
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
		проблема проблемы в на с и не ну он так там то что чё а как за ни у я это но ты все по же из бы уже его мой про меня
        вот до был было еще ещё или только если й они где есть мне даже когда да нет их вы ага для тоже прям него чтобы тут
		а б в г д е ж з`

	scriber := NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: text,
		Chat: &telegram.WebhookRequestMessageChat{Id: 1},
	}, 0)
	time.Sleep(100 * time.Millisecond) // TODO none reliable
	for len(scriber.packs) != 0 {
		time.Sleep(100 * time.Millisecond) // TODO none reliable
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
		"не":         0,
		"ну":         0,
		"он":         0,
		"так":        0,
		"там":        0,
		"то":         0,
		"что":        0,
		"чё":         0,
		"а":          0,
		"как":        0,
		"за":         0,
		"ни":         0,
		"у":          0,
		"я":          0,
		"это":        0,
		"но":         0,
		"ты":         0,
		"все":        0,
		"по":         0,
		"же":         0,
		"из":         0,
		"бы":         0,
		"уже":        0,
		"его":        0,
		"мой":        0,
		"про":        0,
		"меня":       0,
		"вот":        0,
		"до":         0,
		"был":        0,
		"было":       0,
		"ещё":        0,
		"еще":        0,
		"или":        0,
		"только":     0,
		"если":       0,
		"й":          0,
		"они":        0,
		"где":        0,
		"есть":       0,
		"мне":        0,
		"даже":       0,
		"когда":      0,
		"б":          0,
		"г":          0,
		"д":          0,
		"да":         0,
		"нет":        0,
		"их":         0,
		"вы":         0,
		"ага":        0,
		"для":        0,
		"тоже":       0,
		"прям":       0,
		"него":       0,
		"чтобы":      0,
		"тут":        0,
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

func Test_statisticsRussianTextWithPunctuation(t *testing.T) {
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
	}, 0)
	time.Sleep(100 * time.Millisecond) // TODO none reliable
	for len(scriber.packs) != 0 {
		time.Sleep(10 * time.Millisecond) // TODO none reliable
	}
	expectedWordStatistics := map[string]int{
		"очевидна":    1,
		"развития":    3,
		"активизации": 1,
		"забывать":    1,
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
