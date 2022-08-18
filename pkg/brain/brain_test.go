package brain

import (
	"strconv"
	"testing"
)

func Test_responseOnlyToWhitelisted(t *testing.T) {
	brain := NewBrain(NewMemory(), false)
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

func Test_returnsOnSomeText(t *testing.T) {
	brain := NewBrain(NewMemory(), false)
	data := map[string]string{
		"gg": "gg",
		"GG": "gg",
		"gG": "gg",
		"Gg": "gg",

		"нет": "пидора ответ",
		"НЕТ": "пидора ответ",
		"Нет": "пидора ответ",
		"НеТ": "пидора ответ",
		"Heт": "пидора ответ",
		"Нeт": "пидора ответ",

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

		"купил":             "А не пиздишь? Аренда это не покупка",
		"бла бла бла купил": "А не пиздишь? Аренда это не покупка",
		"бла бла бла купил бла бла бла": "А не пиздишь? Аренда это не покупка",

		"spotify":  "Эти пидоры Антону косарик должны за подписку",
		"Spotify":  "Эти пидоры Антону косарик должны за подписку",
		"спотифай": "Эти пидоры Антону косарик должны за подписку",
		"бла бла бла spotify бла бла бла":  "Эти пидоры Антону косарик должны за подписку",
		"бла бла бла спотифай бла бла бла": "Эти пидоры Антону косарик должны за подписку",

		"devops": "Девопсы не нужны",
		"девопс": "Девопсы не нужны",

		"трансформация":                "оргия гомогеев",
		"трансформацию":                "оргию гомогеев",
		"трансформации":                "оргии гомогеев",
		"сейчас трансформация пройдет": "сейчас оргия гомогеев пройдет",
		"сейчас трансформацию пройдет": "сейчас оргию гомогеев пройдет",
		"сейчас трансформации пройдет": "сейчас оргии гомогеев пройдет",
		"с этими трансформациями уже":  "с этими оргиями гомогеев уже",

		"Я не понял - вот делают трансформацию, ч": "я не понял - вот делают оргию гомогеев, ч",

		"java":          "джава-хуява, а я работаю на го",
		"бла java бла":  "джава-хуява, а я работаю на го",
		"джава":         "джава-хуява, а я работаю на го",
		"бла джава бла": "джава-хуява, а я работаю на го",
		"джаба":         "джава-хуява, а я работаю на го",
		"бла джаба бла": "джава-хуява, а я работаю на го",

		"да там же уже решено, теперь надо новое что-то заблокировать": "пусть себе анус заблокируют",
		"они это вчера заблокировали и ждут":                           "пусть себе анус заблокируют",
		"это меня блокирует сильно":                                    "пусть себе анус заблокируют",

		"я думал сначала Медведев это опять. А тут какой то давыдов": "не опять, а снова",

		"у нас проблема":           "у меня есть 5-10 солюшенов этой проблемы",
		"у нас проблема, товарищи": "у меня есть 5-10 солюшенов этой проблемы",
		"проблема в том":           "у меня есть 5-10 солюшенов этой проблемы",
		"куча проблем":             "у меня есть 5-10 солюшенов этой проблемы",
	}
	for k, expected := range data {
		respond, response := brain.Decision(0, k)
		if !respond || response != expected {
			t.Errorf("Expected: \"%s\" but got: \"%s\"", expected, response)
		}
	}
}

func Test_returnsOnNotElderRing(t *testing.T) {
	brain := NewBrain(NewMemory(), true)
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

func Test_returnsForLuckyKhaleesifiedText(t *testing.T) {
	brain := NewBrain(NewMemory(), true)
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

func _Test_avoidsSensitiveTopics(t *testing.T) {
	brain := NewBrain(NewMemory(), true)
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
	for i := 0; i < 500; i++ {
		for _, text := range data {
			respond, response := brain.Decision(0, text)
			if respond {
				t.Error("Response for ", text, ": not expected != ", response)
			}
		}
	}
}
