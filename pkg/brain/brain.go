package brain

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	m "github.com/avvero/the_gamers_guild_bot/pkg/memory"
	u "github.com/avvero/the_gamers_guild_bot/pkg/utils"
)

type Brain struct {
	memory *m.Memory
}

func NewBrain(memory *m.Memory) *Brain {
	return &Brain{memory: memory}
}

func (b *Brain) Decision(chatId int64, text string, rnd bool) (respond bool, response string) {
	if rnd {
		if randomUpTo(50) == 0 {
			phrase := b.memory.GetSenslessPhrases()[randomUpTo(len(b.memory.GetSenslessPhrases()))]
			return true, phrase
		}
		if len(text) > 14 && randomUpTo(50) == 0 {
			phrase := b.khaleesify(text)
			return true, phrase
		}
		if !u.Contains([]string{"0", "-1001733786877", "245851441", "-578279468"}, strconv.FormatInt(chatId, 10)) {
			return true, "Mr Moony presents his compliments to Professor Snape, and begs him to keep his abnormally large nose out of other people’s business."
		}
	}
	//
	text = strings.ToLower(text)
	if text == "gg" {
		return true, "gg"
	}
	if b.normalizeRu(text) == "нет" {
		return true, "пидора ответ"
	}
	if strings.Contains(text, "morrowind") ||
		strings.Contains(text, "морровинд") ||
		strings.Contains(text, "моровинд") {
		return true, "Morrowind - одна из лучших игр эва"
	}
	if text == "er" ||
		text == "ер" ||
		text == "эр" ||
		strings.Contains(text, "elden ring") ||
		strings.Contains(text, " er ") ||
		strings.Contains(text, " ер ") ||
		strings.Contains(text, " эр ") {
		return true, "Elden Ring - это величие"
	}
	if strings.Contains(text, "купил") {
		return true, "А не пиздишь? Аренда это не покупка"
	}
	if strings.Contains(text, "spotify") || strings.Contains(text, "спотифай") {
		return true, "Эти пидоры Антону косарик должны за подписку"
	}
	if strings.Contains(b.normalizeEn(text), "devops") ||
		strings.Contains(b.normalizeRu(text), "девопс") {
		return true, "Девопсы не нужны"
	}
	return false, ""
}

func randomUpTo(max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max)
}

func (b *Brain) normalizeRu(text string) string {
	result := text
	for k, v := range b.memory.GetNormalizationMap() {
		result = strings.Replace(result, k, v, -1)
	}
	return result
}

func (b *Brain) normalizeEn(text string) string {
	result := text
	for k, v := range b.memory.GetNormalizationMap() {
		result = strings.Replace(result, v, k, -1)
	}
	return result
}

func (brain Brain) khaleesify(text string) string {
	result := strings.ToLower(text)
	for _, k := range brain.memory.GetMokingMapKeys() {
		result = strings.Replace(result, k, brain.memory.GetMockingMap()[k], -1)
	}
	return result
}
