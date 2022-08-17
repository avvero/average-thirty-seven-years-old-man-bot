package brain

import (
	"github.com/avvero/the_gamers_guild_bot/internal/knowledge"
	"sort"
)

type Memory struct {
	senselessPhrases []string
	mockingMap       map[string]string
	mockingMapKeys   []string
	normalisationMap map[string]string
}

func NewMemory() *Memory {
	mockingMapKeys := make([]string, len(knowledge.MockingMap))
	i := 0
	for k := range knowledge.MockingMap {
		mockingMapKeys[i] = k
		i++
	}
	sort.Strings(mockingMapKeys)
	//
	memory := &Memory{
		senselessPhrases: knowledge.SenselessPhrases,
		mockingMap:       knowledge.MockingMap,
		mockingMapKeys:   mockingMapKeys,
		normalisationMap: knowledge.NormalisationMap,
	}
	return memory
}
