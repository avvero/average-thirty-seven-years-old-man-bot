package memory

import "sort"

type Memory struct {
	senselessPhrases []string
	mockingMap       map[string]string
	mockingMapKeys   []string
	normalisationMap map[string]string
}

func NewMemory() *Memory {
	//
	senselessPhrases := []string{"хуйню не неси", "база", "мда", "вообще похую", "ничего нового", "хули нам кабанам",
		"кринж", "норм", "такое себе", "априори", "не комильфо от слова совсем", "ебаный цирк", "болие лимение", "оп-хуй",
		"вишенка на торте", "мякотка", "пруфай", "найс", "несолоно хлебавши", "и что не так?",
		"имеет место быть", "с моей колокольни", "рабочий кейс", "лалка", "я тебя услышал", "ну такое", "внимательно",
		"сам-то понял что сказал?", "шта", "ржомба", "литерали", "провёл ресеч", "бывает", "частый кейс", "по факту",
		"двачую", "охуел с твоей истории", "это другое", "ну это провал", "храни тебя господь", "твою бога душу мать",
		"претенциозно", "вот раньше лучше было", "рил", "по кд", "тоси-боси", "ор"}
	//
	mockingMap := map[string]string{
		"ль": "й",
		"ри": "и",
		"ре": "ри",
		"ра": "ья",
		"за": "зя",
		"ол": "ой",
		"ме": "ми",
		"мн": "мен",
		"те": "ти",
		"не": "ни",
		"се": "си",
		//"го": "во",
		"ыл": "ыль",
		"он": "онь",
		"вс": "фс",
		"го": "кхо",
		"а":  "я",
		//"е":  "и",
		"р": "л",
		"ж": "з",
		//"в": "ф",
		"": "",
	}
	mockingMapKeys := make([]string, len(mockingMap))
	i := 0
	for k := range mockingMap {
		mockingMapKeys[i] = k
		i++
	}
	sort.Strings(mockingMapKeys)
	//
	normalisationMap := map[string]string{
		"e": "е",
		"o": "о",
		"h": "н",
		"a": "а",
		"t": "т",
		"k": "к",
		"c": "с",
		"b": "б",
		"m": "м",
	}
	//
	memory := &Memory{
		senselessPhrases: senselessPhrases,
		mockingMap:       mockingMap,
		mockingMapKeys:   mockingMapKeys,
		normalisationMap: normalisationMap,
	}
	return memory
}

func (m *Memory) GetSenslessPhrases() []string {
	return m.senselessPhrases
}

func (m *Memory) GetMockingMap() map[string]string {
	return m.mockingMap
}

func (m *Memory) GetMokingMapKeys() []string {
	return m.mockingMapKeys
}

func (m *Memory) GetNormalizationMap() map[string]string {
	return m.normalisationMap
}
