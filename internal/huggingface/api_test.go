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

	apiClient := NewApiClient(ts.URL, "key")
	score, err := apiClient.ToxicityScore("any")
	expected := 0.0014224671758711338
	if err != nil {
		t.Error("Error: ", fmt.Sprintf("%s", err))
	}
	if score != expected {
		t.Error("Score expected: ", expected, " != ", score)
	}
}

func Test_toxicityScoreReturnsErrorIfCantParse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{dsf}`)
	}))
	defer ts.Close()

	apiClient := NewApiClient(ts.URL, "key")
	score, err := apiClient.ToxicityScore("any")
	expected := "Can't parse response: {dsf}\n: invalid character 'd' looking for beginning of object key string"
	if fmt.Sprintf("%s", err) != expected {
		t.Error("Error expected: ", expected, " != ", fmt.Sprintf("%s", err))
	}
	if score != 0 {
		t.Error("Score expected: ", 0, " != ", score)
	}
}

func Test_toxicityScoreReturnsErrorIfLabel1IsNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[[{"label":"LABEL_0","score":0.9985774755477905}]]`)
	}))
	defer ts.Close()

	apiClient := NewApiClient(ts.URL, "key")
	score, err := apiClient.ToxicityScore("any")
	expected := "Can't find LABEL_1: [[{LABEL_0 0.9985774755477905}]]"
	if fmt.Sprintf("%s", err) != expected {
		t.Error("Error expected: ", expected, " != ", fmt.Sprintf("%s", err))
	}
	if score != 0 {
		t.Error("Score expected: ", 0, " != ", score)
	}
}

func Test_toxicityScoreReturnsErrorOn503(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
		fmt.Fprintln(w, `{"error": "Some error", "estimated_time": 23}`)
	}))
	defer ts.Close()

	apiClient := NewApiClient(ts.URL, "key")
	score, err := apiClient.ToxicityScore("any")
	expected := fmt.Sprintf("Response from: %s: 503: Some error, estimated time: 23.000000", apiClient.url)
	if fmt.Sprintf("%s", err) != expected {
		t.Error("Error expected: ", expected, " != ", fmt.Sprintf("%s", err))
	}
	if score != 0 {
		t.Error("Score expected: ", 0, " != ", score)
	}
}

func Test_toxicityScoreReturnsErrorOn502(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(502)
		fmt.Fprintln(w, `{"error": "Some error", "estimated_time": 23}`)
	}))
	defer ts.Close()

	apiClient := NewApiClient(ts.URL, "key")
	score, err := apiClient.ToxicityScore("any")
	expected := fmt.Sprintf("Response from: %s: %d", apiClient.url, 502)
	if fmt.Sprintf("%s", err) != expected {
		t.Error("Error expected: ", expected, " != ", fmt.Sprintf("%s", err))
	}
	if score != 0 {
		t.Error("Score expected: ", 0, " != ", score)
	}
}
