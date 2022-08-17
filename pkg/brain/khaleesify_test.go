package brain

import (
	"testing"
)

func Test_khaleesifiesText(t *testing.T) {
	intention := &KhaleesifyIntention{memory: NewMemory()}
	data := map[string]string{
		"Позвольте мне сражаться за Вас, Кхалиси":                                   "позвойти мени слязяться зя вяс, кхялиси",
		"дерись за меня, дракон":                                                    "делись зя миня, дляконь",
		"Мне кажется ягодки это какой-то сайт для секс знакомств должен быть":       "мени кязется якходки это кякой-то сяйт для сикс знякомств дойзен быть",
		"не время бухтеть":                                                          "ни влемя бухтить",
		"Сегодня был созвон со всеми разработчиками и всем осветили будущую модель": "сикходня быль созвонь со фсими лязляботчикями и фсим осветили будущую модей",
	}
	for k, expected := range data {
		_, result := intention.Express(k)
		if result != expected {
			t.Error("Expected and got:", expected, " != ", result)
		}
	}
}
