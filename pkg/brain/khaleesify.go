package brain

import (
	"sort"
	"strings"
)

var mockingMap = map[string]string{
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

type KhaleesifyIntention struct {
	mockingMap     map[string]string
	mockingMapKeys []string
}

func NewKhaleesifyIntention() *KhaleesifyIntention {
	mockingMapKeys := make([]string, len(mockingMap))
	i := 0
	for k := range mockingMap {
		mockingMapKeys[i] = k
		i++
	}
	sort.Strings(mockingMapKeys)
	//
	return &KhaleesifyIntention{mockingMap: mockingMap, mockingMapKeys: mockingMapKeys}
}

func (this KhaleesifyIntention) Express(text string) (has bool, response string) {
	result := strings.ToLower(text)
	for _, k := range this.mockingMapKeys {
		result = strings.Replace(result, k, this.mockingMap[k], -1)
	}
	return true, result
}
