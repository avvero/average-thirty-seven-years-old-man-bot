package brain

import "github.com/avvero/the_gamers_guild_bot/internal/utils"

type SenselessPhrasesIntention struct {
}

var senselessPhrases = []string{
	"хуйню не неси",
	"база",
	"мда",
	"вообще похую",
	"ничего нового",
	"хули нам кабанам",
	"хули нам канбанам",
	"кринж",
	"норм",
	"такое себе",
	"априори",
	"не комильфо от слова совсем",
	"ебаный цирк",
	"болие лимение",
	"оп-хуй",
	"вишенка на торте",
	"мякотка",
	"пруфай",
	"найс",
	"несолоно хлебавши",
	"и что не так?",
	"имеет место быть",
	"с моей колокольни",
	"рабочий кейс",
	"лалка",
	"я тебя услышал",
	"ну такое",
	"внимательно",
	"сам-то понял что сказал?",
	"шта",
	"ржомба",
	"литерали",
	"провёл ресеч",
	"бывает",
	"частый кейс",
	"по факту",
	"двачую",
	"охуел с твоей истории",
	"это другое",
	"ну это провал",
	"храни тебя господь",
	"твою бога душу мать",
	"претенциозно",
	"вот раньше лучше было",
	"рил",
	"по кд",
	"тоси-боси",
	"ор",
	"действительно",
	"ну точно",
	"сказал, как в лужу пёрнул",
	"в общем и целом",
	"держи в курсе, ебло ослиное",
}

func (this SenselessPhrasesIntention) Express(text string) (has bool, response string) {
	return true, senselessPhrases[utils.RandomUpTo(len(senselessPhrases))]
}
