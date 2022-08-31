package brain

import "github.com/avvero/the_gamers_guild_bot/internal/utils"

type ConfuciusPhrasesIntention struct {
}

var confuciusPhrases = []string{
	"First, by reflection, which is noblest;\nSecond, by imitation, which is easiest;\nThird by experience, which is the bitterest.",
	"It does not matter how slowly you go as long as you do not stop.",
	"Wheresoever you go, go with all your heart.",
	"The man who moves a mountain begins by carrying away small stones.",
	"Life is really simple, but we insist on making it complicated.",
	"Study the past if you would define the future.",
	"If you make a mistake and do not correct it, this is called a mistake.",
	"Before you embark on a journey of revenge, dig two graves.",
	"Respect yourself and others will respect you.",
	"The man who asks a question is a fool for a minute, the man who does not ask is a fool for life.",
	"Silence is a true friend who never betrays.",
	"I hear and I forget. I see and I remember. I do and I understand.",
	"The hardest thing of all is to find a black cat in a dark room, especially if there is no cat.",
	"Give a bowl of rice to a man and you will feed him for a day.",
	"Teach him how to grow his own rice and you will save his life.",
	"We have two lives, and the second begins when we realize we only have one.",
}

func (this ConfuciusPhrasesIntention) Express(text string) (has bool, response string) {
	return true, confuciusPhrases[utils.RandomUpTo(len(confuciusPhrases))]
}
