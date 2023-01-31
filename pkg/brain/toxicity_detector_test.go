package brain

import (
	"fmt"
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"testing"
	"time"
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
			thisRespond, thisResponse, _ := brain.Decision(0, "user", "any")
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

func Test_returnsForFullToxic(t *testing.T) {
	//setup
	now := time.Now()
	dateTimeFormat := "2006-01-02 15:04:05"
	scriber := statistics.NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
		Chat: &telegram.WebhookRequestMessageChat{Id: 0},
	}, 0)
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "second"}, Text: "two",
		Chat: &telegram.WebhookRequestMessageChat{Id: 0},
	}, 0)
	time.Sleep(10 * time.Millisecond) // TODO none reliable
	scriber.SetUserLastMessageDateTime(0, "user", now.Add(-time.Minute*time.Duration(1)).Format(dateTimeFormat))
	scriber.SetUserLastMessageDateTime(0, "third", now.Add(-time.Minute*time.Duration(1)).Format(dateTimeFormat))
	scriber.SetUserLastMessageDateTime(0, "fourth", now.Add(-time.Minute*time.Duration(4)).Format(dateTimeFormat))
	scriber.SetUserLastMessageDateTime(0, "fifth", now.Add(-time.Minute*time.Duration(5)).Format(dateTimeFormat))
	scriber.SetUserLastMessageDateTime(0, "sixth", now.Add(-time.Minute*time.Duration(6)).Format(dateTimeFormat))
	scriber.SetUserLastMessageDateTime(0, "seventh", now.Add(-time.Minute*time.Duration(7)).Format(dateTimeFormat))
	scriber.SetUserLastMessageDate(0, "fourth", "")
	// when
	toxicityDetector := &ToxicityDetectorNoop{score: 1.00}
	brain := NewBrain(true, scriber, toxicityDetector, nil)
	expected := `Выражаю глубокую озабоченность касательно токсичного поведения user, такое поведение нанесло моральный ущерб некоторым гражданам. Им будет выплачена компенсация за моральный ущерб: 
 - first: +10
 - fourth: +10
 - second: +10
 - third: +10
`
	respond, response, _ := brain.Decision(0, "user", "any")
	if !respond || response != expected {
		t.Error("Response for score "+fmt.Sprintf("%v", 1.00)+" : ", expected, " != ", response)
	}
}
func Test_returnsForFullToxicNoReparation(t *testing.T) {
	//setup
	scriber := statistics.NewScriber()
	// when
	toxicityDetector := &ToxicityDetectorNoop{score: 1.00}
	brain := NewBrain(true, scriber, toxicityDetector, nil)
	expected := `Выражаю глубокую озабоченность касательно токсичного поведения user, такое поведение могло нанесло моральный ущерб некоторым гражданам, 
но к счастью все отделались легким негативом.`
	respond, response, _ := brain.Decision(0, "user", "any")
	if !respond || response != expected {
		t.Error("Response for score "+fmt.Sprintf("%v", 1.00)+" : ", expected, " != ", response)
	}
}
