package brain

import (
	"github.com/avvero/the_gamers_guild_bot/internal/knowledge"
)

type Memory struct {
	normalisationMap map[string]string
}

func NewMemory() *Memory {
	//
	memory := &Memory{
		normalisationMap: knowledge.NormalisationMap,
	}
	return memory
}
