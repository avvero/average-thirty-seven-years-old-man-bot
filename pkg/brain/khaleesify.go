package brain

import (
	"strings"
)

type KhaleesifyIntention struct {
	memory *Memory
}

func NewKhaleesifyIntention(memory *Memory) *KhaleesifyIntention {
	return &KhaleesifyIntention{memory: memory}
}

func (this KhaleesifyIntention) Express(text string) (has bool, response string) {
	result := strings.ToLower(text)
	for _, k := range this.memory.mockingMapKeys {
		result = strings.Replace(result, k, this.memory.mockingMap[k], -1)
	}
	return true, result
}
