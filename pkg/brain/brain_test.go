package brain

import (
	"errors"
	"fmt"
	"github.com/avvero/the_gamers_guild_bot/internal/data"
	"github.com/avvero/the_gamers_guild_bot/internal/openai"
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func Test_responseOnlyToWhitelisted(t *testing.T) {
	brain := NewBrain(false, statistics.NewScriber(), &ToxicityDetectorNoop{}, nil)
	data := map[string]string{
		"-1001733786877": "gg",
		"245851441":      "gg",
		"-578279468":     "gg",
		"1":              "1: Mr Moony presents his compliments to Professor Snape, and begs him to keep his abnormally large nose out of other people’s business.",
	}
	for k, expected := range data {
		chatId, _ := strconv.ParseInt(k, 10, 64)
		respond, response, _ := brain.Decision(chatId, "user", "gg")
		if !respond || response != expected {
			t.Error("Response for ", k, ": \n", expected, "\n != \n", response)
		}
	}
}

func Test_responseOnCommandInfo(t *testing.T) {
	brain := NewBrain(false, statistics.NewScriber(), &ToxicityDetectorNoop{}, nil)
	respond, response, _ := brain.Decision(0, "user", "/info")
	expected := "I'm bot"
	if !respond || response != expected {
		t.Error("Response for : ", expected, " != ", response)
	}
}

func _Test_responseOnCommandStatisticsOnEmptyStatistics(t *testing.T) {
	brain := NewBrain(false, statistics.NewScriber(), &ToxicityDetectorNoop{}, nil)
	respond, response, _ := brain.Decision(0, "user", "статистика хуистика")
	if respond || response != "" {
		t.Error("Expected {false, nil} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_responseOnCommandStatistics(t *testing.T) {
	scriber := statistics.NewScriberWithData(&data.Data{ChatStatistics: make(map[int64]*data.ChatStatistics)}, "http://url")
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
		Chat: &telegram.WebhookRequestMessageChat{Id: 0},
	}, 0)
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "two",
		Chat: &telegram.WebhookRequestMessageChat{Id: 0},
	}, 0)
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "fifth"}, Text: "two",
		Chat: &telegram.WebhookRequestMessageChat{Id: 0},
	}, 0)
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "fourth"}, Text: "two",
		Chat: &telegram.WebhookRequestMessageChat{Id: 0},
	}, 0)
	time.Sleep(100 * time.Millisecond) // TODO none reliable
	scriber.SetUserTension(0, "third", 1)
	scriber.SetUserLastMessageDate(0, "fifth", "")
	scriber.SetUserLastMessageDate(0, "fourth", "2022-01-01")
	//
	brain := NewBrain(false, scriber, &ToxicityDetectorNoop{}, nil)
	respond, response, _ := brain.Decision(0, "user", "статистика хуистика")
	date := time.Now().Format("2006-01-02")
	expected := `Top 10 users:
 - first: 2 (t: 0.00)

Last 10 days:
 - ` + date + `: 4 (t: 0.00)

Top 10 infuriating persons:
 - third: tension = 1

Oblivion: 3
 - fifth: 1 (t: 0.00), last message date: unknown
 - fourth: 1 (t: 0.00), last message date: 2022-01-01
 - third: 0 (t: 0.00), last message date: unknown

To get more information visit: http://url?id=0`
	if !respond || response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_SkillCosts(t *testing.T) {
	// setup
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, `{"choices": [{"text": "Сообщение от AI"}]}`)
	}))
	defer ts.Close()
	apiClient := openai.NewApiClient(ts.URL, "key")
	scriber := statistics.NewScriberWithData(&data.Data{ChatStatistics: make(map[int64]*data.ChatStatistics)}, "http://url")
	brain := NewBrain(false, scriber, &ToxicityDetectorNoop{}, &apiClient)
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
		Chat: &telegram.WebhookRequestMessageChat{Id: 0},
	}, 0)
	time.Sleep(100 * time.Millisecond) // TODO none reliable
	// when
	respond, response, _ := brain.Decision(0, "first", "интелекто ебанина текст текст текст")
	// then
	expected := `Сообщение от AI`
	if response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
	// when
	for i := 0; i <= 10; i++ {
		scriber.Keep(&telegram.WebhookRequestMessage{
			From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
			Chat: &telegram.WebhookRequestMessageChat{Id: 0},
		}, 0)
	}
	time.Sleep(100 * time.Millisecond) // TODO none reliable
	respond, response, _ = brain.Decision(0, "first", "интелекто ебанина текст текст текст")
	// then
	expected = `Сообщение от AI`
	if response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
	// and
	counter := scriber.GetUserMessageCount(0, "first")
	if counter != 10 {
		t.Errorf("Test failed.\nExpected: \"%d\" \nbut got : \"%d\"", 10, counter)
	}
}

