package brain

import (
	"sort"
)

type Memory struct {
	senselessPhrases []string
	mockingMap       map[string]string
	mockingMapKeys   []string
	normalisationMap map[string]string
}

func NewMemory() *Memory {
	senselessPhrases := make([]string, 10)
	mockingMap := make(map[string]string)
	mockingMapKeys := make([]string, 10)
	normalisationMap := make(map[string]string)

	memory := &Memory{
		senselessPhrases: senselessPhrases,
		mockingMap:       mockingMap,
		mockingMapKeys:   mockingMapKeys,
		normalisationMap: normalisationMap,
	}
	return memory
}

func (m *Memory) GetSenslessPhrases() []string {
	return m.senselessPhrases
}

func (m *Memory) GetMockingMap() map[string]string {
	return m.mockingMap
}

func (m *Memory) GetMokingMapKeys() []string {
	return m.mockingMapKeys
}

func (m *Memory) GetNormalizationMap() map[string]string {
	return m.normalisationMap
}

func (m *Memory) SetSenslessPhrases(list []string) {
	for i := 0; i < len(list); i++ {
		m.senselessPhrases = append(m.senselessPhrases, list[i])
	}

}

func (m *Memory) SetMockingMap(mockingMap map[string]string) {
	newMap := make(map[string]string)
	keys := make([]string, len(mockingMap))
	i := 0
	for k, v := range mockingMap {
		newMap[k] = v
		keys[i] = k
		i++
	}
	/*
		i := 0
		for k := range mockingMap {
			keys[i] = k
			i++
		}
	*/
	sort.Strings(keys)
	m.mockingMap = newMap
	m.mockingMapKeys = keys
}

func (m *Memory) SetNormalizationMap(normalisationMap map[string]string) {
	newMap := make(map[string]string)
	for k, v := range normalisationMap {
		newMap[k] = v
	}
	m.normalisationMap = newMap
}
