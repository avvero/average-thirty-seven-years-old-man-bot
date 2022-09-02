package huggingface

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_toxicityScoreReturnsScore(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[[{"label":"LABEL_0","score":0.9985774755477905},{"label":"LABEL_1","score":0.0014224671758711338}]]`)
	}))
	defer ts.Close()

	apiClient := NewHuggingFaceApiClient(ts.URL, "key")
	score, _ := apiClient.ToxicityScore("any")
	expected := 0.0014224671758711338
	if score != expected {
		t.Error("Score: ", expected, " != ", score)
	}
}

func Test_toxicityScoreReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
		fmt.Fprintln(w, `{"error": "Some error", "estimated_time": 23}`)
	}))
	defer ts.Close()

	apiClient := NewHuggingFaceApiClient(ts.URL, "key")
	score, err := apiClient.ToxicityScore("any")
	expected := "Response code 503: Some error, estimated time: 23.000000"
	if fmt.Sprintf("%s", err) != expected {
		t.Error("Error: ", expected, " != ", fmt.Sprintf("%s", err))
	}
	if score != 0 {
		t.Error("Score: ", 0, " != ", score)
	}
}
