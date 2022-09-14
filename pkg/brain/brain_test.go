package brain

import (
	"errors"
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"strconv"
	"testing"
	"time"
)

func Test_responseOnlyToWhitelisted(t *testing.T) {
	brain := NewBrain(false, statistics.NewScriber(), &ToxicityDetectorNoop{})
	data := map[string]string{
		"-1001733786877": "gg",
		"245851441":      "gg",
		"-578279468":     "gg",
		"1":              "Mr Moony presents his compliments to Professor Snape, and begs him to keep his abnormally large nose out of other people’s business.",
	}
	for k, expected := range data {
		chatId, _ := strconv.ParseInt(k, 10, 64)
		respond, response := brain.Decision(chatId, "gg")
		if !respond || response != expected {
			t.Error("Response for ", k, ": ", expected, " != ", response)
		}
	}
}

func Test_responseOnCommandInfo(t *testing.T) {
	brain := NewBrain(false, statistics.NewScriber(), &ToxicityDetectorNoop{})
	respond, response := brain.Decision(0, "/info")
	expected := "I'm bot"
	if !respond || response != expected {
		t.Error("Response for : ", expected, " != ", response)
	}
}

func Test_responseOnCommandStatisticsOnEmptyStatistics(t *testing.T) {
	brain := NewBrain(false, statistics.NewScriber(), &ToxicityDetectorNoop{})
	respond, response := brain.Decision(0, "/statistics")
	if respond || response != "" {
		t.Error("Expected {false, nil} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_responseOnCommandStatistics(t *testing.T) {
	scriber := statistics.NewScriber()
	scriber.Keep(&telegram.WebhookRequestMessage{
		From: &telegram.WebhookRequestMessageSender{Username: "first"}, Text: "one",
		Chat: &telegram.WebhookRequestMessageChat{Id: 0},
	})
	time.Sleep(100 * time.Millisecond) // TODO none reliable
	brain := NewBrain(false, scriber, &ToxicityDetectorNoop{})
	respond, response := brain.Decision(0, "/statistics")
	date := time.Now().Format("2006-01-02")
	expected := "Statistics by user:\n - first: 1\nStatistics by day:\n - " + date + ": 1\n"
	if !respond || response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_responseOnCommandToxicityWithoutPhrase(t *testing.T) {
	scriber := statistics.NewScriber()
	brain := NewBrain(false, scriber, &ToxicityDetectorNoop{score: 0.0})
	respond, response := brain.Decision(0, "/toxicity")
	expected := ""
	if respond || response != expected {
		t.Error("Expected {false, \"\"} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_responseOnCommandToxicityFailed(t *testing.T) {
	scriber := statistics.NewScriber()
	brain := NewBrain(false, scriber, &ToxicityDetectorNoop{score: 0, err: errors.New("что-то не вышло")})
	respond, response := brain.Decision(0, "/toxicity что-то на токсичном")
	expected := "определить уровень токсичности не удалось, быть может вы - черт, попробуйте позже"
	if !respond || response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_responseOnCommandToxicity(t *testing.T) {
	scriber := statistics.NewScriber()
	brain := NewBrain(false, scriber, &ToxicityDetectorNoop{score: 0.98})
	respond, response := brain.Decision(0, "/toxicity что-то на токсичном")
	expected := "уровень токсичности 98%"
	if !respond || response != expected {
		t.Error("Expected {true, " + expected + "} but got {" + strconv.FormatBool(respond) + ", " + response + "}")
	}
}

func Test_returnsOnSomeText(t *testing.T) {
	brain := NewBrain(false, statistics.NewScriber(), &ToxicityDetectorNoop{})
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

		"devops": "Девопсы не нужны",
		"девопс": "Девопсы не нужны",

		"scrum": "Скрам - это пережиток",
		"скрам": "Скрам - это пережиток",

		"java":          "джава-хуява, а я работаю на го",
		"бла java бла":  "джава-хуява, а я работаю на го",
		"джава":         "джава-хуява, а я работаю на го",
		"бла джава бла": "джава-хуява, а я работаю на го",
		"джаба":         "джава-хуява, а я работаю на го",
		"бла джаба бла": "джава-хуява, а я работаю на го",

		"да там же уже решено, теперь надо новое что-то заблокировать": "пусть себе анус заблокируют",
		"они это вчера заблокировали и ждут":                           "пусть себе анус заблокируют",
		"это меня блокирует сильно":                                    "пусть себе анус заблокируют",

		"что-то про mass effect и так далее": "Шепард умрет",
		"что-то про масс эффект и так далее": "Шепард умрет",
	}
	for origin, expected := range data {
		respond, response := brain.Decision(0, origin)
		if !respond || response != expected {
			t.Error("Response for", origin, ":", expected, "!=", response)
		}
	}
}

func Test_returnsOnSomeTextWithRandomFactor(t *testing.T) {
	brain := NewBrain(true, statistics.NewScriber(), &ToxicityDetectorNoop{score: 0})
	data := map[string]string{
		"я думал сначала Медведев это опять. А тут какой то давыдов": "не опять, а снова",

		"двуличие": "хуичие",

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

		"Я не понял - вот делают трансформацию, ч": "я не понял - вот делают пертурбацию, ч",

		"у нас проблема":           "у меня есть 5-10 солюшенов этой проблемы",
		"у нас проблема, товарищи": "у меня есть 5-10 солюшенов этой проблемы",
		"проблема в том":           "у меня есть 5-10 солюшенов этой проблемы",
		"куча проблем":             "у меня есть 5-10 солюшенов этой проблемы",
	}
	for origin, expected := range data {
		respond := false
		response := ""
		for i := 0; i < 500; i++ {
			thisRespond, thisResponse := brain.Decision(0, origin)
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
	brain := NewBrain(true, statistics.NewScriber(), &ToxicityDetectorNoop{})
	data := []string{
		"pERt",
		"sdfERdfd",
		"аааЕРваа",
		"трансформационный1",
	}
	for _, text := range data {
		respond, response := brain.Decision(0, text)
		if respond {
			t.Errorf("Not expected: \"%s\"", response)
		}
	}
}

func _Test_returnsForLuckyKhaleesifiedText(t *testing.T) {
	brain := NewBrain(true, statistics.NewScriber(), &ToxicityDetectorNoop{})
	respond := false
	response := ""
	expected := "делись зя миня, дляконь"
	for i := 0; i < 500; i++ {
		thisRespond, thisResponse := brain.Decision(0, "дерись за меня, дракон")
		if thisRespond && thisResponse == expected {
			respond = thisRespond
			response = thisResponse
		}
	}
	if !respond {
		t.Error("Expected and got:", expected, " != ", response)
	}
}

func Test_censorTests(t *testing.T) {
	brain := NewBrain(true, statistics.NewScriber(), &ToxicityDetectorNoop{})
	data := []string{
		"Россия",
		"Росия",
		"Россию",
		"Росию",
		"Путин",
		"Украина",
		"Украине",
		"Аллах",
		"Алах",
		"Мухаммед",
		"Мухамед",
		"бог",
		"Иисус",
		"Исус",
	}
	for _, text := range data {
		respond, response := brain.Decision(0, text)
		if respond {
			t.Error("Response for ", text, ": not expected != ", response)
		}
	}
}
