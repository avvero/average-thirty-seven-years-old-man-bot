package knowledge

var SenselessPhrases = []string{"хуйню не неси", "база", "мда", "вообще похую", "ничего нового", "хули нам кабанам",
	"кринж", "норм", "такое себе", "априори", "не комильфо от слова совсем", "ебаный цирк", "болие лимение", "оп-хуй",
	"вишенка на торте", "мякотка", "пруфай", "найс", "несолоно хлебавши", "и что не так?",
	"имеет место быть", "с моей колокольни", "рабочий кейс", "лалка", "я тебя услышал", "ну такое", "внимательно",
	"сам-то понял что сказал?", "шта", "ржомба", "литерали", "провёл ресеч", "бывает", "частый кейс", "по факту",
	"двачую", "охуел с твоей истории", "это другое", "ну это провал", "храни тебя господь", "твою бога душу мать",
	"претенциозно", "вот раньше лучше было", "рил", "по кд", "тоси-боси", "ор"}

var MockingMap = map[string]string{
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

var NormalisationMap = map[string]string{
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
