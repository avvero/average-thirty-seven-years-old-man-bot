package brain

import (
	"fmt"
	"github.com/avvero/the_gamers_guild_bot/internal/openai"
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func Test_returnsForToxicResponse(t *testing.T) {
	data := map[float64]string{
		//0.92:  "осторожнее",
		//0.93:  "осторожнее",
		//0.98:  "на грани щас",
		//0.981: "на грани щас",
		//0.99:  "токсик",
		//1.00:  "токсик",
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
	// setup
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, `{"choices": [{"text": "Выражаю глубокую озабоченность касательно токсичного поведения user, такое поведение нанесло моральный ущерб некоторым гражданам. Им будет выплачена компенсация:"}]}`)
	}))
	defer ts.Close()
	apiClient := openai.NewApiClient(ts.URL, "key")
	now := time.Now()
	dateTimeFormat := "2006-01-02 15:04:05"
	scriber := statistics.NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
		Chat: &telegram.WebhookRequestMessageChat{Id: 0},
	}, 0)
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
	brain := NewBrain(true, scriber, toxicityDetector, &apiClient)
	expected := `Выражаю глубокую озабоченность касательно токсичного поведения user, такое поведение нанесло моральный ущерб некоторым гражданам. Им будет выплачена компенсация:
 - first: +10
 - fourth: +10
 - second: +10
 - third: +10
`
	respond, response, _ := brain.Decision(0, "user", "any")
	if !respond || response != expected {
		t.Error("Response for score "+fmt.Sprintf("%v", 1.00)+" : ", expected, " != ", response)
	}

	//
	respond, response, _ = brain.Decision(0, "user", "статистика хуистика")
	date := time.Now().Format("2006-01-02")
	expected = `Top 10 users:
 - first: 12 (t: 0.00)
 - second: 11 (t: 0.00)

Last 10 days:
 - ` + date + `: 3 (t: 0.00)

Houses scores:
 - Гриффиндор: -10

Between-Morrowind-And-Skyrim: 6

To get more information visit: ?id=0`
	if !respond || response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}
func Test_returnsForFullToxicNoReparation(t *testing.T) {
	//setup
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, `{"choices": [{"text": "Выражаю глубокую озабоченность касательно токсичного поведения user, такое поведение могло нанесло моральный ущерб некоторым гражданам, 
но к счастью все отделались легким негативом."}]}`)
	}))
	defer ts.Close()
	apiClient := openai.NewApiClient(ts.URL, "key")
	scriber := statistics.NewScriber()
	// when
	toxicityDetector := &ToxicityDetectorNoop{score: 1.00}
	brain := NewBrain(true, scriber, toxicityDetector, &apiClient)
	expected := `Выражаю глубокую озабоченность касательно токсичного поведения user, такое поведение могло нанесло моральный ущерб некоторым гражданам, 
но к счастью все отделались легким негативом.`
	respond, response, _ := brain.Decision(0, "user", "any")
	if !respond || response != expected {
		t.Error("Response for score "+fmt.Sprintf("%v", 1.00)+" : ", expected, " != ", response)
	}
}