func Test_responseOnCommandToxicityWithoutPhrase(t *testing.T) {
	scriber := statistics.NewScriber()
	brain := NewBrain(false, scriber, &ToxicityDetectorNoop{score: 0.0}, nil)
	respond, response, _ := brain.Decision(0, "user", "токсик ревиленто")
	expected := ""
	if respond || response != expected {
		t.Error("Expected {false, \"\"} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_responseOnCommandToxicityFailed(t *testing.T) {
	scriber := statistics.NewScriber()
	brain := NewBrain(false, scriber, &ToxicityDetectorNoop{score: 0, err: errors.New("что-то не вышло")}, nil)
	respond, response, _ := brain.Decision(0, "user", "токсик ревиленто что-то на токсичном")
	expected := "определить уровень токсичности не удалось, быть может вы - черт, попробуйте позже"
	if !respond || response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_responseOnCommandToxicity(t *testing.T) {
	scriber := statistics.NewScriber()
	brain := NewBrain(false, scriber, &ToxicityDetectorNoop{score: 0.98}, nil)
	respond, response, _ := brain.Decision(0, "user", "токсик ревиленто что-то на токсичном")
	expected := "уровень токсичности 98%"
	if !respond || response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_returnsOnSomeText(t *testing.T) {
	brain := NewBrain(false, statistics.NewScriber(), &ToxicityDetectorNoop{}, nil)
	data := map[string]string{
		"gg": "gg",
		"GG": "gg",
		"gG": "gg",
		"Gg": "gg",

		"morrowind":             "Morrowind - одна из лучших игр эва",
		"моровинд":              "Morrowind - одна из лучших игр эва",
		"морровинд":             "Morrowind - одна из лучших игр эва",
		"бла бла бла morrowind": "Morrowind - одна из лучших игр эва",
		"бла бла бла morrowind бла бла бла": "Morrowind - одна из лучших игр эва",

		"Elden Ring": "Elden Ring - это величие",
		"бла бла бла Elden Ring бла бла бла": "Elden Ring - это величие",
		"elden ring": "Elden Ring - это величие",
		"бла бла бла elden ring бла бла бла": "Elden Ring - это величие",
		"ER": "Elden Ring - это величие",
		"бла бла бла ER бла бла бла": "Elden Ring - это величие",
		"ЕР": "Elden Ring - это величие",
		"бла бла бла ЕР бла бла бла": "Elden Ring - это величие",
		"ЭР": "Elden Ring - это величие",
		"бла бла бла ЭР бла бла бла": "Elden Ring - это величие",
		"Элден ринг":                 "Elden Ring - это величие",
		"Елден РИНГ":                 "Elden Ring - это величие",

		"spotify":  "Эти н'вахи Антону косарик должны за подписку",
		"Spotify":  "Эти н'вахи Антону косарик должны за подписку",
		"спотифай": "Эти н'вахи Антону косарик должны за подписку",
		"бла бла бла spotify бла бла бла":  "Эти н'вахи Антону косарик должны за подписку",
		"бла бла бла спотифай бла бла бла": "Эти н'вахи Антону косарик должны за подписку",

		"scrum": "Скрам - это пережиток",
		"скрам": "Скрам - это пережиток",

		//"java":          "джава-хуява, а я работаю на го",
		//"бла java бла":  "джава-хуява, а я работаю на го",
		//"джава":         "джава-хуява, а я работаю на го",
		//"бла джава бла": "джава-хуява, а я работаю на го",
		//"джаба":         "джава-хуява, а я работаю на го",
		//"бла джаба бла": "джава-хуява, а я работаю на го",

		//"да там же уже решено, теперь надо новое что-то заблокировать": "пусть себе анус заблокируют",
		//"они это вчера заблокировали и ждут":                           "пусть себе анус заблокируют",
		//"это меня блокирует сильно":                                    "пусть себе анус заблокируют",

		"что-то про mass effect и так далее": "Шепард умрет",
		"что-то про масс эффект и так далее": "Шепард умрет",
	}
	for origin, expected := range data {
		respond, response, _ := brain.Decision(0, "user", origin)
		if !respond || response != expected {
			t.Error("Response for", origin, ":", expected, "!=", response)
		}
	}
}

func Test_returnsOnSomeTextWithRandomFactor(t *testing.T) {
	brain := NewBrain(true, statistics.NewScriber(), &ToxicityDetectorNoop{score: 0}, nil)
	data := map[string]string{
		"я думал сначала Медведев это опять. А тут какой то давыдов": "не опять, а снова",

		//"двуличие": "хуичие",

		//"Потому что надо было не шифры ваши писать": "хуёму что хуядо хуило не хуифры хуяши хуисать",

		//"купил":             "А не пиздишь? Аренда это не покупка",
		//"бла бла бла купил": "А не пиздишь? Аренда это не покупка",
		//"бла бла бла купил бла бла бла": "А не пиздишь? Аренда это не покупка",

		"трансформация":                "пертурбация",
		"трансформацию":                "пертурбацию",
		"трансформации":                "пертурбации",
		"сейчас трансформация пройдет": "сейчас пертурбация пройдет",
		"сейчас трансформацию пройдет": "сейчас пертурбацию пройдет",
		"сейчас трансформации пройдет": "сейчас пертурбации пройдет",
		"с этими трансформациями уже":  "с этими пертурбациями уже",

		"у нас проблема":           "у меня есть 5-10 солюшенов этой проблемы",
		"у нас проблема, товарищи": "у меня есть 5-10 солюшенов этой проблемы",
		"проблема в том":           "у меня есть 5-10 солюшенов этой проблемы",
		"куча проблем":             "у меня есть 5-10 солюшенов этой проблемы",

		"devops": "Девопсы не нужны",
		"девопс": "Девопсы не нужны",
	}
	for origin, expected := range data {
		respond := false
		response := ""
		for i := 0; i < 500; i++ {
			thisRespond, thisResponse, _ := brain.Decision(0, "user", origin)
			if thisRespond && thisResponse == expected {
				respond = thisRespond
				response = thisResponse
			}
		}
		if !respond || response != expected {
			t.Error("Response for ", origin, ": ", expected, " != ", response)
		}
	}
}

func Test_returnsOnNotElderRing(t *testing.T) {
	brain := NewBrain(true, statistics.NewScriber(), &ToxicityDetectorNoop{}, nil)
	data := []string{
		"pERt",
		"sdfERdfd",
		"аааЕРваа",
		"трансформационный1",
	}
	for _, text := range data {
		respond, response, _ := brain.Decision(0, "user", text)
		if respond {
			t.Errorf("Not expected: \"%s\"", response)
		}
	}
}

func _Test_returnsForLuckyKhaleesifiedText(t *testing.T) {
	brain := NewBrain(true, statistics.NewScriber(), &ToxicityDetectorNoop{}, nil)
	respond := false
	response := ""
	expected := "делись зя миня, дляконь"
	for i := 0; i < 500; i++ {
		thisRespond, thisResponse, _ := brain.Decision(0, "user", "дерись за меня, дракон")
		if thisRespond && thisResponse == expected {
			respond = thisRespond
			response = thisResponse
		}
	}
	if !respond {
		t.Error("Expected and got:", expected, " != ", response)
	}
}

func Test_NotificationUnrecognized(t *testing.T) {
	// setup
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, `{"choices": [{"text": "{\"action\": \"жопа жопа тра та та\", \"time\": \"\"}"}]}`)
	}))
	defer ts.Close()
	apiClient := openai.NewApiClient(ts.URL, "key")
	scriber := statistics.NewScriberWithData(&data.Data{ChatStatistics: make(map[int64]*data.ChatStatistics)}, "http://url")
	brain := NewBrain(false, scriber, &ToxicityDetectorNoop{}, &apiClient)
	// when
	respond, response, _ := brain.Decision(0, "first", "мементо foo bar")
	// then
	expected := `Не могу разобрать, что ты написал, напиши нормально`
	if response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_NotificationDeformed(t *testing.T) {
	// setup
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, `{"choices": [{"text": "неожиданный текст"}]}`)
	}))
	defer ts.Close()
	apiClient := openai.NewApiClient(ts.URL, "key")
	scriber := statistics.NewScriberWithData(&data.Data{ChatStatistics: make(map[int64]*data.ChatStatistics)}, "http://url")
	brain := NewBrain(false, scriber, &ToxicityDetectorNoop{}, &apiClient)
	// when
	respond, response, _ := brain.Decision(0, "first", "мементо foo bar")
	// then
	expected := `Не могу разобрать, что ты написал, напиши нормально`
	if response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_NotificationOverdue(t *testing.T) {
	// setup
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, `{"choices": [{"text": "{\"action\": \"мементо ответить Сереге через 1 час и 15 минут\", \"time\": \"2023-01-24 21:07\"}"}]}`)
	}))
	defer ts.Close()
	apiClient := openai.NewApiClient(ts.URL, "key")
	scriber := statistics.NewScriberWithData(&data.Data{ChatStatistics: make(map[int64]*data.ChatStatistics)}, "http://url")
	brain := NewBrain(false, scriber, &ToxicityDetectorNoop{}, &apiClient)
	// when
	respond, response, _ := brain.Decision(0, "first", "мементо ответить Сереге через 1 час и 15 минут назад")
	// then
	expected := `И как ты себе это представляешь, пес?`
	if response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_NotificationSucceeded(t *testing.T) {
	// setup
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, `{"choices": [{"text": "{\"action\": \"ответить Сереге через 1 час и 15 минут\", \"time\": \"2030-01-24 23:37\"}"}]}`)
	}))
	defer ts.Close()
	apiClient := openai.NewApiClient(ts.URL, "key")
	scriber := statistics.NewScriberWithData(&data.Data{ChatStatistics: make(map[int64]*data.ChatStatistics)}, "http://url")
	brain := NewBrain(true, scriber, &ToxicityDetectorNoop{}, &apiClient)
	// when
	respond, response, _ := brain.Decision(0, "first", "мементо ответить Сереге через 1 час и 15 минут")
	// then
	expected := `Напомню «ответить Сереге через 1 час и 15 минут» в 2030-01-24 23:37`
	if response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
		return
	}
	expectedAction := "ответить Сереге через 1 час и 15 минут"
	gotAction := ""
	if scriber.GetNotifications(0) != nil && scriber.GetNotifications(0)["2030-01-24 23:37"] != nil {
		gotAction = scriber.GetNotifications(0)["2030-01-24 23:37"].Action
	}
	if gotAction != expectedAction {
		t.Error("Expected: " + expectedAction + ", but got:" + gotAction)
	}
	//when
	respond, response, _ = brain.Decision(0, "first", "мементос ревиленто")
	// then
	expected = `Notifications:
 - 2030-01-24 23:37: ответить Сереге через 1 час и 15 минут
`
	if response != expected {
		t.Errorf("Expected: \"%s\" but got: \"%s\"", expected, response)
		return
	}
	//when
	scriber.RemoveNotification(0, "2030-01-24 23:37")
	respond, response, _ = brain.Decision(0, "first", "мементос ревиленто")
	// then
	expected = `Ничего нет, просто ничего, просто 0, ни ху я`
	if response != expected {
		t.Errorf("Expected: \"%s\" but got: \"%s\"", expected, response)
		return
	}
}

func Test_NotificationSucceeded2(t *testing.T) {
	// setup
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, `{"choices": [{"text": "{'action': 'ответить Сереге через 1 час и 15 минут', 'time': '2030-01-24 23:37'}"}]}`)
	}))
	defer ts.Close()
	apiClient := openai.NewApiClient(ts.URL, "key")
	scriber := statistics.NewScriberWithData(&data.Data{ChatStatistics: make(map[int64]*data.ChatStatistics)}, "http://url")
	brain := NewBrain(true, scriber, &ToxicityDetectorNoop{}, &apiClient)
	// when
	respond, response, _ := brain.Decision(0, "first", "мементо ответить Сереге через 1 час и 15 минут")
	// then
	expected := `Напомню «ответить Сереге через 1 час и 15 минут» в 2030-01-24 23:37`
	if response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
		return
	}
}
