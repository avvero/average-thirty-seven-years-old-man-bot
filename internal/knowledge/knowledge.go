package knowledge

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

var Vowels = []rune{'а', 'е', 'ё', 'и', 'о', 'у', 'ы', 'э', 'ю', 'я'}
var VowelsSoftenMap = map[rune]rune{
	'о': 'ё',
	'ы': 'и',
	'а': 'я',
	'у': 'ю',
}
var Delimiters = []rune{' ', '.', ',', ':', '!', '?', '/', ';', '\'', '"', '#', '$', '(', ')', '-'}
