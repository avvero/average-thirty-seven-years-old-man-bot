package huggingface

import (
	"testing"
)

func Test_parseTest(t *testing.T) {
	text := `[[{"label":"LABEL_0","score":0.9985774755477905},{"label":"LABEL_1","score":0.0014224671758711338}]]`
	labels, err := Parse(text)
	if err != nil {
		t.Error("Parse error ", err)
	}
	found := false
	for _, label := range labels {
		if label.Title == "LABEL_1" {
			found = true
		}
	}
	if !found {
		t.Error("Parse unexpected result: ", labels)
	}
}

func Test_apiTest(t *testing.T) {
	accessKey := "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"
	apiClient := NewHuggingFaceApiClient(accessKey)
	toxicityScore, err := apiClient.ToxicityScore("привет")
	if err != nil {
		t.Error("Parse error ", err)
	}
	t.Log("Toxicity score", toxicityScore)
	if toxicityScore == 0 {
		t.Error("Unexpected toxicity score ", toxicityScore)
	}
}
